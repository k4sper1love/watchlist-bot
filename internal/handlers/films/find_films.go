package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleFindFilmsCommand(app models.App, session *models.Session) {
	if metadata, err := getFilms(app, session); err != nil {
		app.SendMessage(messages.FilmsFailure(session), keyboards.FilmsNotFound(session))
	} else {
		app.SendMessage(messages.FindFilms(session, metadata), keyboards.FindFilms(session, metadata.CurrentPage, metadata.LastPage))
	}
}

func HandleFindFilmsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallbackFindFilmsBack:
		session.ClearAllStates()
		HandleFilmsCommand(app, session)

	case states.CallbackFindFilmsAgain:
		session.ClearAllStates()
		handleFilmsFindByTitle(app, session)

	case states.CallbackFindFilmsPageNext, states.CallbackFindFilmsPagePrev,
		states.CallbackFindFilmsPageLast, states.CallbackFindFilmsPageFirst:
		handleFindFilmsPagination(app, session, callback)
	}
}

func handleFindFilmsPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallbackFindFilmsPageNext:
		if session.FilmsState.CurrentPage >= session.FilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage++

	case states.CallbackFindFilmsPagePrev:
		if session.FilmsState.CurrentPage <= 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage--

	case states.CallbackFindFilmsPageLast:
		if session.FilmsState.CurrentPage == session.FilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage = session.FilmsState.LastPage

	case states.CallbackFindFilmsPageFirst:
		if session.FilmsState.CurrentPage == 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage = 1
	}

	HandleFindFilmsCommand(app, session)
}
