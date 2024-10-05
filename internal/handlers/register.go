package handlers

import (
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/pkg/utils"
)

func startRegistrationProcess(app *config.App, session *models.Session) {
	utils.SendMessage(app.Bot, app.Upd, "Введите желаемый юзернейм")
	session.State = RegistrationStateAwatingUsername
}

func handleRegistrationProcess(app *config.App, session *models.Session) {
	switch session.State {
	case RegistrationStateAwatingUsername:
		session.AuthState.Username = app.Upd.Message.Text
		utils.SendMessage(app.Bot, app.Upd, "Введите ваш email")
		session.State = RegistrationStateAwatingEmail
	case RegistrationStateAwatingEmail:
		session.AuthState.Email = app.Upd.Message.Text
		utils.SendMessage(app.Bot, app.Upd, "Введите ваш пароль")
		session.State = RegistrationStateAwatingPassword
	case RegistrationStateAwatingPassword:
		session.AuthState.Password = app.Upd.Message.Text
		if err := watchlist.RegisterUser(app, session); err != nil {
			utils.SendMessage(app.Bot, app.Upd, "Неудачная регистрация")
		} else {
			utils.SendMessage(app.Bot, app.Upd, "Успешная регистрация")
			session.IsLogged = true
		}
		session.AuthState.Password = ""
		session.State = ""
	}
}
