package watchlist

import (
	"encoding/json"
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"net/http"
)

func GetCollectionFilms(app models.App, session *models.Session) (*models.CollectionFilmsResponse, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("%s/collections/%d/films?page=%d&page_size=%d", app.Vars.BaseURL, session.CollectionDetailState.ObjectID, session.CollectionDetailState.CurrentPage, session.CollectionDetailState.PageSize)

	resp, err := SendRequest(requestURL, http.MethodGet, nil, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get_films failed: %s", resp.Status)
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

	requestURL := fmt.Sprintf("%s/collections/%d/films", app.Vars.BaseURL, session.CollectionDetailState.ObjectID)

	resp, err := SendRequest(requestURL, http.MethodPost, session.CollectionFilmState, headers)
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

func DeleteCollectionFilm(app models.App, session *models.Session) error {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("%s/collections/%d/films/%d", app.Vars.BaseURL, session.CollectionDetailState.ObjectID, session.CollectionFilmState.Object.ID)

	resp, err := SendRequest(requestURL, http.MethodDelete, nil, headers)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete_collection_film failed: %s", resp.Status)
	}

	return nil
}
