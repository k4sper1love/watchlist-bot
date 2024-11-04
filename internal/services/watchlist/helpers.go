package watchlist

import (
	"encoding/json"
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"io"
)

func parseAuth(dest *models.Session, data io.Reader) error {
	var responseMap map[string]interface{}
	if err := json.NewDecoder(data).Decode(&responseMap); err != nil {
		return err
	}

	userData := responseMap["user"].(map[string]interface{})

	if id, ok := userData["id"].(float64); ok {
		dest.UserID = int(id)
	} else {
		return fmt.Errorf("invalid id format")
	}
	dest.AccessToken = userData["access_token"].(string)
	dest.RefreshToken = userData["refresh_token"].(string)

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
