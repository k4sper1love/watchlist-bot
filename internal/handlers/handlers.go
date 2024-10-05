package handlers

import (
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/pkg/utils"
	"strings"
)

func HandleUserInput(app *config.App) {
	session, err := postgres.GetSessionByTelegramID(app.Upd.Message.From.ID)
	if err != nil {
		utils.SendMessage(app.Bot, app.Upd, "Произошла ошибка при получении пользователя")
		return
	}

	if session.State == "" {
		switch app.Upd.Message.Command() {
		case "profile":
			handleUserGet(app, session)
		case "register":
			startRegistrationProcess(app, session)
		default:
			utils.SendMessage(app.Bot, app.Upd, "Введите /register")
		}
	} else {
		switch {
		case strings.HasPrefix(session.State, "registration_"):
			handleRegistrationProcess(app, session)
		default:
			utils.SendMessage(app.Bot, app.Upd, "Неизвестное состояние")
		}
	}

	postgres.Save(&session)
}
