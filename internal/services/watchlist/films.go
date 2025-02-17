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

func GetFilms(app models.App, session *models.Session) (*models.FilmsResponse, error) {
	return getFilmsRequest(app, session, -1, session.FilmsState.CurrentPage, session.FilmsState.PageSize)
}

func GetFilmsExcludeCollection(app models.App, session *models.Session) (*models.FilmsResponse, error) {
	return getFilmsRequest(app, session, session.CollectionDetailState.Collection.ID, session.CollectionFilmsState.CurrentPage, session.CollectionFilmsState.PageSize)
}

//func GetFilmsByTitle(app models.App, session *models.Session) (*models.FilmsResponse, error) {
//	return getFilmsRequest(app, session, -1, session.FilmsState.CurrentPage, session.FilmsState.PageSize, session.FilmsState.Title)
//}

func getFilmsRequest(app models.App, session *models.Session, collectionID, currentPage, pageSize int) (*models.FilmsResponse, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodGet,
			URL:                buildGetFilmsURL(app, session, collectionID, currentPage, pageSize),
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	var filmsResponse models.FilmsResponse
	if err = json.NewDecoder(resp.Body).Decode(&filmsResponse); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return &filmsResponse, nil
}

func GetFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodGet,
			URL:                fmt.Sprintf("%s/api/v1/films/%d", app.Vars.Host, session.FilmDetailState.Film.ID),
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	film := &apiModels.Film{}
	if err = parseFilm(film, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return film, nil
}

func UpdateFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodPut,
			URL:                fmt.Sprintf("%s/api/v1/films/%d", app.Vars.Host, session.FilmDetailState.Film.ID),
			Body:               session.FilmDetailState,
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	film := &apiModels.Film{}
	if err = parseFilm(film, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return film, nil
}

func CreateFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodPost,
			URL:                fmt.Sprintf("%s/api/v1/films", app.Vars.Host),
			Body:               session.FilmDetailState,
			ExpectedStatusCode: http.StatusCreated,
		},
	)
	if err != nil {
		return nil, err
	}
	defer utils.CloseBody(resp.Body)

	film := &apiModels.Film{}
	if err = parseFilm(film, resp.Body); err != nil {
		utils.LogParseJSONError(err, resp.Request.Method, resp.Request.URL.String())
		return nil, err
	}

	return film, nil
}

func DeleteFilm(app models.App, session *models.Session) error {
	resp, err := client.Do(
		&client.CustomRequest{
			HeaderType:         client.HeaderAuthorization,
			HeaderValue:        session.AccessToken,
			Method:             http.MethodDelete,
			URL:                fmt.Sprintf("%s/api/v1/films/%d", app.Vars.Host, session.FilmDetailState.Film.ID),
			ExpectedStatusCode: http.StatusOK,
		},
	)
	if err != nil {
		return err
	}
	defer utils.CloseBody(resp.Body)

	return nil
}

func buildGetFilmsURL(app models.App, session *models.Session, collectionID, currentPage, pageSize int) string {
	baseURL := fmt.Sprintf("%s/api/v1/films", app.Vars.Host)
	queryParams := url.Values{}

	if collectionID >= 0 {
		queryParams.Add("exclude_collection", fmt.Sprintf("%d", collectionID))
	}

	state := session.FilmsState

	queryParams = addFilmsBasicParams(queryParams, state.Title, currentPage, pageSize)

	queryParams = addFilmsFilterAndSortingParams(queryParams, state.FilmFilters, state.FilmSorting)

	requestURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	return requestURL
}

func addFilmsBasicParams(queryParams url.Values, title string, currentPage, pageSize int) url.Values {
	if title != "" {
		queryParams.Add("title", title)
	}

	if currentPage > 0 {
		queryParams.Add("page", fmt.Sprintf("%d", currentPage))
	}

	if pageSize > 0 {
		queryParams.Add("page_size", fmt.Sprintf("%d", pageSize))
	}

	return queryParams
}

func addFilmsFilterAndSortingParams(queryParams url.Values, filter *models.FiltersFilm, sorting *models.Sorting) url.Values {
	if filter.Rating != "" {
		queryParams.Add("rating", filter.Rating)
	}

	if filter.UserRating != "" {
		queryParams.Add("user_rating", filter.UserRating)
	}

	if filter.Year != "" {
		queryParams.Add("year", filter.Year)
	}

	if filter.IsViewed != nil {
		queryParams.Add("is_viewed", fmt.Sprintf("%t", *filter.IsViewed))
	}

	if filter.IsFavorite != nil {
		queryParams.Add("is_favorite", fmt.Sprintf("%t", *filter.IsFavorite))
	}

	if filter.HasURL != nil {
		queryParams.Add("has_url", fmt.Sprintf("%t", *filter.HasURL))
	}

	if sorting.Sort != "" {
		queryParams.Add("sort", sorting.Sort)
	}

	return queryParams
}
