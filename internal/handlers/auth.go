package handlers

import (
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/pkg/utils"
)

func startRegistrationProcess(app config.App, session *models.Session) {
	utils.SendMessage(app.Bot, app.Upd, "Введите желаемый юзернейм")
	session.State = RegistrationStateAwatingUsername
}

func handleRegistrationProcess(app config.App, session *models.Session) {
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
		if err := watchlist.Register(app, session); err != nil {
			utils.SendMessage(app.Bot, app.Upd, "Неудачная регистрация")
		} else {
			utils.SendMessage(app.Bot, app.Upd, "Успешная регистрация")
			session.IsLogged = true
		}
		session.AuthState.Password = ""
		session.State = ""
	}
}

func startLoginProcess(app config.App, session *models.Session) {
	utils.SendMessage(app.Bot, app.Upd, "Введите ваш email")
	session.State = LoginStateAwaitingEmail
}

func handleLoginProcess(app config.App, session *models.Session) {
	switch session.State {
	case LoginStateAwaitingEmail:
		session.AuthState.Email = app.Upd.Message.Text
		utils.SendMessage(app.Bot, app.Upd, "Введите ваш пароль")
		session.State = LoginStateAwaitingPassword
	case LoginStateAwaitingPassword:
		session.AuthState.Password = app.Upd.Message.Text
		if err := watchlist.Login(app, session); err != nil {
			utils.SendMessage(app.Bot, app.Upd, "Неудачный вход")
		} else {
			utils.SendMessage(app.Bot, app.Upd, "Успешный вход")
			session.IsLogged = true
		}
		session.AuthState.Password = ""
		session.State = ""
	}
}

func startLogoutProcess(app config.App, session *models.Session) {
	utils.SendMessage(app.Bot, app.Upd, "Вы точно хотите из аккаунта? Отправьте \"+\" для выхода")
	session.State = LogoutStateAwaitingConfirm
}

func handleLogoutProcess(app config.App, session *models.Session) {
	switch session.State {
	case LogoutStateAwaitingConfirm:
		if app.Upd.Message.Text == "+" {
			if err := watchlist.Logout(app, session); err != nil {
				utils.SendMessage(app.Bot, app.Upd, "Неудачный выход из системы")
			} else {
				utils.SendMessage(app.Bot, app.Upd, "Успешно вышли из системы")
				session.IsLogged = false
			}
		} else {
			utils.SendMessage(app.Bot, app.Upd, "Отмена выхода из системы")
		}
		session.State = ""
	}
}
