package watchlist

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"net/http"
)

func GetUser(app models.App, session *models.Session) (*apiModels.User, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodGet,
			URL:                fmt.Sprintf("%s/api/v1/user", app.Vars.Host),
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	user := &apiModels.User{}
	if err = parseUser(user, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return user, nil
}

func UpdateUser(app models.App, session *models.Session) (*apiModels.User, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodPut,
			URL:                fmt.Sprintf("%s/api/v1/user", app.Vars.Host),
			Body:               session.ProfileState,
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	user := &apiModels.User{}
	if err = parseUser(user, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return user, nil
}

func DeleteUser(app models.App, session *models.Session) error {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodDelete,
			URL:                fmt.Sprintf("%s/api/v1/user", app.Vars.Host),
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return err
	}
	defer utils.CloseBody(resp.Body)

	return nil
}
