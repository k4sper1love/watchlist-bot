package users

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

var updateProfileButtons = []builders.Button{
	{"Username", states.CallbackUpdateProfileSelectUsername},
	{"Email", states.CallbackUpdateProfileSelectEmail},
}

func HandleUpdateProfileCommand(app models.App, session *models.Session) {
	msg := builders.BuildProfileMessage(&session.User)
	msg += "\nВыберите, какое поле вы хотите изменить?"

	keyboard := builders.NewKeyboard(1).
		AddSeveral(updateProfileButtons).
		AddBack(states.CallbackUpdateProfileSelectBack).
		Build()

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
	msg := "Введите новый username"

	keyboard := builders.NewKeyboard(1).
		AddCancel().
		Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateProfileAwaitingUsername)
}

func parseUpdateProfileUsername(app models.App, session *models.Session) {
	session.ProfileState.Username = utils.ParseMessageString(app.Upd)

	finishUpdateProfileProcess(app, session)
}

func handleUpdateProfileEmail(app models.App, session *models.Session) {
	msg := "Введите адрес электронной почты"

	keyboard := builders.NewKeyboard(1).
		AddCancel().
		Build()

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
		app.SendMessage("Не удалось обновить профиль", nil)
		return
	}

	session.User = *user
	app.SendMessage("Профиль успешно обновлен", nil)
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
