package watchlist

import (
	"encoding/json"
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"io"
	"net/http"
)

type collectionRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
}

func GetCollections(app config.App, session *models.Session) (*models.CollectionsResponse, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("/collections?page=%d&page_size=%d", session.CollectionState.CurrentPage, session.CollectionState.PageSize)

	resp, err := SendRequest(app.Vars.BaseURL, requestURL, http.MethodGet, nil, headers)
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

func CreateCollection(app config.App, session *models.Session) (*apiModels.Collection, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	data := collectionRequest{
		Name:        session.CollectionState.Name,
		Description: session.CollectionState.Description,
	}

	resp, err := SendRequest(app.Vars.BaseURL, "/collections", http.MethodPost, data, headers)
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

func parseCollection(dest *apiModels.Collection, data io.Reader) error {
	return json.NewDecoder(data).Decode(&struct {
		Collection *apiModels.Collection `json:"collection"`
	}{
		Collection: dest,
	})
}
