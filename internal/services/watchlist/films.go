package watchlist

import (
	"encoding/json"
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
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

func UpdateFilm(app models.App, session *models.Session) (*apiModels.Film, error) {
	headers := map[string]string{
		"Authorization": session.AccessToken,
	}

	requestURL := fmt.Sprintf("%s/films/%d", app.Vars.BaseURL, session.CollectionFilmState.Object.ID)

	resp, err := SendRequest(requestURL, http.MethodPut, session.CollectionFilmState, headers)
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
