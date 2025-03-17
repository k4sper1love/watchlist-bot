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

// HandleUpdateCollectionCommand handles the command for updating a collection.
// Sends a message with options to update the collection's name or description.
func HandleUpdateCollectionCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.UpdateCollection(session), keyboards.CollectionUpdate(session))
}

// HandleUpdateCollectionButtons handles button interactions related to updating a collection.
// Supports actions like going back, updating the collection's name, or updating its description.
func HandleUpdateCollectionButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallUpdateCollectionBack:
		HandleManageCollectionCommand(app, session)

	case states.CallUpdateCollectionName:
		handleUpdateCollectionName(app, session)

	case states.CallUpdateCollectionDescription:
		handleUpdateCollectionDescription(app, session)
	}
}

// HandleUpdateCollectionProcess processes the workflow for updating a collection.
// Handles states like awaiting input for the collection's name or description.
func HandleUpdateCollectionProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleUpdateCollectionCommand(app, session)
		return
	}

	switch session.State {
	case states.AwaitUpdateCollectionName:
		parser.ParseCollectionName(app, session, handleUpdateCollectionName, finishUpdateCollectionProcess)

	case states.AwaitUpdateCollectionDescription:
		parser.ParseCollectionDescription(app, session, handleUpdateCollectionDescription, finishUpdateCollectionProcess)
	}
}

// finishUpdateCollectionProcess finalizes the update of a collection.
func finishUpdateCollectionProcess(app models.App, session *models.Session) {
	HandleUpdateCollection(app, session, HandleUpdateCollectionCommand)
}

// handleUpdateCollectionName prompts the user to enter a new name for the collection.
func handleUpdateCollectionName(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestCollectionName(session), keyboards.Cancel(session))
	session.SetState(states.AwaitUpdateCollectionName)
}

// handleUpdateCollectionDescription prompts the user to enter a new description for the collection.
func handleUpdateCollectionDescription(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestCollectionDescription(session), keyboards.Cancel(session))
	session.SetState(states.AwaitUpdateCollectionDescription)
}

// HandleUpdateCollection updates the collection using the Watchlist service and resets the session state.
func HandleUpdateCollection(app models.App, session *models.Session, back func(models.App, *models.Session)) {
	if err := updateCollectionAndState(app, session); err != nil {
		app.SendMessage(messages.UpdateCollectionFailure(session), nil)
	} else {
		app.SendMessage(messages.UpdateCollectionSuccess(session), nil)
	}

	session.ClearAllStates()
	back(app, session)
}

// updateCollectionAndState updates the collection in the database and synchronizes the session state with the updated data.
func updateCollectionAndState(app models.App, session *models.Session) error {
	collection, err := watchlist.UpdateCollection(app, session)
	if err != nil {
		return err
	}

	session.CollectionDetailState.Collection = *collection
	return nil
}
