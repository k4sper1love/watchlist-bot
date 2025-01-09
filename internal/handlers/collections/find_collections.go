package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleFindCollectionsCommand(app models.App, session *models.Session) {
	metadata, err := GetCollections(app, session)
	if err != nil {
		app.SendMessage(err.Error(), nil)
		return
	}

	if metadata.TotalRecords == 0 {
		msg := translator.Translate(session.Lang, "collectionsNotFound", nil, nil)
		keyboard := keyboards.NewKeyboard().AddAgain(states.CallbackFindCollectionsAgain).AddBack(states.CallbackFindCollectionsBack).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		return
	}

	msg := messages.BuildCollectionsMessage(session, metadata, true)

	keyboard := keyboards.BuildFindCollectionsKeyboard(session, metadata.CurrentPage, metadata.LastPage)

	app.SendMessage(msg, keyboard)
}

func HandleFindCollectionsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)

	switch {
	case callback == states.CallbackFindCollectionsBack:
		session.ClearAllStates()
		HandleCollectionsCommand(app, session)
		return

	case callback == states.CallbackFindCollectionsAgain:
		session.ClearAllStates()
		handleCollectionsFindByName(app, session)
		return

	case callback == states.CallbackFindCollectionsNextPage:
		if session.CollectionsState.CurrentPage < session.CollectionsState.LastPage {
			session.CollectionsState.CurrentPage++
			HandleCollectionsCommand(app, session)
		} else {
			msg := translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackFindCollectionsPrevPage:
		if session.CollectionsState.CurrentPage > 1 {
			session.CollectionsState.CurrentPage--
			HandleCollectionsCommand(app, session)
		} else {
			msg := translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}
	}
}
