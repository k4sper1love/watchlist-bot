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
	"github.com/k4sper1love/watchlist-bot/pkg/security"
	"io"
	"log/slog"
	"net/http"
	"net/url"
)

// GetFilmFromKinopoisk fetches a single film from the Kinopoisk API using the provided URL.
// It extracts the query and ID from the URL, makes an HTTP request to the Kinopoisk API,
// and parses the response into an `models.Film` object.
func GetFilmFromKinopoisk(session *models.Session, url string) (*apiModels.Film, error) {
	queryKey, id, err := utils.ExtractKinopoiskQuery(url)
	if err != nil {
		sl.Log.Error("failed to extract query", slog.Any("error", err), slog.String("url", url))
		return nil, err
	}

	// Construct the API URL using the extracted query key and ID.
	apiURL := fmt.Sprintf("https://api.kinopoisk.dev/v1.4/movie?%s=%s", queryKey, id)

	resp, err := getDataFromKinopoisk(session, apiURL)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	film, err := parseFilmFromKinopoisk(resp.Body)
	if err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return film, nil
}

// GetFilmsFromKinopoisk fetches a list of films from the Kinopoisk API based on the current session state.
// It constructs a search query using the session's title, page, and page size, and parses the response
// into a list of `models.Film` objects along with metadata.
func GetFilmsFromKinopoisk(session *models.Session) ([]apiModels.Film, *filters.Metadata, error) {
	state := session.FilmsState
	query := url.QueryEscape(state.Title) // Escape the title to ensure it's URL-safe.

	// Construct the API URL using the extracted query key and ID.
	apiURL := fmt.Sprintf("https://api.kinopoisk.dev/v1.4/movie/search?page=%d&limit=%d&query=%s", state.CurrentPage, state.PageSize, query)

	resp, err := getDataFromKinopoisk(session, apiURL)
	if err != nil {
		return nil, nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	films, metadata, err := parseFilmsFromKinopoisk(resp.Body)
	if err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, nil, err
	}

	return films, metadata, nil
}

// getDataFromKinopoisk sends an HTTP GET request to the Kinopoisk API with the required API key.
// The API key is decrypted from the session's KinopoiskAPIToken before being used.
func getDataFromKinopoisk(session *models.Session, url string) (*http.Response, error) {
	token, err := security.Decrypt(session.KinopoiskAPIToken)
	if err != nil {
		utils.LogDecryptError(err)
		return nil, err
	}

	return client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderExternalAPIKey, // Use the external API key header.
			HeaderValue:        token,
			Method:             http.MethodGet, // HTTP GET method for fetching data.
			URL:                url,
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
		},
	)
}

// parseFilmFromKinopoisk parses a single film from the Kinopoisk API response.
func parseFilmFromKinopoisk(data io.Reader) (*apiModels.Film, error) {
	var response map[string]interface{}
	if err := json.NewDecoder(data).Decode(&response); err != nil {
		return nil, err
	}

	docs, ok := response["docs"].([]interface{})
	if !ok || len(docs) == 0 {
		return nil, fmt.Errorf("film not found in response")
	}

	return parseFilmDataKinopoisk(docs[0].(map[string]interface{})), nil
}

// parseFilmsFromKinopoisk parses a list of films and metadata from the Kinopoisk API response.
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

		film := parseFilmDataKinopoisk(filmData)
		film.URL = fmt.Sprintf("https://www.kinopoisk.ru/film/%d/", film.ID)

		films = append(films, *film)
	}

	// Extract pagination metadata from the response.
	metadata := filters.Metadata{
		TotalRecords: getIntFromMap(response, "total", 0),
		PageSize:     getIntFromMap(response, "limit", 0),
		CurrentPage:  getIntFromMap(response, "page", 0),
		LastPage:     getIntFromMap(response, "pages", 0),
	}

	return films, &metadata, nil
}

// parseFilmDataKinopoisk extracts film details from the Kinopoisk API response data.
func parseFilmDataKinopoisk(data map[string]interface{}) *apiModels.Film {
	return &apiModels.Film{
		ID:          getIntFromMap(data, "id", 0),                      // Extract the film ID.
		Title:       getStringFromMap(data, "name", "Unknown"),         // Extract the film title.
		Year:        getIntFromMap(data, "year", 0),                    // Extract the release year.
		Genre:       getFirstGenre(data, ""),                           // Extract the first genre.
		Description: getStringFromMap(data, "description", ""),         // Extract the description.
		Rating:      getFloatFromNestedMap(data, "rating", "kp", 0.0),  // Extract the Kinopoisk rating.
		ImageURL:    getStringFromNestedMap(data, "poster", "url", ""), // Extract the poster image URL.
	}
}

// getFirstGenre extracts the first genre from the "genres" field in the Kinopoisk API response.
func getFirstGenre(data map[string]interface{}, defaultValue string) string {
	if genres, ok := data["genres"].([]interface{}); ok && len(genres) > 0 {
		if firstGenre, ok := genres[0].(map[string]interface{}); ok {
			if genreName, ok := firstGenre["name"].(string); ok {
				return genreName // Return the name of the first genre.
			}
		}
	}
	return defaultValue // Return the default value if no genre is found.
}
