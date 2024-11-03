package watchlist

import (
	"encoding/json"
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"io"
	"net/http"
)

func GetUser(app models.App, session *models.Session) (*apiModels.User, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	resp, err := SendRequest(app.Vars.BaseURL+"/user", http.MethodGet, nil, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get_user failed: %s", resp.Status)
	}

	user := &apiModels.User{}
	if err := parseUser(user, resp.Body); err != nil {
		return nil, err
	}

	return user, nil
}

func parseUser(dest *apiModels.User, data io.Reader) error {
	return json.NewDecoder(data).Decode(&struct {
		User *apiModels.User `json:"user"`
	}{
		User: dest,
	})
}
