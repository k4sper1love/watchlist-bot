package general

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// HandleLogoutCommand handles the command for logging out the user.
// Sends a confirmation message and sets the session state to await user confirmation.
func HandleLogoutCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.Logout(session), keyboards.Survey(session))
	session.SetState(states.AwaitLogoutConfirm)
}

// HandleLogoutProcess processes the workflow for logging out the user.
// Handles states like awaiting confirmation from the user.
func HandleLogoutProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.AwaitLogoutConfirm:
		handleLogoutConfirm(app, session)
	}
}

// handleLogoutConfirm processes the user's response to the logout confirmation.
func handleLogoutConfirm(app models.App, session *models.Session) {
	if !utils.IsAgree(app.Update) {
		app.SendMessage(messages.CancelAction(session), nil)
		session.ClearState()
		HandleMenuCommand(app, session)
		return
	}

	if err := watchlist.Logout(app, session); err != nil {
		app.SendMessage(messages.LogoutFailure(session), keyboards.Back(session, ""))
		session.ClearState()
		return
	}

	app.SendMessage(messages.LogoutSuccess(session), nil)
	session.Logout()
}
