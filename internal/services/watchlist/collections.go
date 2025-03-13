package watchlist

import (
	"encoding/json"
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"net/http"
	"net/url"
)

func GetCollections(app models.App, session *models.Session) (*models.CollectionsResponse, error) {
	return getCollectionsRequest(app, session, -1, -1, session.CollectionsState.CurrentPage, session.CollectionsState.PageSize)
}

func GetCollectionsExcludeFilm(app models.App, session *models.Session) (*models.CollectionsResponse, error) {
	return getCollectionsRequest(app, session, -1, session.FilmDetailState.Film.ID, session.CollectionFilmsState.CurrentPage, session.CollectionFilmsState.PageSize)
}

func getCollectionsRequest(app models.App, session *models.Session, filmID, excludeFilmID, currentPage, pageSize int) (*models.CollectionsResponse, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodGet,
			URL:                buildGetCollectionsURL(app, session, filmID, excludeFilmID, currentPage, pageSize),
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	var collectionsResponse models.CollectionsResponse
	if err = json.NewDecoder(resp.Body).Decode(&collectionsResponse); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &collectionsResponse, nil
}

func CreateCollection(app models.App, session *models.Session) (*apiModels.Collection, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodPost,
			URL:                app.Config.APIHost + "/api/v1/collections",
			Body:               session.CollectionDetailState,
			ExpectedStatusCode: http.StatusCreated,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	collection := &apiModels.Collection{}
	if err = parseCollection(collection, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return collection, nil
}

func UpdateCollection(app models.App, session *models.Session) (*apiModels.Collection, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodPut,
			URL:                fmt.Sprintf("%s/api/v1/collections/%d", app.Config.APIHost, session.CollectionDetailState.Collection.ID),
			Body:               session.CollectionDetailState,
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	collection := &apiModels.Collection{}
	if err = parseCollection(collection, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return collection, nil
}

func DeleteCollection(app models.App, session *models.Session) error {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodDelete,
			URL:                fmt.Sprintf("%s/api/v1/collections/%d", app.Config.APIHost, session.CollectionDetailState.Collection.ID),
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return err
	}
	defer utils.CloseBody(resp.Body)

	return nil
}

func buildGetCollectionsURL(app models.App, session *models.Session, filmID, excludeFilmID, currentPage, pageSize int) string {
	baseURL := fmt.Sprintf("%s/api/v1/collections", app.Config.APIHost)
	queryParams := url.Values{}

	if filmID >= 0 {
		queryParams.Add("film", fmt.Sprintf("%d", filmID))
	}
	if excludeFilmID >= 0 {
		queryParams.Add("exclude_film", fmt.Sprintf("%d", excludeFilmID))
	}
	if currentPage > 0 {
		queryParams.Add("page", fmt.Sprintf("%d", currentPage))
	}
	if pageSize > 0 {
		queryParams.Add("page_size", fmt.Sprintf("%d", pageSize))
	}
	if session.CollectionsState.Name != "" {
		queryParams.Add("name", session.CollectionsState.Name)
	}
	if session.CollectionsState.Sorting.Sort != "" {
		queryParams.Add("sort", session.CollectionsState.Sorting.Sort)
	}

	return fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())
}
