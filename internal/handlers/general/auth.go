package general

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleAuthProcess(app models.App, session *models.Session) error {
	if session.AccessToken != "" {
		if watchlist.IsTokenValid(app, session.AccessToken) {
			return nil
		} else if session.RefreshToken != "" {
			if err := watchlist.RefreshAccessToken(app, session); err == nil {
				return nil
			}
		}
	}

	err := watchlist.Login(app, session)
	if err == nil {
		//msg := translator.Translate(session.Lang, "loginSuccess", map[string]interface{}{
		//  "Username": session.User.Username,
		//}, nil)
		//
		//app.SendMessage(msg, nil)
		return nil
	}

	err = watchlist.Register(app, session)
	if err == nil {
		msg := "✅ " + translator.Translate(session.Lang, "registrationSuccess", map[string]interface{}{
			"Username": session.User.Username,
		}, nil)

		app.SendMessage(msg, nil)
		return nil
	}

	return err
}

func HandleLogoutCommand(app models.App, session *models.Session) {
	msg := "⚠️ " + translator.Translate(session.Lang, "logoutConfirm", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)

	keyboard := keyboards.NewKeyboard().AddSurvey().Build(session.Lang)

	app.SendMessage(msg, keyboard)
	session.SetState(states.ProcessLogoutAwaitingConfirm)
}

func HandleLogoutProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.ProcessLogoutAwaitingConfirm:
		parseLogoutConfirm(app, session)
	}
}

func parseLogoutConfirm(app models.App, session *models.Session) {
	switch utils.IsAgree(app.Update) {
	case true:
		username := session.User.Username

		if err := watchlist.Logout(app, session); err != nil {
			msg := "🚨 " + translator.Translate(session.Lang, "logoutFailure", map[string]interface{}{
				"Username": username,
			}, nil)
			keyboard := keyboards.NewKeyboard().AddBack("").Build(session.Lang)
			app.SendMessage(msg, keyboard)
			break
		}

		msg := "🚪 " + translator.Translate(session.Lang, "logoutSuccess", map[string]interface{}{
			"Username": username,
		}, nil)

		app.SendMessage(msg, nil)
		session.Logout()

	case false:
		msg := "🚫 " + translator.Translate(session.Lang, "cancelAction", nil, nil)
		app.SendMessage(msg, nil)
	}

	session.ClearState()
	HandleMenuCommand(app, session)
}
