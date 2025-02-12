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

func GetFilmFromKinopoisk(app models.App, url string) (*apiModels.Film, error) {
	queryKey, id, err := utils.ExtractKinopoiskQuery(url)
	if err != nil {
		sl.Log.Error("failed to extract query", slog.Any("error", err), slog.String("url", url))
		return nil, err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderExternalAPIKey,
			HeaderValue:        app.Vars.KinopoiskAPIToken,
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

func GetFilmsFromKinopoisk(app models.App, session *models.Session) ([]apiModels.Film, *filters.Metadata, error) {
	state := session.FilmsState
	query := url.QueryEscape(state.Title)

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderExternalAPIKey,
			HeaderValue:        app.Vars.KinopoiskAPIToken,
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
		return nil, fmt.Errorf("not found film in response")
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
			return nil, nil, err
		}

		film.URL = fmt.Sprintf("https://www.kinopoisk.ru/film/%d/", film.ID)
		films = append(films, *film)
	}

	var metadata filters.Metadata

	if total, ok := response["total"].(float64); ok {
		metadata.TotalRecords = int(total)
	} else {
		return nil, nil, fmt.Errorf("failed to parse total")
	}

	if limit, ok := response["limit"].(float64); ok {
		metadata.PageSize = int(limit)
	} else {
		return nil, nil, fmt.Errorf("failed to parse limit")
	}

	if page, ok := response["page"].(float64); ok {
		metadata.CurrentPage = int(page)
	} else {
		return nil, nil, fmt.Errorf("failed to parse page")
	}

	if pages, ok := response["pages"].(float64); ok {
		metadata.LastPage = int(pages)
	} else {
		return nil, nil, fmt.Errorf("failed to parse pages")
	}

	return films, &metadata, nil
}

func parseFilmDataKinopoisk(data map[string]interface{}) (*apiModels.Film, error) {
	var film apiModels.Film

	if id, ok := data["id"].(float64); ok {
		film.ID = int(id)
	} else {
		return nil, fmt.Errorf("failed to parse id")
	}

	if title, ok := data["name"].(string); ok {
		film.Title = title
	} else {
		return nil, fmt.Errorf("failed to parse title")
	}

	if year, ok := data["year"].(float64); ok {
		film.Year = int(year)
	} else {
		return nil, fmt.Errorf("failed to parse year")
	}

	if genres, ok := data["genres"].([]interface{}); ok && len(genres) > 0 {
		if firstGenre, ok := genres[0].(map[string]interface{}); ok {
			film.Genre, _ = firstGenre["name"].(string)
		} else {
			return nil, fmt.Errorf("failed to parse genre")
		}
	} else {
		return nil, fmt.Errorf("failed to parse genres")
	}

	if description, ok := data["description"].(string); ok {
		film.Description = description
	} else {
		return nil, fmt.Errorf("failed to parse description")
	}

	if rating, ok := data["rating"].(map[string]interface{}); ok {
		if kpRating, ok := rating["kp"].(float64); ok {
			film.Rating = kpRating
		} else {
			return nil, fmt.Errorf("failed to parse rating")
		}
	} else {
		return nil, fmt.Errorf("failed to parse rating")
	}

	if image, ok := data["poster"].(map[string]interface{}); ok {
		if imageURL, ok := image["url"].(string); ok {
			film.ImageURL = imageURL
		} else {
			return nil, fmt.Errorf("failed to parse image url")
		}
	} else {
		return nil, fmt.Errorf("failed to parse poster")
	}

	return &film, nil
}
