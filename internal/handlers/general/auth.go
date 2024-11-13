package general

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
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
		app.SendMessage("Успешный вход!", nil)
		return nil
	}

	err = watchlist.Register(app, session)
	if err == nil {
		app.SendMessage("Успешная регистрация", nil)
		return nil
	}

	return err
}

func HandleLogoutCommand(app models.App, session *models.Session) {
	msg := "Вы точно хотите выйти из аккаунта?"

	keyboard := builders.NewKeyboard(1).AddSurvey().Build()

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
	switch utils.IsAgree(app.Upd) {
	case true:
		if err := watchlist.Logout(app, session); err != nil {
			app.SendMessage("Неудачный выход из системы", nil)
			break
		}
		app.SendMessage("Успешно вышли из системы", nil)
		session.ClearFull()

	case false:
		app.SendMessage("Отмена выхода из системы", nil)
	}

	session.ClearState()
	HandleMenuCommand(app, session)
}
