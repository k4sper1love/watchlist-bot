package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleFindFilmsCommand(app models.App, session *models.Session) {
	metadata, err := GetFilms(app, session)
	if err != nil {
		app.SendMessage(messages.BuildFilmsFailureMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackFindFilmsBack))
		return
	}

	if metadata.TotalRecords == 0 {
		app.SendMessage(messages.BuildFilmsNotFoundMessage(session), keyboards.BuildFilmsNotFoundKeyboard(session))
		return
	}

	app.SendMessage(messages.BuildFindFilmsMessage(session, metadata), keyboards.BuildFindFilmsKeyboard(session, metadata.CurrentPage, metadata.LastPage))
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

	case states.CallbackFindFilmsNextPage, states.CallbackFindFilmsPrevPage,
		states.CallbackFindFilmsLastPage, states.CallbackFindFilmsFirstPage:
		handleFindFilmsPagination(app, session, callback)
	}
}

func handleFindFilmsPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallbackFindFilmsNextPage:
		if session.FilmsState.CurrentPage >= session.FilmsState.LastPage {
			app.SendMessage(messages.BuildLastPageAlertMessage(session), nil)
			return
		}
		session.FilmsState.CurrentPage++

	case states.CallbackFindFilmsPrevPage:
		if session.FilmsState.CurrentPage <= 1 {
			app.SendMessage(messages.BuildFirstPageAlertMessage(session), nil)
			return
		}
		session.FilmsState.CurrentPage--

	case states.CallbackFindFilmsLastPage:
		if session.FilmsState.CurrentPage == session.FilmsState.LastPage {
			app.SendMessage(messages.BuildLastPageAlertMessage(session), nil)
			return
		}
		session.FilmsState.CurrentPage = session.FilmsState.LastPage

	case states.CallbackFindFilmsFirstPage:
		if session.FilmsState.CurrentPage == 1 {
			app.SendMessage(messages.BuildFirstPageAlertMessage(session), nil)
			return
		}
		session.FilmsState.CurrentPage = 1

	}

	HandleFindFilmsCommand(app, session)
}
