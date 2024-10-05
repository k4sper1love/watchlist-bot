package handlers

import (
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/pkg/utils"
	"strings"
)

func HandleUserInput(app config.App) {
	session, err := postgres.GetSessionByTelegramID(app.Upd.Message.From.ID)
	if err != nil {
		utils.SendMessage(app.Bot, app.Upd, "Произошла ошибка при получении сессии")
		return
	}

	if session.State == "" {
		switch app.Upd.Message.Command() {
		case "start":
			handleStartCommand(app, session)
		case "help":
			handleHelpCommand(app, session)
		case "profile":
			requireAuth(app, session, handleProfileCommand)
		case "register":
			requireNoAuth(app, session, startRegistrationProcess)
		case "login":
			requireNoAuth(app, session, startLoginProcess)
		case "logout":
			requireAuth(app, session, startLogoutProcess)
		default:
			utils.SendMessage(app.Bot, app.Upd, "Введите /register")
		}
	} else {
		switch {
		case strings.HasPrefix(session.State, "registration_"):
			handleRegistrationProcess(app, session)
		case strings.HasPrefix(session.State, "login_"):
			handleLoginProcess(app, session)
		case strings.HasPrefix(session.State, "logout_"):
			requireAuth(app, session, handleLogoutProcess)
		default:
			utils.SendMessage(app.Bot, app.Upd, "Неизвестное состояние")
		}
	}

	postgres.Save(&session)
}
