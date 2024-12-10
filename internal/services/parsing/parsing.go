package parsing

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"strings"
)

func GetFilmByURL(app models.App, url string) (*apiModels.Film, error) {
	switch {
	case strings.Contains(url, "kinopoisk"):
		return GetFilmFromKinopoisk(app, url)

	case strings.Contains(url, "rezka"):
		return GetFilmFromRezka(url)

	case strings.Contains(url, "kinoafisha") && strings.Contains(url, "movies"):
		return GetFilmFromKinoafisha(url)

	case strings.Contains(url, "kinoafisha") && strings.Contains(url, "series"):
		return GetSeriesFromKinoafisha(url)

	default:
		return nil, fmt.Errorf("unsupported URL")
	}
}
