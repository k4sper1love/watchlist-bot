package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleNewCollectionCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildCollectionRequestNameMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessNewCollectionAwaitingName)
}

func HandleNewCollectionProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleCollectionsCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessNewCollectionAwaitingName:
		parseCollectionName(app, session, HandleNewCollectionCommand, requestNewCollectionDescription)

	case states.ProcessNewCollectionAwaitingDescription:
		parseCollectionDescription(app, session, requestNewCollectionDescription, finishNewCollectionProcess)
	}
}

func requestNewCollectionDescription(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildCollectionRequestDescriptionMessage(session), keyboards.BuildKeyboardWithSkipAndCancel(session))
	session.SetState(states.ProcessNewCollectionAwaitingDescription)
}

func finishNewCollectionProcess(app models.App, session *models.Session) {
	collection, err := watchlist.CreateCollection(app, session)
	session.ClearAllStates()
	if err != nil {
		app.SendMessage(messages.BuildCreateCollectionFailureMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackCollectionsNew))
		return
	}

	session.CollectionDetailState.ObjectID = collection.ID
	app.SendMessage(messages.BuildCreateCollectionSuccessMessage(session), nil)
	setContextAndHandleFilms(app, session)
}
