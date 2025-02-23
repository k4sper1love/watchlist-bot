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

func GetCollectionFilms(app models.App, session *models.Session) (*models.CollectionFilmsResponse, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodGet,
			URL:                buildGetCollectionFilmsURL(app, session),
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	var collectionFilmsResponse models.CollectionFilmsResponse
	if err = json.NewDecoder(resp.Body).Decode(&collectionFilmsResponse); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &collectionFilmsResponse, nil
}

func CreateCollectionFilm(app models.App, session *models.Session) (*apiModels.CollectionFilm, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodPost,
			URL:                fmt.Sprintf("%s/api/v1/collections/%d/films", app.Config.APIHost, session.CollectionDetailState.Collection.ID),
			Body:               session.FilmDetailState,
			ExpectedStatusCode: http.StatusCreated,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	collectionFilm := &apiModels.CollectionFilm{}
	if err = parseCollectionFilm(collectionFilm, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return collectionFilm, nil
}

func AddCollectionFilm(app models.App, session *models.Session) (*apiModels.CollectionFilm, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodPost,
			URL:                fmt.Sprintf("%s/api/v1/collections/%d/films/%d", app.Config.APIHost, session.CollectionDetailState.Collection.ID, session.FilmDetailState.Film.ID),
			ExpectedStatusCode: http.StatusCreated,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	collectionFilm := &apiModels.CollectionFilm{}
	if err = parseCollectionFilm(collectionFilm, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return collectionFilm, nil
}

func DeleteCollectionFilm(app models.App, session *models.Session) error {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodDelete,
			URL:                fmt.Sprintf("%s/api/v1/collections/%d/films/%d", app.Config.APIHost, session.CollectionDetailState.Collection.ID, session.FilmDetailState.Film.ID),
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return err
	}
	defer utils.CloseBody(resp.Body)

	return nil
}

func buildGetCollectionFilmsURL(app models.App, session *models.Session) string {
	baseURL := fmt.Sprintf("%s/api/v1/collections/%d/films", app.Config.APIHost, session.CollectionDetailState.ObjectID)
	state := session.FilmsState
	queryParams := url.Values{}

	queryParams = addFilmsBasicParams(queryParams, state.Title, state.CurrentPage, state.PageSize)
	queryParams = addFilmsFilterAndSortingParams(queryParams, state.CollectionFilters, state.CollectionSorting)

	return fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())
}
