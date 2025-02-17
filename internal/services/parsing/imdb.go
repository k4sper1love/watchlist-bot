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

	dest.Title = client.GetStringFromMap(response, "Title", "Unknown")
	dest.Year = client.GetIntFromStringMap(response, "Year", 0)
	dest.Genre = getFirstGenreFromString(response, "Genre", "")
	dest.Description = client.GetStringFromMap(response, "Plot", "")
	dest.Rating = client.GetFloatFromStringMap(response, "imdbRating", 0.0)
	dest.ImageURL = client.GetStringFromMap(response, "Poster", "")

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

func getFirstGenreFromString(data map[string]interface{}, key string, defaultValue string) string {
	if value, ok := data[key].(string); ok {
		genres := strings.Split(value, ",")
		if len(genres) > 0 {
			return strings.TrimSpace(genres[0])
		}
	}
	return defaultValue