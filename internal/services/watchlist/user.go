package watchlist

import (
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/security"
	"net/http"
)

func GetUser(app models.App, session *models.Session) (*apiModels.User, error) {
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(err)
		return nil, err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodGet,
			URL:                app.Config.APIHost + "/api/v1/user",
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
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(err)
		return nil, err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodPut,
			URL:                app.Config.APIHost + "/api/v1/user",
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
	token, err := security.Decrypt(session.AccessToken)
	if err != nil {
		utils.LogDecryptError(err)
		return err
	}

	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        token,
			Method:             http.MethodDelete,
			URL:                app.Config.APIHost + "/api/v1/user",
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return err
	}
	defer utils.CloseBody(resp.Body)

	return nil
}
