package handlers

import (
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/utils"
)

func handleStartCommand(app config.App, session *models.Session) {
	msg := ""

	if !session.IsLogged {
		msg += "Добро пожаловать в Watchlist Bot. \n" +
			"На данный момент, бот находится в разработке.\n" +
			"Используйте /help, чтобы получить список доступных вам команд"
	} else {
		msg += "С возвращением в Watchlist Bot.\n" +
			"Используйте /help, чтобы получить список доступных вам команд"
	}

	utils.SendMessage(app.Bot, app.Upd, msg)
}

func handleHelpCommand(app config.App, session *models.Session) {
	msg := "Список доступных вам команд:\n" +
		"/start - привественное сообщение от бота\n" +
		"/help - список доступных команд\n"

	if !session.IsLogged {
		msg += "/register - создать новый аккаунт\n" +
			"/login - войти в существующий аккаунт"
	} else {
		msg += "/profile - получить информацию об аккаунте\n" +
			"/logout - выйти из системы"
	}

	utils.SendMessage(app.Bot, app.Upd, msg)
}
