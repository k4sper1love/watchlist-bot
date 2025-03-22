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

// GetCollections fetches the list of collections from the API.
// It supports pagination and filtering based on the session's state.
func GetCollections(app models.App, session *models.Session) (*models.CollectionsResponse, error) {
	return getCollectionsRequest(app, session, -1, -1, session.CollectionsState.CurrentPage, session.CollectionsState.PageSize)
}

// GetCollectionsExcludeFilm fetches the list of collections excluding a specific film.
// It supports pagination and filtering based on the session's state.
func GetCollectionsExcludeFilm(app models.App, session *models.Session) (*models.CollectionsResponse, error) {
	return getCollectionsRequest(app, session, -1, session.FilmDetailState.Film.ID, session.CollectionFilmsState.CurrentPage, session.CollectionFilmsState.PageSize)
}

// getCollectionsRequest is a helper function to send requests for fetching collections.
// It constructs the URL with query parameters for filtering, sorting, and pagination,
// decrypts the access token, and parses the response into a `models.CollectionsResponse` object.
func getCollectionsRequest(app models.App, session *models.Session, filmID, excludeFilmID, currentPage, pageSize int) (*models.CollectionsResponse, error) {
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(session.TelegramID, err)
		return nil, err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodGet, // HTTP GET method for fetching data.
			URL:                buildGetCollectionsURL(app, session, filmID, excludeFilmID, currentPage, pageSize),
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
			TelegramID:         session.TelegramID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	var collectionsResponse models.CollectionsResponse
	if err = json.NewDecoder(resp.Body).Decode(&collectionsResponse); err != nil {
		utils.LogParseJSONError(session.TelegramID, err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &collectionsResponse, nil
}

// CreateCollection creates a new collection by sending a POST request to the API.
// It decrypts the access token, sends the request with the collection details in the body,
// and parses the response into an `models.Collection` object.
func CreateCollection(app models.App, session *models.Session) (*apiModels.Collection, error) {
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(session.TelegramID, err)
		return nil, err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodPost, // HTTP POST method for creating data.
			URL:                app.Config.APIHost + "/api/v1/collections",
			Body:               session.CollectionDetailState,
			ExpectedStatusCode: http.StatusCreated, // Expecting a 201 Created response.
			TelegramID:         session.TelegramID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	collection := &apiModels.Collection{}
	if err = parseCollection(collection, resp.Body); err != nil {
		utils.LogParseJSONError(session.TelegramID, err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return collection, nil
}

// UpdateCollection updates an existing collection by sending a PUT request to the API.
// It decrypts the access token, sends the request with the updated collection details in the body,
// and parses the response into an `models.Collection` object.
func UpdateCollection(app models.App, session *models.Session) (*apiModels.Collection, error) {
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(session.TelegramID, err)
		return nil, err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodPut, // HTTP PUT method for updating data.
			URL:                fmt.Sprintf("%s/api/v1/collections/%d", app.Config.APIHost, session.CollectionDetailState.Collection.ID),
			Body:               session.CollectionDetailState,
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
			TelegramID:         session.TelegramID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	collection := &apiModels.Collection{}
	if err = parseCollection(collection, resp.Body); err != nil {
		utils.LogParseJSONError(session.TelegramID, err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return collection, nil
}

// DeleteCollection deletes a collection by sending a DELETE request to the API.
// It decrypts the access token, sends the request with the collection ID in the URL,
// and handles the response.
func DeleteCollection(app models.App, session *models.Session) error {
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(session.TelegramID, err)
		return err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodDelete, // HTTP DELETE method for removing data.
			URL:                fmt.Sprintf("%s/api/v1/collections/%d", app.Config.APIHost, session.CollectionDetailState.Collection.ID),
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
			TelegramID:         session.TelegramID,
		},
	)
	if err != nil {
		return err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	return nil
}

// buildGetCollectionsURL constructs the URL for fetching collections.
// It includes query parameters for filtering, sorting, and pagination.
func buildGetCollectionsURL(app models.App, session *models.Session, filmID, excludeFilmID, currentPage, pageSize int) string {
	baseURL := fmt.Sprintf("%s/api/v1/collections", app.Config.APIHost)
	queryParams := url.Values{}

	// Add optional query parameters if they are provided.
	if filmID >= 0 {
		queryParams.Add("film", fmt.Sprintf("%d", filmID))
	}
	if excludeFilmID >= 0 {
		queryParams.Add("exclude_film", fmt.Sprintf("%d", excludeFilmID))
	}
	if currentPage > 0 {
		queryParams.Add("page", fmt.Sprintf("%d", currentPage))
	}
	if pageSize > 0 {
		queryParams.Add("page_size", fmt.Sprintf("%d", pageSize))
	}
	if session.CollectionsState.Name != "" {
		queryParams.Add("name", session.CollectionsState.Name)
	}
	if session.CollectionsState.Sorting.Sort != "" {
		queryParams.Add("sort", session.CollectionsState.Sorting.Sort)
	}

	// Encode the query parameters and append them to the base URL.
	return fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())
}
