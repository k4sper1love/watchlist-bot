package profile

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleDeleteProfileCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.DeleteProfile(session), keyboards.Survey(session))
	session.SetState(states.AwaitDeleteProfileConfirm)
}

func HandleDeleteProfileProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.AwaitDeleteProfileConfirm:
		handleDeleteProfileConfirm(app, session)
		session.ClearState()
	}
}

func handleDeleteProfileConfirm(app models.App, session *models.Session) {
	if !utils.IsAgree(app.Update) {
		app.SendMessage(messages.CancelAction(session), nil)
		HandleProfileCommand(app, session)
		return
	}

	if err := watchlist.DeleteUser(app, session); err != nil {
		app.SendMessage(messages.DeleteProfileFailure(session), keyboards.Back(session, states.CallMenuProfile))
		return
	}

	app.SendMessage(messages.DeleteProfileSuccess(session), nil)
	session.Logout()
}
