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

func GetCollectionFilms(app models.App, session *models.Session) (*models.CollectionFilmsResponse, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := buildGetCollectionFilmsURL(app, session)

	resp, err := client.SendRequestWithOptions(requestURL, http.MethodGet, nil, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get_collection_films failed: %s", resp.Status)
	}

	var collectionFilmsResponse models.CollectionFilmsResponse
	if err := json.NewDecoder(resp.Body).Decode(&collectionFilmsResponse); err != nil {
		return nil, err
	}

	return &collectionFilmsResponse, nil
}

func CreateCollectionFilm(app models.App, session *models.Session) (*apiModels.CollectionFilm, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("%s/api/v1/collections/%d/films", app.Vars.Host, session.CollectionDetailState.Collection.ID)

	resp, err := client.SendRequestWithOptions(requestURL, http.MethodPost, session.FilmDetailState, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create_collection_film failed: %d", resp.StatusCode)
	}

	collectionFilm := &apiModels.CollectionFilm{}
	if err := parseCollectionFilm(collectionFilm, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to parse collection film: %w", err)
	}

	return collectionFilm, nil
}

func AddCollectionFilm(app models.App, session *models.Session) (*apiModels.CollectionFilm, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("%s/api/v1/collections/%d/films/%d", app.Vars.Host, session.CollectionDetailState.Collection.ID, session.FilmDetailState.Film.ID)

	resp, err := client.SendRequestWithOptions(requestURL, http.MethodPost, nil, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("add_collection_film failed: %d", resp.StatusCode)
	}

	collectionFilm := &apiModels.CollectionFilm{}
	if err := parseCollectionFilm(collectionFilm, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to parse collection film: %w", err)
	}

	return collectionFilm, nil
}

func DeleteCollectionFilm(app models.App, session *models.Session) error {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("%s/api/v1/collections/%d/films/%d", app.Vars.Host, session.CollectionDetailState.Collection.ID, session.FilmDetailState.Film.ID)

	resp, err := client.SendRequestWithOptions(requestURL, http.MethodDelete, nil, headers)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete_collection_film failed: %s", resp.Status)
	}

	return nil
}

func buildGetCollectionFilmsURL(app models.App, session *models.Session) string {
	baseURL := fmt.Sprintf("%s/api/v1/collections/%d/films", app.Vars.Host, session.CollectionDetailState.ObjectID)
	queryParams := url.Values{}

	state := session.FilmsState

	queryParams = addFilmsBasicParams(queryParams, state.Title, state.CurrentPage, state.PageSize)

	queryParams = addFilmsFilterAndSortingParams(queryParams, state.CollectionFilters, state.CollectionSorting)

	requestURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	return requestURL
}
