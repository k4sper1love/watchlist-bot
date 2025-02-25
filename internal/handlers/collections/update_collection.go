package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleUpdateCollectionCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildUpdateCollectionMessage(session), keyboards.BuildCollectionUpdateKeyboard(session))
}

func HandleUpdateCollectionButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallbackUpdateCollectionSelectBack:
		HandleManageCollectionCommand(app, session)

	case states.CallbackUpdateCollectionSelectName:
		handleUpdateCollectionName(app, session)

	case states.CallbackUpdateCollectionSelectDescription:
		handleUpdateCollectionDescription(app, session)
	}
}

func HandleUpdateCollectionProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleUpdateCollectionCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessUpdateCollectionAwaitingName:
		parseCollectionName(app, session, handleUpdateCollectionName, finishUpdateCollectionProcess)

	case states.ProcessUpdateCollectionAwaitingDescription:
		parseCollectionDescription(app, session, handleUpdateCollectionDescription, finishUpdateCollectionProcess)
	}
}

func finishUpdateCollectionProcess(app models.App, session *models.Session) {
	HandleUpdateCollection(app, session, HandleUpdateCollectionCommand)
}

func handleUpdateCollectionName(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildCollectionRequestNameMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateCollectionAwaitingName)
}

func handleUpdateCollectionDescription(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildCollectionRequestDescriptionMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessUpdateCollectionAwaitingDescription)
}

func HandleUpdateCollection(app models.App, session *models.Session, back func(models.App, *models.Session)) {
	if err := updateCollectionAndState(app, session); err != nil {
		app.SendMessage(messages.BuildUpdateCollectionFailureMessage(session), nil)
	} else {
		app.SendMessage(messages.BuildUpdateCollectionSuccessMessage(session), nil)
	}

	session.ClearAllStates()
	back(app, session)
}

func updateCollectionAndState(app models.App, session *models.Session) error {
	collection, err := watchlist.UpdateCollection(app, session)
	if err != nil {
		return err
	}

	session.CollectionDetailState.Collection = *collection
	return nil
}
