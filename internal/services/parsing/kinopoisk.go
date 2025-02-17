package parsing

import (
	"encoding/json"
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

func GetFilmFromKinopoisk(session *models.Session, url string) (*apiModels.Film, error) {
	queryKey, id, err := utils.ExtractKinopoiskQuery(url)
	if err != nil {
		sl.Log.Error("failed to extract query", slog.Any("error", err), slog.String("url", url))
		return nil, err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderExternalAPIKey,
			HeaderValue:        session.KinopoiskAPIToken,
			Method:             http.MethodGet,
			URL:                fmt.Sprintf("https://api.kinopoisk.dev/v1.4/movie?%s=%s", queryKey, id),
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	film, err := parseFilmFromKinopoisk(resp.Body)
	if err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return film, nil
}

func GetFilmsFromKinopoisk(session *models.Session) ([]apiModels.Film, *filters.Metadata, error) {
	state := session.FilmsState
	query := url.QueryEscape(state.Title)

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderExternalAPIKey,
			HeaderValue:        session.KinopoiskAPIToken,
			Method:             http.MethodGet,
			URL:                fmt.Sprintf("https://api.kinopoisk.dev/v1.4/movie/search?page=%d&limit=%d&query=%s", state.CurrentPage, state.PageSize, query),
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return nil, nil, err
	}
	defer utils.CloseBody(resp.Body)

	films, metadata, err := parseFilmsFromKinopoisk(resp.Body)
	if err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, nil, err
	}

	return films, metadata, nil
}

func parseFilmFromKinopoisk(data io.Reader) (*apiModels.Film, error) {
	var response map[string]interface{}
	if err := json.NewDecoder(data).Decode(&response); err != nil {
		return nil, err
	}

	docs, ok := response["docs"].([]interface{})
	if !ok || len(docs) == 0 {
		return nil, fmt.Errorf("film not found in response")
	}

	filmData, ok := docs[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid film data format")
	}

	return parseFilmDataKinopoisk(filmData)
}

func parseFilmsFromKinopoisk(data io.Reader) ([]apiModels.Film, *filters.Metadata, error) {
	var response map[string]interface{}
	if err := json.NewDecoder(data).Decode(&response); err != nil {
		return nil, nil, err
	}

	docs, ok := response["docs"].([]interface{})
	if !ok || len(docs) == 0 {
		return []apiModels.Film{}, &filters.Metadata{}, nil
	}

	var films []apiModels.Film
	for _, doc := range docs {
		filmData, ok := doc.(map[string]interface{})
		if !ok {
			continue
		}

		film, err := parseFilmDataKinopoisk(filmData)
		if err != nil {
			continue
		}

		film.URL = fmt.Sprintf("https://www.kinopoisk.ru/film/%d/", film.ID)
		films = append(films, *film)
	}

	metadata := filters.Metadata{
		TotalRecords: client.GetIntFromMap(response, "total", 0),
		PageSize:     client.GetIntFromMap(response, "limit", 0),
		CurrentPage:  client.GetIntFromMap(response, "page", 0),
		LastPage:     client.GetIntFromMap(response, "pages", 0),
	}

	return films, &metadata, nil
}

func parseFilmDataKinopoisk(data map[string]interface{}) (*apiModels.Film, error) {
	film := apiModels.Film{
		ID:          client.GetIntFromMap(data, "id", 0),
		Title:       client.GetStringFromMap(data, "name", "Unknown"),
		Year:        client.GetIntFromMap(data, "year", 0),
		Genre:       getFirstGenre(data, ""),
		Description: client.GetStringFromMap(data, "description", ""),
		Rating:      client.GetFloatFromNestedMap(data, "rating", "kp", 0.0),
		ImageURL:    client.GetStringFromNestedMap(data, "poster", "url", ""),
	}

	return &film, nil
}

func getFirstGenre(data map[string]interface{}, defaultValue string) string {
	if genres, ok := data["genres"].([]interface{}); ok && len(genres) > 0 {
		if firstGenre, ok := genres[0].(map[string]interface{}); ok {
			if genreName, ok := firstGenre["name"].(string); ok {
				return genreName
			}
		}
	}
	return defaultValue
}
