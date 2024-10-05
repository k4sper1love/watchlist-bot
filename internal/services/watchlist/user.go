package watchlist

import (
	"encoding/json"
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"net/http"
	"time"
)

func GetUser(app *config.App, session *models.Session) (*apiModels.User, error) {
	resp, err := SendRequest(app.Vars.BaseURL, "/user", http.MethodGet, session.AccessToken, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("не удалось, код: %s", resp.Status)
	}
	sl.Log.Debug("resp", resp)

	var responseData map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&responseData); err != nil {
		return nil, err
	}

	return parseUser(responseData)
}

func parseUser(data map[string]interface{}) (*apiModels.User, error) {
	userData, ok := data["user"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("данные пользователя недоступны")
	}

	id := int(userData["id"].(float64))
	username := userData["username"].(string)
	email := userData["email"].(string)
	createdAtStr := userData["created_at"].(string)

	createdAt, err := time.Parse(time.RFC3339, createdAtStr)
	if err != nil {
		return nil, err
	}

	user := &apiModels.User{
		Id:        id,
		Username:  username,
		Email:     email,
		CreatedAt: createdAt,
	}

	return user, nil
}
