package watchlist

import (
	"encoding/json"
	"fmt"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"net/http"
)

func GetCollectionFilms(app config.App, session *models.Session) (*models.CollectionFilmsResponse, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestUrl := fmt.Sprintf("/collections/%d/films?page=%d", session.CollectionState.ObjectID, session.CollectionFilmState.CurrentPage)

	resp, err := SendRequest(app.Vars.BaseURL, requestUrl, http.MethodGet, nil, headers)
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
