package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strings"
)

// HandleFindCollectionsCommand handles the command for searching and listing collections by name.
// Retrieves paginated collections matching the search criteria and sends a message with their details and navigation buttons.
func HandleFindCollectionsCommand(app models.App, session *models.Session) {
	if metadata, err := getCollections(app, session); err != nil {
		app.SendMessage(messages.CollectionsFailure(session), keyboards.Back(session, states.CallFindCollectionsBack))
	} else {
		app.SendMessage(messages.Collections(session, metadata, true), keyboards.FindCollections(session, metadata.CurrentPage, metadata.LastPage))
	}
}

// HandleFindCollectionsButtons handles button interactions related to the collection search results.
// Supports actions like going back, refreshing the search, and pagination.
func HandleFindCollectionsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallFindCollectionsBack:
		session.ClearAllStates()
		HandleCollectionsCommand(app, session)

	case states.CallFindCollectionsAgain:
		session.ClearAllStates()
		handleCollectionsFindByName(app, session)

	default:
		if strings.HasPrefix(callback, states.FindCollectionsPage) {
			handleFindCollectionPagination(app, session, callback)
		}
	}
}

// handleFindCollectionPagination processes pagination actions for the search results of collections.
// Updates the current page in the session and reloads the collections list.
func handleFindCollectionPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallFindCollectionsPageNext:
		if session.CollectionsState.CurrentPage >= session.CollectionsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.CollectionsState.CurrentPage++

	case states.CallFindCollectionsPagePrev:
		if session.CollectionsState.CurrentPage <= 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.CollectionsState.CurrentPage--

	case states.CallFindCollectionsPageLast:
		if session.CollectionsState.CurrentPage == session.CollectionsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.CollectionsState.CurrentPage = session.CollectionsState.LastPage

	case states.CallFindCollectionsPageFirst:
		if session.CollectionsState.CurrentPage == 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.CollectionsState.CurrentPage = 1
	}

	HandleCollectionsCommand(app, session)
}
