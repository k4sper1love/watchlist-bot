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

// GetFilmFromIMDB fetches film details from the IMDB API using the provided URL.
// It extracts the film ID from the URL, makes an HTTP request to the OMDB API,
// and parses the response into an `models.Film` object.
func GetFilmFromIMDB(app models.App, url string) (*apiModels.Film, error) {
	id, err := parseIDFromIMDB(url)
	if err != nil {
		sl.Log.Error("failed to parse IMDB id", slog.Any("error", err), slog.String("url", url))
		return nil, err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			Method:             http.MethodGet, // HTTP GET method for fetching data.
			URL:                fmt.Sprintf("http://www.omdbapi.com/?apikey=%s&i=%s&plot=full", app.Config.IMDBAPIToken, id),
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	var film apiModels.Film
	if err = parseFilmFromIMDB(&film, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &film, nil
}

// parseFilmFromIMDB parses film details from the IMDB API response into an `models.Film` object.
func parseFilmFromIMDB(dest *apiModels.Film, data io.Reader) error {
	var response map[string]interface{} // Temporary map to hold the JSON response.
	if err := json.NewDecoder(data).Decode(&response); err != nil {
		return err
	}

	// Extract film details from the response.
	dest.Title = getStringFromMap(response, "Title", "Unknown")      // Default title is "Unknown".
	dest.Year = getIntFromStringMap(response, "Year", 0)             // Default year is 0.
	dest.Genre = getFirstGenreFromString(response, "Genre", "")      // Extract the first genre from the list.
	dest.Description = getStringFromMap(response, "Plot", "")        // Default description is an empty string.
	dest.Rating = getFloatFromStringMap(response, "imdbRating", 0.0) // Default rating is 0.0.
	dest.ImageURL = getStringFromMap(response, "Poster", "")         // Default image URL is an empty string.

	return nil
}

// parseIDFromIMDB extracts the IMDB ID from the given URL.
func parseIDFromIMDB(url string) (string, error) {
	shortURL := strings.TrimPrefix(url, "https://www.imdb.com/") // Remove the base URL.
	parts := strings.Split(shortURL, "/")                        // Split the remaining path by "/".
	if len(parts) > 1 {
		return parts[1], nil // The second part of the path is the IMDB ID.
	}
	return "", fmt.Errorf("id not found") // Return an error if the ID cannot be extracted.
}

// getFirstGenreFromString extracts the first genre from a comma-separated string in the map.
// If the key exists and contains genres, it splits the string and returns the first genre.
// Otherwise, it returns the default value.
func getFirstGenreFromString(data map[string]interface{}, key string, defaultValue string) string {
	if value, ok := data[key].(string); ok {
		if genres := strings.Split(value, ","); len(genres) > 0 {
			return strings.TrimSpace(genres[0]) // Return the first genre after trimming whitespace.
		}
	}
	return defaultValue // Return the default value if no genres are found.
}
