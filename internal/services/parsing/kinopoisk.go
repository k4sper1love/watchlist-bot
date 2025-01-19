package parsing

import (
	"encoding/json"
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"io"
	"log"
	"net/http"
	"net/url"
)

func GetFilmFromKinopoisk(app models.App, url string) (*apiModels.Film, error) {
	queryKey, id, err := utils.ExtractKinopoiskQuery(url)
	if err != nil {
		return nil, err
	}

	headers := map[string]string{
		"X-API-KEY": app.Vars.KinopoiskAPIToken,
	}

	reqUrl := fmt.Sprintf("https://api.kinopoisk.dev/v1.4/movie?%s=%s", queryKey, id)

	resp, err := client.SendRequestWithOptions(reqUrl, http.MethodGet, nil, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed response. Status is %s", resp.Status)
	}

	film, err := parseFilmFromKinopoisk(resp.Body)
	if err != nil {
		return nil, err
	}

	return film, nil
}

func GetFilmsFromKinopoisk(app models.App, session *models.Session) ([]apiModels.Film, *filters.Metadata, error) {
	headers := map[string]string{
		"X-API-KEY":  app.Vars.KinopoiskAPIToken,
		"User-Agent": "PostmanRuntime/7.31.1",
	}

	state := session.FilmsState

	query := url.QueryEscape(state.Title)
	reqUrl := fmt.Sprintf("https://api.kinopoisk.dev/v1.4/movie/search?page=%d&limit=%d&query=%s", state.CurrentPage, state.PageSize, query)

	resp, err := client.SendRequestWithOptions(reqUrl, http.MethodGet, nil, headers)
	if err != nil {
		return nil, nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, nil, fmt.Errorf("failed response. Status is %s", resp.Status)
	}

	films, metadata, err := parseFilmsFromKinopoisk(resp.Body)
	if err != nil {
		return nil, nil, err
	}

	log.Println(films)

	return films, metadata, nil
}

func parseFilmFromKinopoisk(data io.Reader) (*apiModels.Film, error) {
	var response map[string]interface{}
	if err := json.NewDecoder(data).Decode(&response); err != nil {
		return nil, err
	}

	docs, ok := response["docs"].([]interface{})
	if !ok || len(docs) == 0 {
		return nil, fmt.Errorf("no films found in response")
	}

	filmData, ok := docs[0].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid film data format")
	}

	film := parseFilmDataKinopoisk(filmData)

	return film, nil
}

func parseFilmsFromKinopoisk(data io.Reader) ([]apiModels.Film, *filters.Metadata, error) {
	var response map[string]interface{}
	if err := json.NewDecoder(data).Decode(&response); err != nil {
		return nil, nil, err
	}

	docs, ok := response["docs"].([]interface{})
	if !ok || len(docs) == 0 {
		return nil, nil, fmt.Errorf("no films found in response")
	}

	var films []apiModels.Film
	for _, doc := range docs {
		filmData, ok := doc.(map[string]interface{})
		if !ok {
			continue
		}
		film := parseFilmDataKinopoisk(filmData)
		film.URL = fmt.Sprintf("https://www.kinopoisk.ru/film/%d/", film.ID)
		films = append(films, *film)
	}

	var metadata filters.Metadata
	if total, ok := response["total"].(float64); ok {
		metadata.TotalRecords = int(total)
	}
	if limit, ok := response["limit"].(float64); ok {
		metadata.PageSize = int(limit)
	}
	if page, ok := response["page"].(float64); ok {
		metadata.CurrentPage = int(page)
	}
	if pages, ok := response["pages"].(float64); ok {
		metadata.LastPage = int(pages)
	}

	return films, &metadata, nil
}

func parseFilmDataKinopoisk(data map[string]interface{}) *apiModels.Film {
	var film apiModels.Film

	if id, ok := data["id"].(float64); ok {
		film.ID = int(id)
	}

	if title, ok := data["name"].(string); ok {
		film.Title = title
	}

	if year, ok := data["year"].(float64); ok {
		film.Year = int(year)
	}

	if genres, ok := data["genres"].([]interface{}); ok && len(genres) > 0 {
		if firstGenre, ok := genres[0].(map[string]interface{}); ok {
			film.Genre, _ = firstGenre["name"].(string)
		}
	}

	if description, ok := data["description"].(string); ok {
		film.Description = description
	}

	if rating, ok := data["rating"].(map[string]interface{}); ok {
		if kpRating, ok := rating["kp"].(float64); ok {
			film.Rating = kpRating
		}
	}

	if image, ok := data["poster"].(map[string]interface{}); ok {
		if imageURL, ok := image["url"].(string); ok {
			film.ImageURL = imageURL
		}
	}

	return &film
}
