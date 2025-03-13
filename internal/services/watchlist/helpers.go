package watchlist

import (
	"encoding/json"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"io"
)

func parseAuth(dest *models.Session, data io.Reader) error {
	auth := &apiModels.AuthResponse{}

	err := json.NewDecoder(data).Decode(&struct {
		Auth *apiModels.AuthResponse `json:"user"`
	}{
		Auth: auth,
	})
	if err != nil {
		return err
	}

	dest.User = *auth.User
	dest.AccessToken = auth.AccessToken
	dest.RefreshToken = auth.RefreshToken

	return nil
}

func parseUser(dest *apiModels.User, data io.Reader) error {
	return json.NewDecoder(data).Decode(&struct {
		User *apiModels.User `json:"user"`
	}{
		User: dest,
	})
}

func parseFilm(dest *apiModels.Film, data io.Reader) error {
	return json.NewDecoder(data).Decode(&struct {
		Film *apiModels.Film `json:"film"`
	}{
		Film: dest,
	})
}

func parseCollection(dest *apiModels.Collection, data io.Reader) error {
	return json.NewDecoder(data).Decode(&struct {
		Collection *apiModels.Collection `json:"collection"`
	}{
		Collection: dest,
	})
}

func parseCollectionFilm(dest *apiModels.CollectionFilm, data io.Reader) error {
	return json.NewDecoder(data).Decode(&struct {
		CollectionFilm *apiModels.CollectionFilm `json:"collection_film"`
	}{
		CollectionFilm: dest,
	})
}

func parseImageURL(data io.Reader) (string, error) {
	var result struct {
		ImageURL string `json:"image_url"`
	}

	err := json.NewDecoder(data).Decode(&result)
	if err != nil {
		return "", err
	}

	return result.ImageURL, nil
}
