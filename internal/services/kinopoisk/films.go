package kinopoisk

import (
	"encoding/json"
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"io"
	"log"
	"net/http"
)

func GetFilmByID(app models.App, filmID int) (*apiModels.Film, error) {
	headers := map[string]string{
		"X-API-KEY": app.Vars.KinopoiskAPIToken,
	}

	reqUrl := fmt.Sprintf("https://api.kinopoisk.dev/v1.4/movie/%d", filmID)

	resp, err := client.SendRequest(reqUrl, http.MethodGet, nil, headers)
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		log.Println(resp.Status)
		return nil, err
	}

	var film apiModels.Film
	if err := parseFilm(&film, resp.Body); err != nil {
		log.Println(err)
		return nil, err
	}

	return &film, nil
}

func parseFilm(dest *apiModels.Film, data io.Reader) error {
	var response map[string]interface{}
	if err := json.NewDecoder(data).Decode(&response); err != nil {
		return err
	}

	if title, ok := response["name"].(string); ok {
		dest.Title = title
	}

	if year, ok := response["year"].(float64); ok {
		dest.Year = int(year)
	}

	if genres, ok := response["genres"].([]interface{}); ok && len(genres) > 0 {
		if firstGenre, ok := genres[0].(map[string]interface{}); ok {
			dest.Genre, _ = firstGenre["name"].(string)
		}
	}
	if description, ok := response["description"].(string); ok {
		dest.Description = description
	}
	if rating, ok := response["rating"].(map[string]interface{}); ok {
		if kpRating, ok := rating["kp"].(float64); ok {
			dest.Rating = kpRating
		}
	}
	if image, ok := response["poster"].(map[string]interface{}); ok {
		if url, ok := image["url"].(string); ok {
			dest.ImageURL = url
		}
	}

	log.Printf("Fetched film: %+v\n", dest)
	return nil
}
