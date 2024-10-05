package watchlist

import (
	"encoding/json"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"net/http"
)

func GetCollections(app config.App, session models.Session) (*models.CollectionsResponse, error) {
	resp, err := SendRequest(app.Vars.BaseURL, "/collections", http.MethodGet, session.AccessToken, nil)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var collectionsResponse models.CollectionsResponse
	if err := json.NewDecoder(resp.Body).Decode(&collectionsResponse); err != nil {
		return nil, err
	}

	return &collectionsResponse, nil
}
