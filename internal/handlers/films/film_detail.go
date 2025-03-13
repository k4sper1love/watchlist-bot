package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strings"
)

func HandleFilmDetailCommand(app models.App, session *models.Session) {
	if session.FilmDetailState.HasIndex() {
		session.FilmDetailState.Film = session.FilmsState.Films[session.FilmDetailState.Index]
	}

	app.SendImage(
		session.FilmDetailState.Film.ImageURL,
		messages.FilmDetail(session),
		keyboards.FilmDetail(session),
	)
}

func HandleFilmDetailButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallFilmDetailBack:
		session.FilmsState.CurrentPage = 1
		HandleFilmsCommand(app, session)

	case states.CallFilmDetailViewed:
		HandleViewedFilmCommand(app, session)

	case states.CallFilmDetailFavorite:
		makeFavoriteFilm(app, session)

	default:
		if strings.HasPrefix(callback, states.FilmDetailPage) {
			handleFilmDetailPagination(app, session, callback)
		}
	}
}

func handleFilmDetailPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallFilmDetailPageNext:
		if session.FilmDetailState.Index < getFilmsLastIndex(session) {
			session.FilmDetailState.Index++
		} else if err := updateFilmList(app, session, true); err == nil {
			session.FilmDetailState.Index = 0
		} else {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}

	case states.CallFilmDetailPagePrev:
		if session.FilmDetailState.Index > 0 {
			session.FilmDetailState.Index--
		} else if err := updateFilmList(app, session, false); err == nil {
			session.FilmDetailState.Index = getFilmsLastIndex(session)
		} else {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
	}

	HandleFilmDetailCommand(app, session)
}

func makeFavoriteFilm(app models.App, session *models.Session) {
	session.FilmDetailState.SetFavorite(!session.FilmDetailState.Film.IsFavorite)
	HandleUpdateFilm(app, session, HandleFilmDetailCommand)
}

func getFilmsLastIndex(session *models.Session) int {
	return len(session.FilmsState.Films) - 1
}
