package parsing

import (
	"encoding/json"
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"io"
	"net/http"
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
		return nil, err
	}

	var film apiModels.Film
	if err := parseFilmFromKinopoisk(&film, resp.Body); err != nil {
		return nil, err
	}

	return &film, nil
}

func parseFilmFromKinopoisk(dest *apiModels.Film, data io.Reader) error {
	var response map[string]interface{}
	if err := json.NewDecoder(data).Decode(&response); err != nil {
		return err
	}

	docs, ok := response["docs"].([]interface{})
	if !ok || len(docs) == 0 {
		return fmt.Errorf("no films found in response")
	}

	filmData, ok := docs[0].(map[string]interface{})
	if !ok {
		return fmt.Errorf("invalid film data format")
	}

	if title, ok := filmData["name"].(string); ok {
		dest.Title = title
	}

	if year, ok := filmData["year"].(float64); ok {
		dest.Year = int(year)
	}

	if genres, ok := filmData["genres"].([]interface{}); ok && len(genres) > 0 {
		if firstGenre, ok := genres[0].(map[string]interface{}); ok {
			dest.Genre, _ = firstGenre["name"].(string)
		}
	}

	if description, ok := filmData["description"].(string); ok {
		dest.Description = description
	}

	if rating, ok := filmData["rating"].(map[string]interface{}); ok {
		if kpRating, ok := rating["kp"].(float64); ok {
			dest.Rating = kpRating
		}
	}

	if image, ok := filmData["poster"].(map[string]interface{}); ok {
		if url, ok := image["url"].(string); ok {
			dest.ImageURL = url
		}
	}

	return nil
}
