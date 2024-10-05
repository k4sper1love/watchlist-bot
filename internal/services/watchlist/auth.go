package watchlist

import (
	"encoding/json"
	"fmt"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"io"
	"net/http"
)

func Register(app config.App, session *models.Session) error {
	resp, err := SendRequest(app.Vars.BaseURL, "/auth/register", http.MethodPost, "", &session.AuthState)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("registration failed: %s", resp.Status)
	}

	return parseTokens(resp.Body, session)
}

func Login(app config.App, session *models.Session) error {
	resp, err := SendRequest(app.Vars.BaseURL, "/auth/login", http.MethodPost, "", &session.AuthState)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("login failed: %s", resp.Status)
	}

	return parseTokens(resp.Body, session)
}

func IsTokenValid(app config.App, token string) bool {
	resp, err := SendRequest(app.Vars.BaseURL, "/auth/check-token", http.MethodGet, token, nil)
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusOK
}

func RefreshAccessToken(app config.App, session *models.Session) error {
	resp, err := SendRequest(app.Vars.BaseURL, "/auth/refresh", http.MethodPost, session.RefreshToken, nil)
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
	resp, err := SendRequest(app.Vars.BaseURL, "/auth/logout", http.MethodPost, session.RefreshToken, nil)
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

func parseTokens(data io.Reader, session *models.Session) error {
	var responseMap map[string]interface{}
	if err := json.NewDecoder(data).Decode(&responseMap); err != nil {
		return err
	}

	userData := responseMap["user"].(map[string]interface{})

	session.AccessToken = userData["access_token"].(string)
	session.RefreshToken = userData["refresh_token"].(string)

	return nil
}
