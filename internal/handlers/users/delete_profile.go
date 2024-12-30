package users

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleDeleteProfileCommand(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "deleteProfileConfirm", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)

	keyboard := keyboards.NewKeyboard().
		AddSurvey().
		Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessDeleteProfileAwaitingConfirm)
}

func HandleDeleteProfileProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.ProcessDeleteProfileAwaitingConfirm:
		parseDeleteProfileConfirm(app, session)
	}
}

func parseDeleteProfileConfirm(app models.App, session *models.Session) {
	session.ClearState()

	switch utils.IsAgree(app.Upd) {
	case true:
		if err := watchlist.DeleteUser(app, session); err != nil {
			msg := translator.Translate(session.Lang, "deleteProfileFailure", map[string]interface{}{
				"Username": session.User.Username,
			}, nil)

			app.SendMessage(msg, nil)
			HandleProfileCommand(app, session)
			return
		}

		msg := translator.Translate(session.Lang, "deleteProfileSuccess", map[string]interface{}{
			"Username": session.User.Username,
		}, nil)

		app.SendMessage(msg, nil)
		session.Logout()
		general.HandleMenuCommand(app, session)

	case false:
		msg := translator.Translate(session.Lang, "cancelAction", nil, nil)
		app.SendMessage(msg, nil)
		HandleProfileCommand(app, session)
	}
}
