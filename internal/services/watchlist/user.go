package watchlist

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"net/http"
)

func GetUser(app models.App, session *models.Session) (*apiModels.User, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	resp, err := client.SendRequestWithOptions(app.Vars.Host+"/api/v1/user", http.MethodGet, nil, headers)
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

func UpdateUser(app models.App, session *models.Session) (*apiModels.User, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	resp, err := client.SendRequestWithOptions(app.Vars.Host+"/api/v1/user", http.MethodPut, session.ProfileState, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update_user failed: %s", resp.Status)
	}

	user := &apiModels.User{}
	if err := parseUser(user, resp.Body); err != nil {
		return nil, err
	}

	return user, nil
}

func DeleteUser(app models.App, session *models.Session) error {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	resp, err := client.SendRequestWithOptions(app.Vars.Host+"/api/v1/user", http.MethodDelete, nil, headers)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete_user failed: %s", resp.Status)
	}

	return nil
}
