package users

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

var updateProfileButtons = []keyboards.Button{
	{"", "Имя", states.CallbackUpdateProfileSelectUsername},
	{"", "Email", states.CallbackUpdateProfileSelectEmail},
}

func HandleUpdateProfileCommand(app models.App, session *models.Session) {
	msg := messages.BuildProfileMessage(session)
	msg += translator.Translate(session.Lang, "updateChoiceField", nil, nil)

	keyboard := keyboards.NewKeyboard().
		AddButtons(updateProfileButtons...).
		AddBack(states.CallbackUpdateProfileSelectBack).
		Build(session.Lang)

	app.SendMessage(msg, keyboard)
}

func HandleUpdateProfileButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackUpdateProfileSelectBack:
		HandleProfileCommand(app, session)
	case states.CallbackUpdateProfileSelectUsername:
		handleUpdateProfileUsername(app, session)
	case states.CallbackUpdateProfileSelectEmail:
		handleUpdateProfileEmail(app, session)
	}

}

func HandleUpdateProfileProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
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
	msg := translator.Translate(session.Lang, "updateProfileUsername", nil, nil)

	keyboard := keyboards.NewKeyboard().
		AddCancel().
		Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateProfileAwaitingUsername)
}

func parseUpdateProfileUsername(app models.App, session *models.Session) {
	session.ProfileState.Username = utils.ParseMessageString(app.Upd)

	finishUpdateProfileProcess(app, session)
}

func handleUpdateProfileEmail(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "updateProfileEmail", nil, nil)

	keyboard := keyboards.NewKeyboard().
		AddCancel().
		Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateProfileAwaitingEmail)
}

func parseUpdateProfileEmail(app models.App, session *models.Session) {
	session.ProfileState.Email = utils.ParseMessageString(app.Upd)

	finishUpdateProfileProcess(app, session)
}

func updateProfile(app models.App, session *models.Session) {
	user, err := watchlist.UpdateUser(app, session)
	if err != nil {
		msg := translator.Translate(session.Lang, "updateProfileFailure", map[string]interface{}{
			"Username": session.User.Username,
		}, nil)

		app.SendMessage(msg, nil)
		return
	}
	session.User = *user

	msg := translator.Translate(session.Lang, "updateProfileSuccess", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)

	app.SendMessage(msg, nil)
}

func finishUpdateProfileProcess(app models.App, session *models.Session) {
	state := session.ProfileState
	user := session.User

	if state.Username == "" {
		state.Username = user.Username
	}

	if state.Email == "" {
		state.Email = user.Email
	}

	updateProfile(app, session)
	session.ProfileState.Clear()
	session.ClearState()
	HandleUpdateProfileCommand(app, session)
}
