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

	default:
		return nil, fmt.Errorf("unsupported URL")
	}
}
