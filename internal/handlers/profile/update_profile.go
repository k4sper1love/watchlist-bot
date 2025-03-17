package profile

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// HandleUpdateProfileCommand handles the command for updating the user's profile.
// Sends a message with options to update the username or email.
func HandleUpdateProfileCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.UpdateProfile(session), keyboards.UpdateProfile(session))
}

// HandleUpdateProfileButtons handles button interactions related to updating the user's profile.
// Supports actions like going back, updating the username, or updating the email.
func HandleUpdateProfileButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallUpdateProfileBack:
		HandleProfileCommand(app, session)

	case states.CallUpdateProfileUsername:
		requestUpdateProfileUsername(app, session)

	case states.CallUpdateProfileEmail:
		requestUpdateProfileEmail(app, session)
	}
}

// HandleUpdateProfileProcess processes the workflow for updating the user's profile.
// Handles states like awaiting input for the username or email.
func HandleUpdateProfileProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearState()
		HandleUpdateProfileCommand(app, session)
		return
	}

	switch session.State {
	case states.AwaitUpdateProfileUsername:
		parser.ParseProfileUsername(app, session, requestUpdateProfileUsername, finishUpdateProfile)

	case states.AwaitUpdateProfileEmail:
		parser.ParseProfileEmail(app, session, requestUpdateProfileEmail, finishUpdateProfile)
	}
}

// requestUpdateProfileUsername prompts the user to enter a new username for their profile.
func requestUpdateProfileUsername(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestProfileUsername(session), keyboards.Cancel(session))
	session.SetState(states.AwaitUpdateProfileUsername)
}

// requestUpdateProfileEmail prompts the user to enter a new email for their profile.
func requestUpdateProfileEmail(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestProfileEmail(session), keyboards.Cancel(session))
	session.SetState(states.AwaitUpdateProfileEmail)
}

// finishUpdateProfile finalizes the process of updating the user's profile.
// Calls the Watchlist service to update the profile and navigates back to the update menu.
func finishUpdateProfile(app models.App, session *models.Session) {
	if err := updateProfile(app, session); err != nil {
		app.SendMessage(messages.UpdateProfileFailure(session), nil)
	} else {
		app.SendMessage(messages.UpdateProfileSuccess(session), nil)
	}

	session.ClearAllStates()
	HandleUpdateProfileCommand(app, session)
}

// updateProfile updates the user's profile using the Watchlist service.
func updateProfile(app models.App, session *models.Session) error {
	user, err := watchlist.UpdateUser(app, session)
	if err != nil {
		return err
	}

	session.User = *user
	return nil
}
