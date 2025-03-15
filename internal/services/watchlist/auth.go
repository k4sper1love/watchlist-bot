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
	verificationTokenExpiration = 10 * time.Second
)

func Register(app models.App, session *models.Session) error {
	return sendAuthRequest(app, session, "/api/v1/auth/register/telegram", http.StatusCreated)
}

func Login(app models.App, session *models.Session) error {
	return sendAuthRequest(app, session, "/api/v1/auth/login/telegram", http.StatusOK)
}

func sendAuthRequest(app models.App, session *models.Session, endpoint string, expectedStatusCode int) error {
	token, err := tokens.GenerateToken(app.Config.APISecret, session.TelegramID, verificationTokenExpiration)
	if err != nil {
		sl.Log.Error("failed to generate verification token", slog.Any("error", err), slog.Int("telegram_id", session.TelegramID))
		return err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderVerification,
			HeaderValue:        token,
			Method:             http.MethodPost,
			URL:                app.Config.APIHost + endpoint,
			ExpectedStatusCode: expectedStatusCode,
		},
	)
	if err != nil {
		return err
	}
	defer utils.CloseBody(resp.Body)

	if err = parseAuth(session, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return err
	}

	return nil
}

func IsTokenValid(app models.App, encryptedToken string) bool {
	token, err := security.Decrypt(encryptedToken)
	if err != nil {
		utils.LogDecryptError(err)
		return false
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodGet,
			URL:                app.Config.APIHost + "/api/v1/auth/check",
			ExpectedStatusCode: http.StatusOK,
			WithoutLog:         true,
		},
	)
	if err != nil {
		return false
	}
	defer utils.CloseBody(resp.Body)

	return true
}

func RefreshAccessToken(app models.App, session *models.Session) error {
	token, err := security.Decrypt(session.RefreshToken)
	if err != nil {
		utils.LogDecryptError(err)
		return err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodPost,
			URL:                app.Config.APIHost + "/api/v1/auth/refresh",
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return err
	}
	defer utils.CloseBody(resp.Body)

	var responseMap map[string]interface{}
	if err = json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return err
	}

	encryptedAccessToken, err := security.Encrypt(responseMap["access_token"].(string))
	if err != nil {
		utils.LogDecryptError(err)
		return err
	}

	session.AccessToken = encryptedAccessToken
	return nil
}

func Logout(app models.App, session *models.Session) error {
	token, err := security.Decrypt(session.RefreshToken)
	if err != nil {
		utils.LogDecryptError(err)
		return err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodPost,
			URL:                app.Config.APIHost + "/api/v1/auth/logout",
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return err
	}
	defer utils.CloseBody(resp.Body)

	session.ClearUser()
	return nil
}
