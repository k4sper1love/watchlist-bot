package watchlist

import (
	"encoding/json"
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/tokens"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const (
	verificationTokenExpiration = 10 * time.Second
)

func Register(app config.App, session *models.Session) error {
	token, err := tokens.GenerateToken(app.Vars.Secret, session.TelegramID, verificationTokenExpiration)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Verification": token,
	}

	resp, err := SendRequest(app.Vars.BaseURL, "/auth/register/telegram", http.MethodPost, "", headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("registration failed: %s", resp.Status)
	}

	return parseAuth(session, resp.Body)
}

func Login(app config.App, session *models.Session) error {
	token, err := tokens.GenerateToken(app.Vars.Secret, session.TelegramID, verificationTokenExpiration)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Verification": token,
	}

	resp, err := SendRequest(app.Vars.BaseURL, "/auth/login/telegram", http.MethodPost, "", headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			sl.Log.Error("failed to read response body", slog.Any("err", err))
			return err
		}
		sl.Log.Error("login failed", slog.Any("body", string(body)))
		return fmt.Errorf("login failed: %s", resp.Status)
	}

	return parseAuth(session, resp.Body)
}

func IsTokenValid(app config.App, token string) bool {
	headers := map[string]string{
		"Authorization": token,
	}

	resp, err := SendRequest(app.Vars.BaseURL, "/auth/check-token", http.MethodGet, nil, headers)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func RefreshAccessToken(app config.App, session *models.Session) error {
	headers := map[string]string{
		"Authorization": session.RefreshToken,
	}

	resp, err := SendRequest(app.Vars.BaseURL, "/auth/refresh", http.MethodPost, nil, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("refresh token failed: %s", resp.Status)
	}

	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return err
	}

	session.AccessToken = responseMap["access_token"].(string)

	return nil
}

func Logout(app config.App, session *models.Session) error {
	headers := map[string]string{
		"Authorization": session.RefreshToken,
	}

	resp, err := SendRequest(app.Vars.BaseURL, "/auth/logout", http.MethodPost, nil, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("logout failed: %s", resp.Status)
	}

	session.AccessToken = ""
	session.RefreshToken = ""

	return nil
}

func parseAuth(dest *models.Session, data io.Reader) error {
	var responseMap map[string]interface{}
	if err := json.NewDecoder(data).Decode(&responseMap); err != nil {
		return err
	}

	userData := responseMap["user"].(map[string]interface{})

	if id, ok := userData["id"].(float64); ok {
		dest.UserID = int(id)
	} else {
		return fmt.Errorf("invalid id format")
	}
	dest.AccessToken = userData["access_token"].(string)
	dest.RefreshToken = userData["refresh_token"].(string)

	return nil
}
