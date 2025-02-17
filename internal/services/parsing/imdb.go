package parsing

import (
	"encoding/json"
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"io"
	"log/slog"
	"net/http"
	"strconv"
	"strings"
)

func GetFilmFromIMDB(app models.App, url string) (*apiModels.Film, error) {
	id, err := parseIDFromIMDB(url)
	if err != nil {
		sl.Log.Error("failed to parse IMDB id", slog.Any("error", err), slog.String("url", url))
		return nil, err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			Method:             http.MethodGet,
			URL:                fmt.Sprintf("http://www.omdbapi.com/?apikey=%s&i=%s&plot=full", app.Vars.IMDBAPIToken, id),
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	var film apiModels.Film
	if err = parseFilmFromIMDB(&film, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
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
	} else {
		return fmt.Errorf("failed to parse title")
	}

	if yearStr, ok := response["Year"].(string); ok {
		year, err := strconv.Atoi(yearStr)
		if err != nil {
			return fmt.Errorf("failed to parse year: %v", err)
		}
		dest.Year = year
	} else {
		return fmt.Errorf("failed to parse year")
	}

	if genres, ok := response["Genre"].(string); ok {
		dest.Genre = strings.Split(genres, ",")[0]
	} else {
		return fmt.Errorf("failed to parse genre")
	}

	if description, ok := response["Plot"].(string); ok {
		dest.Description = description
	} else {
		return fmt.Errorf("failed to parse description")
	}

	if ratingStr, ok := response["imdbRating"].(string); ok {
		rating, err := strconv.ParseFloat(ratingStr, 64)
		if err != nil {
			return fmt.Errorf("failed to parse imdbRating: %v", err)
		}
		dest.Rating = rating
	} else {
		return fmt.Errorf("failed to parse imdbRating")
	}

	if url, ok := response["Poster"].(string); ok {
		dest.ImageURL = url
	} else {
		return fmt.Errorf("failed to parse imageURL")
	}

	return nil
}

func parseIDFromIMDB(url string) (string, error) {
	shortURL := strings.TrimPrefix(url, "https://www.imdb.com/")

	parts := strings.Split(shortURL, "/")

	if len(parts) > 0 {
		return parts[1], nil
	}

	return "", fmt.Errorf("id not found")
}
