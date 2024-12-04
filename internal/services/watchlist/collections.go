package watchlist

import (
	"encoding/json"
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/client"
	"net/http"
)

func GetCollections(app models.App, session *models.Session) (*models.CollectionsResponse, error) {
	return getCollectionsRequest(app, session, -1, -1, session.CollectionsState.CurrentPage, session.CollectionsState.PageSize)
}

func GetCollectionsByFilm(app models.App, session *models.Session) (*models.CollectionsResponse, error) {
	return getCollectionsRequest(app, session, session.FilmDetailState.Film.ID, -1, session.CollectionsState.CurrentPage, session.CollectionsState.PageSize)
}

func GetCollectionsExcludeFilm(app models.App, session *models.Session) (*models.CollectionsResponse, error) {
	return getCollectionsRequest(app, session, -1, session.FilmDetailState.Film.ID, session.CollectionFilmsState.CurrentPage, session.CollectionFilmsState.PageSize)
}

func getCollectionsRequest(app models.App, session *models.Session, filmID, excludeFilmID, currentPage, pageSize int) (*models.CollectionsResponse, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("%s/api/v1/collections?film=%d&exclude_film=%d&page=%d&page_size=%d", app.Vars.Host, filmID, excludeFilmID, currentPage, pageSize)

	resp, err := client.SendRequest(requestURL, http.MethodGet, nil, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get_collections failed: %s", resp.Status)
	}

	var collectionsResponse models.CollectionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&collectionsResponse); err != nil {
		return nil, err
	}

	return &collectionsResponse, nil
}

func CreateCollection(app models.App, session *models.Session) (*apiModels.Collection, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	resp, err := client.SendRequest(app.Vars.Host+"/api/v1/collections", http.MethodPost, session.CollectionDetailState, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return nil, fmt.Errorf("create_collection failed: %s", resp.Status)
	}

	collection := &apiModels.Collection{}
	if err := parseCollection(collection, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to parse collection: %w", err)
	}

	return collection, nil
}

func UpdateCollection(app models.App, session *models.Session) (*apiModels.Collection, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("%s/api/v1/collections/%d", app.Vars.Host, session.CollectionDetailState.Collection.ID)

	resp, err := client.SendRequest(requestURL, http.MethodPut, session.CollectionDetailState, headers)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("update_collection failed: %s", resp.Status)
	}

	collection := &apiModels.Collection{}
	if err := parseCollection(collection, resp.Body); err != nil {
		return nil, fmt.Errorf("failed to parse collection: %w", err)
	}

	return collection, nil
}

func DeleteCollection(app models.App, session *models.Session) error {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("%s/api/v1/collections/%d", app.Vars.Host, session.CollectionDetailState.Collection.ID)

	resp, err := client.SendRequest(requestURL, http.MethodDelete, nil, headers)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete_collection failed: %s", resp.Status)
	}

	return nil
}
