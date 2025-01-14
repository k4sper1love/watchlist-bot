package parsing

import (
	"encoding/json"
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"io"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func GetFilmFromIMDB(app models.App, url string) (*apiModels.Film, error) {
	id := parseIDFromIMDB(url)

	reqUrl := fmt.Sprintf("http://www.omdbapi.com/?apikey=%s&i=%s", app.Vars.IMDBAPIToken, id)

	resp, err := client.SendRequestWithOptions(reqUrl, http.MethodGet, nil, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return nil, fmt.Errorf("failed response. Status is %s", resp.Status)
	}

	var film apiModels.Film
	if err := parseFilmFromIMDB(&film, resp.Body); err != nil {
		return nil, err
	}

	return &film, nil
}

func parseFilmFromIMDB(dest *apiModels.Film, data io.Reader) error {
	var response map[string]interface{}
	if err := json.NewDecoder(data).Decode(&response); err != nil {
		return err
	}

	if title, ok := response["Title"].(string); ok {
		dest.Title = title
	}

	if yearStr, ok := response["Year"].(string); ok {
		year, err := strconv.Atoi(yearStr)
		if err == nil {
			dest.Year = year
		} else {
			log.Printf("Failed to parse Year: %v", err)
		}
	}

	if genres, ok := response["Genre"].(string); ok {
		dest.Genre = strings.Split(genres, ",")[0]
	}

	if description, ok := response["Plot"].(string); ok {
		dest.Description = description
	}

	if ratingStr, ok := response["imdbRating"].(string); ok {
		rating, err := strconv.ParseFloat(ratingStr, 64)
		if err == nil {
			dest.Rating = rating
		} else {
			log.Printf("Failed to parse imdbRating: %v", err)
		}
	}

	if url, ok := response["Poster"].(string); ok {
		dest.ImageURL = url
	}

	return nil
}

func parseIDFromIMDB(url string) string {
	shortURL := strings.TrimPrefix(url, "https://www.imdb.com/")

	parts := strings.Split(shortURL, "/")

	if len(parts) > 0 {
		return parts[1]
	}

	return ""
}
