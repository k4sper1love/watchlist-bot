package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// HandleManageCollectionCommand handles the command for managing a collection.
// Sends a message with options to update or delete the selected collection.
func HandleManageCollectionCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.CollectionChoiceAction(session), keyboards.CollectionManage(session))
}

// HandleManageCollectionButtons handles button interactions related to managing a collection.
// Supports actions like going back, updating the collection, or deleting it.
func HandleManageCollectionButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallManageCollectionBack:
		// Sets the context back to "collection" and navigates to the films associated with the collection.
		session.SetContext(states.CtxCollection)
		films.HandleFilmsCommand(app, session)

	case states.CallManageCollectionUpdate:
		HandleUpdateCollectionCommand(app, session)

	case states.CallManageCollectionDelete:
		HandleDeleteCollectionCommand(app, session)
	}
}
