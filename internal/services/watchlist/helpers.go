package watchlist

import (
	"encoding/json"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/security"
	"io"
)

// parseAuth parses the authentication response from the API and populates the session with user data.
// It extracts the user details, encrypts the access and refresh tokens, and stores them in the session.
func parseAuth(dest *models.Session, data io.Reader) error {
	auth := &apiModels.AuthResponse{}

	// Decode the JSON response into a structured format.
	err := json.NewDecoder(data).Decode(&struct {
		Auth *apiModels.AuthResponse `json:"user"`
	}{
		Auth: auth,
	})
	if err != nil {
		return err
	}

	dest.User = *auth.User // Populate the session with user details.

	// Encrypt the access token for secure storage.
	encryptedAccessToken, err := security.Encrypt(auth.AccessToken)
	if err != nil {
		return err
	}

	// Encrypt the refresh token for secure storage.
	encryptedRefreshToken, err := security.Encrypt(auth.RefreshToken)
	if err != nil {
		return err
	}

	dest.AccessToken = encryptedAccessToken
	dest.RefreshToken = encryptedRefreshToken
	return nil
}

// parseUser parses the user details from the API response and populates the provided `models.User` object.
func parseUser(dest *apiModels.User, data io.Reader) error {
	return json.NewDecoder(data).Decode(&struct {
		User *apiModels.User `json:"user"`
	}{
		User: dest,
	})
}

// parseFilm parses the film details from the API response and populates the provided `models.Film` object.
func parseFilm(dest *apiModels.Film, data io.Reader) error {
	return json.NewDecoder(data).Decode(&struct {
		Film *apiModels.Film `json:"film"`
	}{
		Film: dest,
	})
}

// parseCollection parses the collection details from the API response and populates the provided `models.Collection` object.
func parseCollection(dest *apiModels.Collection, data io.Reader) error {
	return json.NewDecoder(data).Decode(&struct {
		Collection *apiModels.Collection `json:"collection"`
	}{
		Collection: dest,
	})
}

// parseCollectionFilm parses the collection-film relationship details from the API response
// and populates the provided `models.CollectionFilm` object.
func parseCollectionFilm(dest *apiModels.CollectionFilm, data io.Reader) error {
	return json.NewDecoder(data).Decode(&struct {
		CollectionFilm *apiModels.CollectionFilm `json:"collection_film"`
	}{
		CollectionFilm: dest,
	})
}

// parseImageURL extracts the image URL from the API response.
func parseImageURL(data io.Reader) (string, error) {
	var result struct {
		ImageURL string `json:"image_url"` // Field to extract the image URL.
	}

	if err := json.NewDecoder(data).Decode(&result); err != nil {
		return "", err
	}
	return result.ImageURL, nil
}
