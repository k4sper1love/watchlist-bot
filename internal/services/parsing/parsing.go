package parsing

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"strings"
)

var supportedServices = []string{
	"IMDB",
	"kinopoisk",
	"rezka",
	"kinoafisha",
	"youtube",
}

func GetFilmByURL(app models.App, session *models.Session, url string) (*apiModels.Film, error) {
	switch {
	case strings.Contains(url, "imdb"):
		return GetFilmFromIMDB(app, url)

	case strings.Contains(url, "kinopoisk"):
		return GetFilmFromKinopoisk(app, url)

	case strings.Contains(url, "rezka"):
		return GetFilmFromRezka(url)

	case strings.Contains(url, "kinoafisha") && strings.Contains(url, "movies"):
		return GetFilmFromKinoafisha(url)

	case strings.Contains(url, "kinoafisha") && strings.Contains(url, "series"):
		return GetSeriesFromKinoafisha(url)

	case strings.Contains(url, "youtube") || strings.Contains(url, "youtu.be"):
		return GetFilmFromYoutube(app, session, url)

	default:
		return nil, fmt.Errorf("unsupported URL")
	}
}

func GetSupportedServicesInline() string {
	return strings.Join(supportedServices, ", ")
}
