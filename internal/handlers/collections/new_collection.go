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

// HandleNewCollectionCommand handles the command for creating a new collection.
// Prompts the user to enter the name of the new collection and sets the session state accordingly.
func HandleNewCollectionCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestCollectionName(session), keyboards.Cancel(session))
	session.SetState(states.AwaitNewCollectionName)
}

// HandleNewCollectionProcess processes the workflow for creating a new collection.
// Handles states like awaiting the collection name and description input.
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

// requestNewCollectionDescription prompts the user to enter the description for the new collection.
func requestNewCollectionDescription(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestCollectionDescription(session), keyboards.SkipAndCancel(session))
	session.SetState(states.AwaitNewCollectionDescription)
}

// finishNewCollectionProcess finalizes the creation of a new collection.
// Sends the data to the Watchlist service, clears the session states, and navigates to the films associated with the collection.
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
