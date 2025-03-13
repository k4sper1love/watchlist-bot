package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleNewCollectionCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestCollectionName(session), keyboards.Cancel(session))
	session.SetState(states.AwaitNewCollectionName)
}

func HandleNewCollectionProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleCollectionsCommand(app, session)
		return
	}

	switch session.State {
	case states.AwaitNewCollectionName:
		parser.ParseCollectionName(app, session, HandleNewCollectionCommand, requestNewCollectionDescription)

	case states.AwaitNewCollectionDescription:
		parser.ParseCollectionDescription(app, session, requestNewCollectionDescription, finishNewCollectionProcess)
	}
}

func requestNewCollectionDescription(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestCollectionDescription(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitNewCollectionDescription)
}

func finishNewCollectionProcess(app models.App, session *models.Session) {
	collection, err := watchlist.CreateCollection(app, session)
	session.ClearAllStates()
	if err != nil {
		app.SendMessage(messages.CreateCollectionFailure(session), keyboards.Back(session, states.CallCollectionsNew))
		return
	}

	session.CollectionDetailState.ObjectID = collection.ID
	app.SendMessage(messages.CreateCollectionSuccess(session), nil)
	setContextAndHandleFilms(app, session)
}
