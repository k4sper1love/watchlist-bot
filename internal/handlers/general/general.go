package general

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"log/slog"
)

func HandleStartCommand(app models.App, session *models.Session) {
	msg := "Добро пожаловать в Watchlist Bot. \n" +
		"На данный момент, бот находится в разработке.\n" +
		"Используйте /help, чтобы получить список доступных вам команд\n"

	app.SendMessage(msg, nil)

	msg = "Выполняем вход в систему"
	app.SendMessage(msg, nil)

	if err := HandleAuthProcess(app, session); err != nil {
		app.SendMessage("Не удалось войти в систему. Отправьте /start заново.", nil)
		sl.Log.Error("error auth process", slog.Any("err", err))
		return
	}
	HandleMenuCommand(app, session)
}

func HandleHelpCommand(app models.App, session *models.Session) {
	msg := "Список доступных вам команд:\n" +
		"/start - начать пользоваться ботом\n" +
		"/help - список доступных команд\n" +
		"/profile - получить информацию об аккаунте\n" +
		"/collections - список ваших коллекций\n" +
		"/settings - настройки\n" +
		"/logout - выйти из системы\n"

	app.SendMessage(msg, nil)
}
