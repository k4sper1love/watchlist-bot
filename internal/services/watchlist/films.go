package watchlist

import (
	"encoding/json"
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/security"
	"net/http"
	"net/url"
)

// GetFilms fetches a paginated list of films from the API.
// It supports filtering and pagination based on the session's state.
func GetFilms(app models.App, session *models.Session) (*models.FilmsResponse, error) {
	return getFilmsRequest(app, session, -1, session.FilmsState.CurrentPage, session.FilmsState.PageSize)
}

// GetFilmsExcludeCollection fetches a paginated list of films excluding those in a specific collection.
// It supports filtering and pagination based on the session's state.
func GetFilmsExcludeCollection(app models.App, session *models.Session) (*models.FilmsResponse, error) {
	return getFilmsRequest(app, session, session.CollectionDetailState.Collection.ID, session.CollectionFilmsState.CurrentPage, session.CollectionFilmsState.PageSize)
}

// getFilmsRequest is a helper function to send requests for fetching films.
// It constructs the URL with query parameters for filtering, sorting, and pagination,
// decrypts the access token, and parses the response into a `models.FilmsResponse` object.
func getFilmsRequest(app models.App, session *models.Session, collectionID, currentPage, pageSize int) (*models.FilmsResponse, error) {
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(err)
		return nil, err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodGet, // HTTP GET method for fetching data.
			URL:                buildGetFilmsURL(app, session, collectionID, currentPage, pageSize),
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	var filmsResponse models.FilmsResponse
	if err = json.NewDecoder(resp.Body).Decode(&filmsResponse); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &filmsResponse, nil
}

// GetFilm fetches a single film by its ID from the API.
// It decrypts the access token, sends the request, and parses the response into an `models.Film` object.
func GetFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(err)
		return nil, err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodGet, // HTTP GET method for fetching data.
			URL:                fmt.Sprintf("%s/api/v1/films/%d", app.Config.APIHost, session.FilmDetailState.Film.ID),
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	film := &apiModels.Film{}
	if err = parseFilm(film, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return film, nil
}

// UpdateFilm updates an existing film by sending a PUT request to the API.
// It decrypts the access token, sends the request with updated film details in the body,
// and parses the response into an `models.Film` object.
func UpdateFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(err)
		return nil, err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodPut, // HTTP PUT method for updating data.
			URL:                fmt.Sprintf("%s/api/v1/films/%d", app.Config.APIHost, session.FilmDetailState.Film.ID),
			Body:               session.FilmDetailState,
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	film := &apiModels.Film{}
	if err = parseFilm(film, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return film, nil
}

// CreateFilm creates a new film by sending a POST request to the API.
// It decrypts the access token, sends the request with the film details in the body,
// and parses the response into an `apiModels.Film` object.
func CreateFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(err)
		return nil, err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodPost, // HTTP POST method for creating data.
			URL:                app.Config.APIHost + "/api/v1/films",
			Body:               session.FilmDetailState,
			ExpectedStatusCode: http.StatusCreated, // Expecting a 201 Created response.
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	film := &apiModels.Film{}
	if err = parseFilm(film, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return film, nil
}

// DeleteFilm deletes a film by sending a DELETE request to the API.
// It decrypts the access token, sends the request with the film ID in the URL,
// and handles the response.
func DeleteFilm(app models.App, session *models.Session) error {
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(err)
		return err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodDelete, // HTTP DELETE method for removing data.
			URL:                fmt.Sprintf("%s/api/v1/films/%d", app.Config.APIHost, session.FilmDetailState.Film.ID),
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
		},
	)
	if err != nil {
		return err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	return nil
}

// buildGetFilmsURL constructs the URL for fetching films.
// It includes query parameters for filtering, sorting, and pagination.
func buildGetFilmsURL(app models.App, session *models.Session, collectionID, currentPage, pageSize int) string {
	baseURL := fmt.Sprintf("%s/api/v1/films", app.Config.APIHost)
	state := session.FilmsState
	queryParams := url.Values{}

	// Add optional query parameters if they are provided.
	if collectionID >= 0 {
		queryParams.Add("exclude_collection", fmt.Sprintf("%d", collectionID))
	}

	queryParams = addFilmsBasicParams(queryParams, state.Title, currentPage, pageSize)
	queryParams = addFilmsFilterAndSortingParams(queryParams, state.FilmFilters, state.FilmSorting)

	// Encode the query parameters and append them to the base URL.
	return fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())
}

// addFilmsBasicParams adds basic query parameters (title, page, page size) to the URL.
func addFilmsBasicParams(queryParams url.Values, title string, currentPage, pageSize int) url.Values {
	if title != "" {
		queryParams.Add("title", title)
	}
	if currentPage > 0 {
		queryParams.Add("page", fmt.Sprintf("%d", currentPage))
	}
	if pageSize > 0 {
		queryParams.Add("page_size", fmt.Sprintf("%d", pageSize))
	}

	return queryParams
}

// addFilmsFilterAndSortingParams adds filter and sorting query parameters to the URL.
func addFilmsFilterAndSortingParams(queryParams url.Values, filter *models.FilmFilters, sorting *models.Sorting) url.Values {
	if filter.Rating != "" {
		queryParams.Add("rating", filter.Rating)
	}
	if filter.UserRating != "" {
		queryParams.Add("user_rating", filter.UserRating)
	}
	if filter.Year != "" {
		queryParams.Add("year", filter.Year)
	}
	if filter.IsViewed != nil {
		queryParams.Add("is_viewed", fmt.Sprintf("%t", *filter.IsViewed))
	}
	if filter.IsFavorite != nil {
		queryParams.Add("is_favorite", fmt.Sprintf("%t", *filter.IsFavorite))
	}
	if filter.HasURL != nil {
		queryParams.Add("has_url", fmt.Sprintf("%t", *filter.HasURL))
	}
	if sorting.Sort != "" {
		queryParams.Add("sort", sorting.Sort)
	}

	return queryParams
}
