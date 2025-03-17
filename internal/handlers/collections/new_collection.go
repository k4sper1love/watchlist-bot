package collections

import (
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/collectionFilms"
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
// It sends the collection data to the Watchlist service, clears the session states,
// and processes the next steps based on the session context.
func finishNewCollectionProcess(app models.App, session *models.Session) {
	collection, err := watchlist.CreateCollection(app, session)
	session.ClearAllStates() // Clears session states before processing the collection context
	if err != nil {
		app.SendMessage(messages.CreateCollectionFailure(session), keyboards.Back(session, states.CallCollectionsNew))
		return
	}

	handleNewCollectionContext(app, session, collection)
}

// handleNewCollectionContext processes the collection context after creation.
// Depending on the session context, it either navigates to the newly created collection's films
// or associates the collection with an existing film.
func handleNewCollectionContext(app models.App, session *models.Session, collection *apiModels.Collection) {
	switch session.Context {
	case states.CtxCollection:
		// If the context is a collection, store its ID and proceed to handling films within it.
		session.CollectionDetailState.ObjectID = collection.ID
		app.SendMessage(messages.CreateCollectionSuccess(session), nil)
		setContextAndHandleFilms(app, session)
	case states.CtxFilm:
		// If the context is a film, associate it with the newly created collection.
		session.CollectionDetailState.Collection.ID = collection.ID
		collectionFilms.AddFilmToCollection(app, session)
	}
}
