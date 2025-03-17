package watchlist

import (
	"encoding/json"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/tokens"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/security"
	"log/slog"
	"net/http"
	"time"
)

const (
	verificationTokenExpiration = 10 * time.Second // Duration for which the verification token is valid.
)

// Register sends a registration request to the API for a Telegram user.
func Register(app models.App, session *models.Session) error {
	return sendAuthRequest(app, session, "/api/v1/auth/register/telegram", http.StatusCreated)
}

// Login sends a login request to the API for a Telegram user.
func Login(app models.App, session *models.Session) error {
	return sendAuthRequest(app, session, "/api/v1/auth/login/telegram", http.StatusOK)
}

// sendAuthRequest is a helper function to send authentication requests (login or register) to the API.
// It generates a verification token, sends it in the request headers, and handles the response.
func sendAuthRequest(app models.App, session *models.Session, endpoint string, expectedStatusCode int) error {
	// Generate a short-lived verification token for the user.
	token, err := tokens.GenerateToken(app.Config.APISecret, session.TelegramID, verificationTokenExpiration)
	if err != nil {
		sl.Log.Error("failed to generate verification token", slog.Any("error", err), slog.Int("telegram_id", session.TelegramID))
		return err
	}

	// Send the authentication request with the verification token.
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderVerification,
			HeaderValue:        token,
			Method:             http.MethodPost, // HTTP POST method for creating data.
			URL:                app.Config.APIHost + endpoint,
			ExpectedStatusCode: expectedStatusCode,
		},
	)
	if err != nil {
		return err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	// Parse the authentication response and populate the session.
	if err = parseAuth(session, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return err
	}

	return nil
}

// IsTokenValid checks if the provided access token is valid by sending a verification request to the API.
func IsTokenValid(app models.App, encryptedToken string) bool {
	// Decrypt the encrypted token for use in the request.
	token, err := security.Decrypt(encryptedToken)
	if err != nil {
		utils.LogDecryptError(err)
		return false
	}

	// Send a request to the API to verify the token.
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodGet, // HTTP GET method for fetching data.
			URL:                app.Config.APIHost + "/api/v1/auth/check",
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
			WithoutLog:         true,          // Suppress logging for this request.
		},
	)
	if err != nil {
		return false
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	return true
}

// RefreshAccessToken refreshes the access token using the refresh token stored in the session.
func RefreshAccessToken(app models.App, session *models.Session) error {
	// Decrypt the refresh token for use in the request.
	token, err := security.Decrypt(session.RefreshToken)
	if err != nil {
		utils.LogDecryptError(err)
		return err
	}

	// Send a request to the API to refresh the access token.
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodPost, // HTTP POST method for creating data.
			URL:                app.Config.APIHost + "/api/v1/auth/refresh",
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
		},
	)
	if err != nil {
		return err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	// Parse the response to extract the new access token.
	var responseMap map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return err
	}

	// Encrypt the new access token and update the session.
	encryptedAccessToken, err := security.Encrypt(responseMap["access_token"].(string))
	if err != nil {
		utils.LogDecryptError(err)
		return err
	}

	session.AccessToken = encryptedAccessToken
	return nil
}

// Logout logs out the user by invalidating their refresh token and clearing their session data.
func Logout(app models.App, session *models.Session) error {
	// Decrypt the refresh token for use in the request.
	token, err := security.Decrypt(session.RefreshToken)
	if err != nil {
		utils.LogDecryptError(err)
		return err
	}

	// Send a request to the API to log out the user.
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodPost, // HTTP POST method for creating data.
			URL:                app.Config.APIHost + "/api/v1/auth/logout",
			ExpectedStatusCode: http.StatusOK, // Expecting a 200 OK response.
		},
	)
	if err != nil {
		return err
	}
	defer utils.CloseBody(resp.Body) // Ensure the response body is closed after use.

	return nil
}
