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

// GetCollectionFilms fetches the list of films in a collection from the API.
// It decrypts the access token, sends a GET request with query parameters for filtering and pagination,
// and parses the response into a `models.CollectionFilmsResponse` object.
func GetCollectionFilms(app models.App, session *models.Session) (*models.CollectionFilmsResponse, error) {
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(session.TelegramID, err)
		return nil, err
	}

	// Build the URL with query parameters for filtering and pagination.
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodGet, // HTTP GET method for fetching data.
			URL:                buildGetCollectionFilmsURL(app, session),
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
			TelegramID:         session.TelegramID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	// Parse the response into a structured format.
	var collectionFilmsResponse models.CollectionFilmsResponse
	if err = json.NewDecoder(resp.Body).Decode(&collectionFilmsResponse); err != nil {
		utils.LogParseJSONError(session.TelegramID, err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &collectionFilmsResponse, nil
}

// CreateCollectionFilm creates a new film in a collection by sending a POST request to the API.
// It decrypts the access token, sends the request with the film details in the body,
// and parses the response into an `models.CollectionFilm` object.
func CreateCollectionFilm(app models.App, session *models.Session) (*apiModels.CollectionFilm, error) {
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(session.TelegramID, err)
		return nil, err
	}

	// Send a POST request to create a new film in the collection.
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodPost, // HTTP POST method for creating data.
			URL:                fmt.Sprintf("%s/api/v1/collections/%d/films", app.Config.APIHost, session.CollectionDetailState.Collection.ID),
			Body:               session.FilmDetailState,
			ExpectedStatusCode: http.StatusCreated, // Expecting a 201 Created response.
			TelegramID:         session.TelegramID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	// Parse the response into a structured format.
	collectionFilm := &apiModels.CollectionFilm{}
	if err = parseCollectionFilm(collectionFilm, resp.Body); err != nil {
		utils.LogParseJSONError(session.TelegramID, err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return collectionFilm, nil
}

// AddCollectionFilm adds an existing film to a collection by sending a POST request to the API.
// It decrypts the access token, sends the request with the film ID and collection ID in the URL,
// and parses the response into an `models.CollectionFilm` object.
func AddCollectionFilm(app models.App, session *models.Session) (*apiModels.CollectionFilm, error) {
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(session.TelegramID, err)
		return nil, err
	}

	// Send a POST request to add an existing film to the collection.
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodPost, // HTTP POST method for adding data.
			URL:                fmt.Sprintf("%s/api/v1/collections/%d/films/%d", app.Config.APIHost, session.CollectionDetailState.Collection.ID, session.FilmDetailState.Film.ID),
			ExpectedStatusCode: http.StatusCreated, // Expecting a 201 Created response.
			TelegramID:         session.TelegramID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	// Parse the response into a structured format.
	collectionFilm := &apiModels.CollectionFilm{}
	if err = parseCollectionFilm(collectionFilm, resp.Body); err != nil {
		utils.LogParseJSONError(session.TelegramID, err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return collectionFilm, nil
}

// DeleteCollectionFilm deletes a film from a collection by sending a DELETE request to the API.
// It decrypts the access token, sends the request with the film ID and collection ID in the URL,
// and handles the response.
func DeleteCollectionFilm(app models.App, session *models.Session) error {
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(session.TelegramID, err)
		return err
	}

	// Send a DELETE request to remove the film from the collection.
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodDelete, // HTTP DELETE method for removing data.
			URL:                fmt.Sprintf("%s/api/v1/collections/%d/films/%d", app.Config.APIHost, session.CollectionDetailState.Collection.ID, session.FilmDetailState.Film.ID),
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

// buildGetCollectionFilmsURL constructs the URL for fetching films in a collection.
// It includes query parameters for filtering, sorting, and pagination.
func buildGetCollectionFilmsURL(app models.App, session *models.Session) string {
	baseURL := fmt.Sprintf("%s/api/v1/collections/%d/films", app.Config.APIHost, session.CollectionDetailState.ObjectID)
	state := session.FilmsState
	queryParams := url.Values{}

	// Add basic parameters (title, page, page size) to the query.
	queryParams = addFilmsBasicParams(queryParams, state.Title, state.CurrentPage, state.PageSize)

	// Add filter and sorting parameters to the query.
	queryParams = addFilmsFilterAndSortingParams(queryParams, state.CollectionFilters, state.CollectionSorting)

	// Encode the query parameters and append them to the base URL.
	return fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())
}
