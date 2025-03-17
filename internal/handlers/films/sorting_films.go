package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strings"
)

// HandleSortingFilmsCommand handles the command for sorting films.
// Sends a message with options to choose a sorting field.
func HandleSortingFilmsCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.ChoiceSorting(session), keyboards.FilmsSorting(session))
}

// HandleSortingFilmsButtons handles button interactions related to sorting films.
// Supports actions like going back, resetting sorting, or selecting a sorting field.
func HandleSortingFilmsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallFilmSortingBack:
		handleWithResetFilmsPage(app, session, HandleFilmsCommand)

	case states.CallFilmSortingAllReset:
		handleSortingFilmsReset(app, session, HandleFilmsCommand)

	default:
		handleSortingFilmsSelect(app, session, callback)
	}
}

// handleSortingFilmsSelect processes the selection of a sorting field.
// Updates the session state with the selected field and prompts the user to choose a sorting direction.
func handleSortingFilmsSelect(app models.App, session *models.Session, callback string) {
	if strings.HasPrefix(callback, states.FilmSortingSelect) {
		session.GetFilmSortingByCtx().Field = strings.TrimPrefix(callback, states.FilmSortingSelect)
		handleSortingFilmsDirection(app, session)
	}
}

// HandleSortingFilmsProcess processes the workflow for sorting films.
// Handles states like awaiting input for the sorting direction.
func HandleSortingFilmsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		handleWithResetFilmsPage(app, session, HandleSortingFilmsCommand)
		return
	}

	switch session.State {
	case states.AwaitFilmSortingDirection:
		parseSortingFilmsDirection(app, session)
	}
}

// handleSortingFilmsDirection prompts the user to choose a sorting direction (ascending or descending).
func handleSortingFilmsDirection(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestSortDirection(session, session.GetFilmSortingByCtx()), keyboards.SortingDirection(session, session.GetFilmSortingByCtx()))
	session.SetState(states.AwaitFilmSortingDirection)
}

// parseSortingFilmsDirection processes the user's choice of sorting direction.
// Updates the session state with the chosen direction and applies the sorting.
func parseSortingFilmsDirection(app models.App, session *models.Session) {
	if utils.IsReset(app.Update) {
		handleSortingFilmsReset(app, session, HandleSortingFilmsCommand)
		return
	}

	if utils.IsDecrease(app.Update) {
		// Sets the sorting direction to descending ("-").
		session.GetFilmSortingByCtx().Direction = "-"
	}

	// Applies the sorting settings.
	session.GetFilmSortingByCtx().SetSort()

	app.SendMessage(messages.SortingApplied(session, session.GetFilmSortingByCtx()), nil)
	handleWithResetFilmsPage(app, session, HandleFilmsCommand)
}

// handleSortingFilmsReset resets the sorting settings and reloads the films list.
func handleSortingFilmsReset(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.GetFilmSortingByCtx().Reset()
	app.SendMessage(messages.ResetSortingSuccess(session), nil)
	handleWithResetFilmsPage(app, session, next)
}

// handleWithResetFilmsPage resets the current page and delegates to the next handler.
func handleWithResetFilmsPage(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	clearStatesAndResetFilmsPage(session)
	next(app, session)
}
