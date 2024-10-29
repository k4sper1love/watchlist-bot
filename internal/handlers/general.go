package handlers

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"log/slog"
)

func handleStartCommand(app config.App, session *models.Session) {
	msg := "Добро пожаловать в Watchlist Bot. \n" +
		"На данный момент, бот находится в разработке.\n" +
		"Используйте /help, чтобы получить список доступных вам команд\n"

	sendMessage(app, msg)

	msg = "Выполняем вход в систему"
	sendMessage(app, msg)

	if err := handleAuthProcess(app, session); err != nil {
		sendMessage(app, "Не удалось войти в систему. Отправьте /start заново.")
		sl.Log.Error("error auth process", slog.Any("err", err))
		return
	}
	sendMessage(app, "Успешный вход в систему.")
	handleHelpCommand(app, session)
}

func handleHelpCommand(app config.App, session *models.Session) {
	msg := "Список доступных вам команд:\n" +
		"/start - начать пользоваться ботом\n" +
		"/help - список доступных команд\n" +
		"/profile - получить информацию об аккаунте\n" +
		"/collections - список ваших коллекций\n" +
		"/new_collection - создать новую коллекцию\n" +
		"/settings - настройки\n" +
		"/logout - выйти из системы\n"

	sendMessage(app, msg)
}
