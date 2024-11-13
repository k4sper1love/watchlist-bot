package watchlist

import (
	"encoding/json"
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"net/http"
)

func GetCollections(app models.App, session *models.Session) (*models.CollectionsResponse, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("%s/api/v1/collections?page=%d&page_size=%d", app.Vars.Host, session.CollectionsState.CurrentPage, session.CollectionsState.PageSize)

	resp, err := SendRequest(requestURL, http.MethodGet, nil, headers)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("get_collection failed: %s", resp.Status)
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

	resp, err := SendRequest(app.Vars.Host+"/api/v1/collections", http.MethodPost, session.CollectionDetailState, headers)
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

	requestURL := fmt.Sprintf("%s/api/v1/collections/%d", app.Vars.Host, session.CollectionDetailState.ObjectID)

	resp, err := SendRequest(requestURL, http.MethodPut, session.CollectionDetailState, headers)
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

	requestURL := fmt.Sprintf("%s/api/v1/collections/%d", app.Vars.Host, session.CollectionDetailState.ObjectID)

	resp, err := SendRequest(requestURL, http.MethodDelete, nil, headers)
	if err != nil {
		return fmt.Errorf("failed to send request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("delete_collection failed: %s", resp.Status)
	}

	return nil
}
