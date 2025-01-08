package watchlist

import (
	"encoding/json"
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
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
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := buildGetFilmsURL(app, session, collectionID, currentPage, pageSize)

	resp, err := client.SendRequestWithOptions(requestURL, http.MethodGet, nil, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get_films failed: %s", resp.Status)
	}

	var filmsResponse models.FilmsResponse
	if err := json.NewDecoder(resp.Body).Decode(&filmsResponse); err != nil {
		return nil, err
	}

	return &filmsResponse, nil
}

func UpdateFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("%s/api/v1/films/%d", app.Vars.Host, session.FilmDetailState.Film.ID)

	resp, err := client.SendRequestWithOptions(requestURL, http.MethodPut, session.FilmDetailState, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update_film failed: %d", resp.StatusCode)
	}

	film := &apiModels.Film{}
	if err := parseFilm(film, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to parse collection film: %w", err)
	}

	return film, nil
}

func CreateFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("%s/api/v1/films", app.Vars.Host)

	resp, err := client.SendRequestWithOptions(requestURL, http.MethodPost, session.FilmDetailState, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create_film failed: %d", resp.StatusCode)
	}

	film := &apiModels.Film{}
	if err := parseFilm(film, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to parse film: %w", err)
	}

	return film, nil
}

func DeleteFilm(app models.App, session *models.Session) error {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("%s/api/v1/films/%d", app.Vars.Host, session.FilmDetailState.Film.ID)

	resp, err := client.SendRequestWithOptions(requestURL, http.MethodDelete, nil, headers)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete_film failed: %s", resp.Status)
	}

	return nil
}

func buildGetFilmsURL(app models.App, session *models.Session, collectionID, currentPage, pageSize int) string {
	baseURL := fmt.Sprintf("%s/api/v1/films", app.Vars.Host)
	queryParams := url.Values{}

	if collectionID >= 0 {
		queryParams.Add("exclude_collection", fmt.Sprintf("%d", collectionID))
	}

	if currentPage > 0 {
		queryParams.Add("page", fmt.Sprintf("%d", currentPage))
	}

	if pageSize > 0 {
		queryParams.Add("page_size", fmt.Sprintf("%d", pageSize))
	}

	if session.FilmsState.Title != "" {
		queryParams.Add("title", session.FilmsState.Title)
	}

	if session.FilmsState.FilmFilters.MinRating > 0 {
		queryParams.Add("rating_min", fmt.Sprintf("%.2f", session.FilmsState.FilmFilters.MinRating))
	}

	if session.FilmsState.FilmFilters.MaxRating > 0 {
		queryParams.Add("rating_max", fmt.Sprintf("%.2f", session.FilmsState.FilmFilters.MaxRating))
	}

	if session.FilmsState.FilmSorting.Sort != "" {
		queryParams.Add("sort", session.FilmsState.FilmSorting.Sort)
	}

	requestURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	return requestURL
}
