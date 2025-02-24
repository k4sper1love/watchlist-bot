package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleFilmsDetailCommand(app models.App, session *models.Session) {
	if session.FilmDetailState.HasIndex() {
		session.FilmDetailState.Film = session.FilmsState.Films[session.FilmDetailState.Index]
	}

	app.SendImage(
		session.FilmDetailState.Film.ImageURL,
		messages.BuildFilmDetailMessage(session),
		keyboards.BuildFilmDetailKeyboard(session),
	)
}

func HandleFilmsDetailButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallbackFilmDetailNextPage, states.CallbackFilmDetailPrevPage:
		handleFilmDetailPagination(app, session, callback)

	case states.CallbackFilmDetailBack:
		session.FilmsState.CurrentPage = 1
		HandleFilmsCommand(app, session)

	case states.CallbackFilmDetailViewed:
		HandleViewedFilmCommand(app, session)

	case states.CallbackFilmDetailFavorite:
		handleFavoriteFilm(app, session)
	}
}

func handleFilmDetailPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallbackFilmDetailNextPage:
		if session.FilmDetailState.Index < getFilmsLastIndex(session) {
			session.FilmDetailState.Index++
		} else if err := UpdateFilmsList(app, session, true); err == nil {
			session.FilmDetailState.Index = 0
		} else {
			app.SendMessage(messages.BuildLastPageAlertMessage(session), nil)
			return
		}

	case states.CallbackFilmDetailPrevPage:
		if session.FilmDetailState.Index > 0 {
			session.FilmDetailState.Index--
		} else if err := UpdateFilmsList(app, session, false); err == nil {
			session.FilmDetailState.Index = getFilmsLastIndex(session)
		} else {
			app.SendMessage(messages.BuildFirstPageAlertMessage(session), nil)
			return
		}
	}

	HandleFilmsDetailCommand(app, session)
}

func handleFavoriteFilm(app models.App, session *models.Session) {
	session.FilmDetailState.SetFavorite(!session.FilmDetailState.Film.IsFavorite)
	HandleUpdateFilm(app, session, HandleFilmsDetailCommand)
}

func getFilmsLastIndex(session *models.Session) int {
	return len(session.FilmsState.Films) - 1
}
