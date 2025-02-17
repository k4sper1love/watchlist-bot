package collections

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleUpdateCollectionCommand(app models.App, session *models.Session) {
	msg := messages.BuildCollectionHeader(session)
	choiceMsg := translator.Translate(session.Lang, "choiceAction", nil, nil)
	msg += fmt.Sprintf("<b>%s</b>", choiceMsg)

	keyboard := keyboards.BuildCollectionUpdateKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func HandleUpdateCollectionButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackUpdateCollectionSelectBack:
		HandleManageCollectionCommand(app, session)

	case states.CallbackUpdateCollectionSelectName:
		handleUpdateCollectionName(app, session)

	case states.CallbackUpdateCollectionSelectDescription:
		handleUpdateCollectionDescription(app, session)
	}
}

func HandleUpdateCollectionProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearState()
		session.CollectionDetailState.Clear()
		HandleUpdateCollectionCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessUpdateCollectionAwaitingName:
		parseUpdateCollectionName(app, session)
	case states.ProcessUpdateCollectionAwaitingDescription:
		parseUpdateCollectionDescription(app, session)
	}
}

func handleUpdateCollectionName(app models.App, session *models.Session) {
	msg := "❓" + translator.Translate(session.Lang, "collectionRequestName", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionAwaitingName)
}

func parseUpdateCollectionName(app models.App, session *models.Session) {
	name := utils.ParseMessageString(app.Upd)
	if ok := utils.ValidStringLength(name, 3, 100); !ok {
		validator.HandleInvalidInputLength(app, session, 3, 100)
		handleUpdateCollectionName(app, session)
		return
	}
	session.CollectionDetailState.Name = name

	finishUpdateCollectionProcess(app, session)
}

func handleUpdateCollectionDescription(app models.App, session *models.Session) {
	msg := "❓" + translator.Translate(session.Lang, "collectionRequestDescription", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionAwaitingDescription)
}

func parseUpdateCollectionDescription(app models.App, session *models.Session) {
	description := utils.ParseMessageString(app.Upd)
	if ok := utils.ValidStringLength(description, 0, 500); !ok {
		validator.HandleInvalidInputLength(app, session, 0, 500)
		handleUpdateCollectionDescription(app, session)
		return
	}
	session.CollectionDetailState.Description = description

	finishUpdateCollectionProcess(app, session)
}

func updateCollection(app models.App, session *models.Session) {
	session.CollectionDetailState.IsFavorite = session.CollectionDetailState.Collection.IsFavorite

	collection, err := watchlist.UpdateCollection(app, session)
	if err != nil {
		msg := "🚨" + translator.Translate(session.Lang, "updateCollectionFailure", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackManageCollectionSelectUpdate).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		return
	}

	session.CollectionDetailState.Collection = *collection

	msg := "✏️ " + translator.Translate(session.Lang, "updateCollectionSuccess", nil, nil)
	app.SendMessage(msg, nil)
}

func finishUpdateCollectionProcess(app models.App, session *models.Session) {
	updateCollection(app, session)
	session.ClearAllStates()
	HandleUpdateCollectionCommand(app, session)
}
