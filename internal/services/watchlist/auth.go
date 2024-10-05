package watchlist

import (
	"encoding/json"
	"fmt"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"net/http"
)

func RegisterUser(app *config.App, session *models.Session) error {
	resp, err := SendRequest(app.Vars.BaseURL, "/auth/register", http.MethodPost, "", &session.AuthState)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("registration failed: %s", resp.Status)
	}

	var responseMap map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseMap); err != nil {
		return err
	}

	session.AccessToken, session.RefreshToken = parseTokens(responseMap)

	return nil
}

func parseTokens(data map[string]interface{}) (string, string) {
	userData := data["user"].(map[string]interface{})

	accessToken := userData["access_token"].(string)
	refreshToken := userData["refresh_token"].(string)

	return accessToken, refreshToken
}
