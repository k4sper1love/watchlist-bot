package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strings"
)

// HandleFindFilmsCommand handles the command for searching and listing films by title.
// Retrieves paginated films matching the search criteria and sends a message with their details and navigation buttons.
func HandleFindFilmsCommand(app models.App, session *models.Session) {
	if metadata, err := getFilms(app, session); err != nil {
		app.SendMessage(messages.FilmsFailure(session), keyboards.FilmsNotFound(session))
	} else {
		app.SendMessage(messages.FindFilms(session, metadata), keyboards.FindFilms(session, metadata.CurrentPage, metadata.LastPage))
	}
}

// HandleFindFilmsButtons handles button interactions related to the search results of films.
// Supports actions like going back, refreshing the search, and pagination.
func HandleFindFilmsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallFindFilmsBack:
		session.ClearAllStates()
		HandleFilmsCommand(app, session)

	case states.CallFindFilmsAgain:
		session.ClearAllStates()
		requestFindFilmsTitle(app, session)

	default:
		if strings.HasPrefix(callback, states.FindFilmsPage) {
			handleFindFilmsPagination(app, session, callback)
		}
	}
}

// handleFindFilmsPagination processes pagination actions for the search results of films.
// Updates the current page in the session and reloads the films list.
func handleFindFilmsPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallFindFilmsPageNext:
		if session.FilmsState.CurrentPage >= session.FilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage++

	case states.CallFindFilmsPagePrev:
		if session.FilmsState.CurrentPage <= 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage--

	case states.CallFindFilmsPageLast:
		if session.FilmsState.CurrentPage == session.FilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage = session.FilmsState.LastPage

	case states.CallFindFilmsPageFirst:
		if session.FilmsState.CurrentPage == 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage = 1
	}

	HandleFindFilmsCommand(app, session)
}
