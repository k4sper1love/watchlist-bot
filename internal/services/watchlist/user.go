package watchlist

import (
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/security"
	"net/http"
)

// GetUser fetches the current user's details from the API.
// It decrypts the access token, sends a GET request to the API, and parses the response into an `models.User` object.
func GetUser(app models.App, session *models.Session) (*apiModels.User, error) {
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
			URL:                app.Config.APIHost + "/api/v1/user",
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
			TelegramID:         session.TelegramID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	user := &apiModels.User{}
	if err = parseUser(user, resp.Body); err != nil {
		utils.LogParseJSONError(session.TelegramID, err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return user, nil
}

// UpdateUser updates the current user's details by sending a PUT request to the API.
// It decrypts the access token, sends the request with updated user details in the body,
// and parses the response into an `models.User` object.
func UpdateUser(app models.App, session *models.Session) (*apiModels.User, error) {
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
			URL:                app.Config.APIHost + "/api/v1/user",
			Body:               session.ProfileState,
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
			TelegramID:         session.TelegramID,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	user := &apiModels.User{}
	if err = parseUser(user, resp.Body); err != nil {
		utils.LogParseJSONError(session.TelegramID, err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return user, nil
}

// DeleteUser deletes the current user's account by sending a DELETE request to the API.
// It decrypts the access token, sends the request, and handles the response.
func DeleteUser(app models.App, session *models.Session) error {
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
			URL:                app.Config.APIHost + "/api/v1/user",
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
