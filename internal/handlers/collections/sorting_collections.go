package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strings"
)

// HandleSortingCollectionsCommand handles the command for sorting collections.
// Sends a message with options to choose a sorting field and direction.
func HandleSortingCollectionsCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.ChoiceSorting(session), keyboards.CollectionsSorting(session))
}

// HandleSortingCollectionsButtons handles button interactions related to sorting collections.
// Supports actions like going back, resetting sorting, or selecting a sorting field.
func HandleSortingCollectionsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallCollectionSortingBack:
		handleWithResetCollectionsPage(app, session, HandleCollectionsCommand)

	case states.CallCollectionSortingAllReset:
		handleSortingCollectionsReset(app, session, HandleCollectionsCommand)

	default:
		handleSortingCollectionsSelect(app, session, callback)
	}
}

// handleSortingCollectionsSelect processes the selection of a sorting field.
// Updates the session state with the selected field and prompts the user to choose a sorting direction.
func handleSortingCollectionsSelect(app models.App, session *models.Session, callback string) {
	if strings.HasPrefix(callback, states.CollectionSortingSelect) {
		session.CollectionsState.Sorting.Field = strings.TrimPrefix(callback, states.CollectionSortingSelect)
		handleSortingCollectionsDirection(app, session)
	}
}

// HandleSortingCollectionsProcess processes the workflow for sorting collections.
// Handles states like awaiting the sorting direction input.
func HandleSortingCollectionsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		handleWithResetCollectionsPage(app, session, HandleSortingCollectionsCommand)
		return
	}

	switch session.State {
	case states.AwaitCollectionSortingDirection:
		parseSortingCollectionsDirection(app, session)
	}
}

// handleSortingCollectionsDirection prompts the user to choose a sorting direction (ascending or descending).
func handleSortingCollectionsDirection(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestSortDirection(session, session.CollectionsState.Sorting), keyboards.SortingDirection(session, session.CollectionsState.Sorting))
	session.SetState(states.AwaitCollectionSortingDirection)
}

// parseSortingCollectionsDirection processes the user's choice of sorting direction.
func parseSortingCollectionsDirection(app models.App, session *models.Session) {
	if utils.IsReset(app.Update) {
		handleSortingCollectionsReset(app, session, HandleSortingCollectionsCommand)
		return
	}

	if utils.IsDecrease(app.Update) {
		// Sets the sorting direction to descending ("-").
		session.CollectionsState.Sorting.Direction = "-"
	}

	// Applies the sorting settings.
	session.CollectionsState.Sorting.SetSort()

	app.SendMessage(messages.SortingApplied(session, session.CollectionsState.Sorting), nil)
	handleWithResetCollectionsPage(app, session, HandleCollectionsCommand)
}

// handleSortingCollectionsReset resets the sorting settings and reloads the collections list.
func handleSortingCollectionsReset(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.CollectionsState.Sorting.Reset()
	app.SendMessage(messages.ResetSortingSuccess(session), nil)
	handleWithResetCollectionsPage(app, session, next)
}

// handleWithResetCollectionsPage resets the current page and delegates to the next handler.
func handleWithResetCollectionsPage(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.ClearAllStates()
	session.CollectionsState.CurrentPage = 1
	next(app, session)
}
