package parsing

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"strings"
)

// supportedServices contains a list of supported external services for film parsing.
var supportedServices = []string{
	"imdb",
	"kinopoisk",
	"rezka",
	"kinoafisha",
	"youtube",
}

// GetFilmByURL parses a film from a given URL based on the supported service.
// It determines the service from the URL and delegates parsing to the corresponding function.
func GetFilmByURL(app models.App, session *models.Session, url string) (*apiModels.Film, error) {
	switch {
	case strings.Contains(url, supportedServices[0]):
		// Parse film from IMDB.
		return GetFilmFromIMDB(app, url)

	case strings.Contains(url, supportedServices[1]):
		// Parse film from Kinopoisk.
		return GetFilmFromKinopoisk(session, url)

	case strings.Contains(url, supportedServices[2]):
		// Parse film from Rezka.
		return GetFilmFromRezka(url)

	case strings.Contains(url, supportedServices[3]) && strings.Contains(url, "movies"):
		// Parse film from Kinoafisha (movies section).
		return GetFilmFromKinoafisha(url)

	case strings.Contains(url, supportedServices[3]) && strings.Contains(url, "series"):
		// Parse series from Kinoafisha (series section).
		return GetSeriesFromKinoafisha(url)

	case strings.Contains(url, supportedServices[4]) || strings.Contains(url, "youtu.be"):
		// Parse film from YouTube or youtu.be links.
		return GetFilmFromYoutube(app, session, url)

	default:
		// Return an error if the URL is not supported.
		return nil, fmt.Errorf("unsupported URL")
	}
}

// GetSupportedServicesInline returns a comma-separated string of supported services.
// This can be used to inform users about the platforms they can provide URLs from.
func GetSupportedServicesInline() string {
	return strings.Join(supportedServices, ", ")
}

// IsKinopoisk checks if the given URL belongs to the Kinopoisk service.
func IsKinopoisk(url string) bool {
	return strings.Contains(url, supportedServices[1])
}
