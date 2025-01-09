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

	requestURL := buildGetCollectionsURL(app, session, filmID, excludeFilmID, currentPage, pageSize)

	resp, err := client.SendRequestWithOptions(requestURL, http.MethodGet, nil, headers)
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

	resp, err := client.SendRequestWithOptions(app.Vars.Host+"/api/v1/collections", http.MethodPost, session.CollectionDetailState, headers)
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

	resp, err := client.SendRequestWithOptions(requestURL, http.MethodPut, session.CollectionDetailState, headers)
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

	resp, err := client.SendRequestWithOptions(requestURL, http.MethodDelete, nil, headers)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete_collection failed: %s", resp.Status)
	}

	return nil
}

func buildGetCollectionsURL(app models.App, session *models.Session, filmID, excludeFilmID, currentPage, pageSize int) string {
	baseURL := fmt.Sprintf("%s/api/v1/collections", app.Vars.Host)
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

	requestURL := fmt.Sprintf("%s?%s", baseURL, queryParams.Encode())

	return requestURL
}
