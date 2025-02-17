package parsing

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"strings"
)

var supportedServices = []string{
	"imdb",
	"kinopoisk",
	"rezka",
	"kinoafisha",
	"youtube",
}

func GetFilmByURL(app models.App, session *models.Session, url string) (*apiModels.Film, error) {
	switch {
	case strings.Contains(url, supportedServices[0]):
		return GetFilmFromIMDB(app, url)

	case strings.Contains(url, supportedServices[1]):
		return GetFilmFromKinopoisk(session, url)

	case strings.Contains(url, supportedServices[2]):
		return GetFilmFromRezka(url)

	case strings.Contains(url, supportedServices[3]) && strings.Contains(url, "movies"):
		return GetFilmFromKinoafisha(url)

	case strings.Contains(url, supportedServices[3]) && strings.Contains(url, "series"):
		return GetSeriesFromKinoafisha(url)

	case strings.Contains(url, supportedServices[4]) || strings.Contains(url, "youtu.be"):
		return GetFilmFromYoutube(app, session, url)

	default:
		return nil, fmt.Errorf("unsupported URL")
	}
}

func GetSupportedServicesInline() string {
	return strings.Join(supportedServices, ", ")
}

func IsKinopoisk(url string) bool {
	return strings.Contains(url, supportedServices[1])
}
