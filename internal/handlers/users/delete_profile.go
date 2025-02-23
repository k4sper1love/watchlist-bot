package users

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleDeleteProfileCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildDeleteProfileMessage(session), keyboards.BuildKeyboardWithSurvey(session))
	session.SetState(states.ProcessDeleteProfileAwaitingConfirm)
}

func HandleDeleteProfileProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.ProcessDeleteProfileAwaitingConfirm:
		parseDeleteProfileConfirm(app, session)
	}
}

func parseDeleteProfileConfirm(app models.App, session *models.Session) {
	if !utils.IsAgree(app.Update) {
		app.SendMessage(messages.BuildCancelActionMessage(session), nil)
		session.ClearState()
		HandleProfileCommand(app, session)
		return
	}

	if err := watchlist.DeleteUser(app, session); err != nil {
		app.SendMessage(messages.BuildDeleteProfileFailureMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackMenuSelectProfile))
		session.ClearState()
		return
	}

	app.SendMessage(messages.BuildDeleteProfileSuccessMessage(session), nil)
	session.Logout()
}
