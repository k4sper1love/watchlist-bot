package watchlist

import (
	"encoding/json"
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-api/pkg/tokens"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"io"
	"log/slog"
	"net/http"
	"time"
)

const (
	verificationTokenExpiration = 10 * time.Second
)

func Register(app models.App, session *models.Session) error {
	token, err := tokens.GenerateToken(app.Vars.Secret, session.TelegramID, verificationTokenExpiration)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Verification": token,
	}

	resp, err := SendRequest(app.Vars.Host+"/api/v1/auth/register/telegram", http.MethodPost, "", headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("registration failed: %s", resp.Status)
	}

	if err := parseAuth(session, resp.Body); err != nil {
		return err
	}

	return nil
}

func Login(app models.App, session *models.Session) error {
	token, err := tokens.GenerateToken(app.Vars.Secret, session.TelegramID, verificationTokenExpiration)
	if err != nil {
		return err
	}

	headers := map[string]string{
		"Verification": token,
	}

	resp, err := SendRequest(app.Vars.Host+"/api/v1/auth/login/telegram", http.MethodPost, "", headers)
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

	if err := parseAuth(session, resp.Body); err != nil {
		return err
	}

	return nil
}

func IsTokenValid(app models.App, token string) bool {
	headers := map[string]string{
		"Authorization": token,
	}

	resp, err := SendRequest(app.Vars.Host+"/api/v1/auth/check-token", http.MethodGet, nil, headers)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func RefreshAccessToken(app models.App, session *models.Session) error {
	headers := map[string]string{
		"Authorization": session.RefreshToken,
	}

	resp, err := SendRequest(app.Vars.Host+"/api/v1/auth/refresh", http.MethodPost, nil, headers)
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

func Logout(app models.App, session *models.Session) error {
	headers := map[string]string{
		"Authorization": session.RefreshToken,
	}

	resp, err := SendRequest(app.Vars.Host+"/api/v1/auth/logout", http.MethodPost, nil, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("logout failed: %s", resp.Status)
	}

	session.ClearUser()

	return nil
}
