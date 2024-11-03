package watchlist

import (
	"encoding/json"
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"net/http"
)

func GetFilms(app models.App, session *models.Session) (*models.FilmsResponse, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("%s/collections/%d/films", app.Vars.BaseURL, session.CollectionDetailState.ObjectID)

	resp, err := SendRequest(requestURL, http.MethodGet, nil, headers)
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
