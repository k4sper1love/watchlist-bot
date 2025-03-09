package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleFindCollectionsCommand(app models.App, session *models.Session) {
	if metadata, err := getCollections(app, session); err != nil {
		app.SendMessage(messages.CollectionsFailure(session), keyboards.BuildKeyboardWithBack(session, states.CallbackFindCollectionsBack))
	} else {
		app.SendMessage(messages.Collections(session, metadata, true), keyboards.BuildFindCollectionsKeyboard(session, metadata.CurrentPage, metadata.LastPage))
	}
}

func HandleFindCollectionsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallbackFindCollectionsBack:
		session.ClearAllStates()
		HandleCollectionsCommand(app, session)

	case states.CallbackFindCollectionsAgain:
		session.ClearAllStates()
		handleCollectionsFindByName(app, session)

	case states.CallbackFindCollectionsNextPage, states.CallbackFindCollectionsPrevPage,
		states.CallbackFindCollectionsLastPage, states.CallbackFindCollectionsFirstPage:
		handleFindCollectionPagination(app, session, callback)
	}
}

func handleFindCollectionPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallbackFindCollectionsNextPage:
		if session.CollectionsState.CurrentPage >= session.CollectionsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.CollectionsState.CurrentPage++

	case states.CallbackFindCollectionsPrevPage:
		if session.CollectionsState.CurrentPage <= 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.CollectionsState.CurrentPage--

	case states.CallbackFindCollectionsLastPage:
		if session.CollectionsState.CurrentPage == session.CollectionsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.CollectionsState.CurrentPage = session.CollectionsState.LastPage

	case states.CallbackFindCollectionsFirstPage:
		if session.CollectionsState.CurrentPage == 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.CollectionsState.CurrentPage = 1
	}

	HandleCollectionsCommand(app, session)
}
