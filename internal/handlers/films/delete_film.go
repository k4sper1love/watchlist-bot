package films

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// HandleDeleteFilmCommand handles the command for deleting a film.
// Sends a confirmation message and sets the session state to await user confirmation.
func HandleDeleteFilmCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.DeleteFilm(session), keyboards.Survey(session))
	session.SetState(states.AwaitDeleteFilmConfirm)
}

// HandleDeleteFilmProcess processes the workflow for deleting a film.
// Handles states like awaiting confirmation from the user.
func HandleDeleteFilmProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.AwaitDeleteFilmConfirm:
		handleDeleteFilmConfirm(app, session)
		session.ClearState()
	}
}

// handleDeleteFilmConfirm processes the user's response to the deletion confirmation.
func handleDeleteFilmConfirm(app models.App, session *models.Session) {
	if !utils.IsAgree(app.Update) {
		app.SendMessage(messages.CancelAction(session), nil)
		HandleManageFilmCommand(app, session)
		return
	}

	if err := DeleteFilm(app, session); err != nil {
		app.SendMessage(messages.DeleteFilmFailure(session), keyboards.Back(session, states.CallFilmsManage))
		return
	}

	app.SendMessage(messages.DeleteFilmSuccess(session), nil)
	HandleFilmsCommand(app, session)
}

// DeleteFilm deletes a film based on the current session context (user or collection).
func DeleteFilm(app models.App, session *models.Session) error {
	switch session.Context {
	case states.CtxFilm:
		return watchlist.DeleteFilm(app, session)
	case states.CtxCollection:
		return watchlist.DeleteCollectionFilm(app, session)
	default:
		return fmt.Errorf("unsupported session context: %s", session.Context)
	}
}
