package users

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleUpdateProfileCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.UpdateProfile(session), keyboards.UpdateProfile(session))
}

func HandleUpdateProfileButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallbackUpdateProfileSelectBack:
		HandleProfileCommand(app, session)
	case states.CallbackUpdateProfileSelectUsername:
		handleUpdateProfileUsername(app, session)
	case states.CallbackUpdateProfileSelectEmail:
		handleUpdateProfileEmail(app, session)
	}
}

func HandleUpdateProfileProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearState()
		HandleUpdateProfileCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessUpdateProfileAwaitingUsername:
		parseUpdateProfileUsername(app, session)
	case states.ProcessUpdateProfileAwaitingEmail:
		parseUpdateProfileEmail(app, session)
	}
}

func handleUpdateProfileUsername(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestProfileUsername(session), keyboards.Cancel(session))
	session.SetState(states.ProcessUpdateProfileAwaitingUsername)
}

func parseUpdateProfileUsername(app models.App, session *models.Session) {
	session.ProfileState.Username = utils.ParseMessageString(app.Update)
	finishUpdateProfileProcess(app, session)
}

func handleUpdateProfileEmail(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestProfileEmail(session), keyboards.Cancel(session))
	session.SetState(states.ProcessUpdateProfileAwaitingEmail)
}

func parseUpdateProfileEmail(app models.App, session *models.Session) {
	session.ProfileState.Email = utils.ParseMessageString(app.Update)
	finishUpdateProfileProcess(app, session)
}

func finishUpdateProfileProcess(app models.App, session *models.Session) {
	if session.ProfileState.Username == "" {
		session.ProfileState.Username = session.User.Username
	}
	if session.ProfileState.Email == "" {
		session.ProfileState.Email = session.User.Email
	}

	if err := updateProfile(app, session); err == nil {
		HandleUpdateProfileCommand(app, session)
	}

	session.ClearAllStates()
}

func updateProfile(app models.App, session *models.Session) error {
	user, err := watchlist.UpdateUser(app, session)
	if err != nil {
		app.SendMessage(messages.UpdateProfileFailure(session), keyboards.Back(session, states.CallbackProfileSelectUpdate))
		return err
	}

	session.User = *user
	app.SendMessage(messages.UpdateProfileSuccess(session), nil)
	return nil
}
