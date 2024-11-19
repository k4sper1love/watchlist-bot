package general

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"log/slog"
)

func HandleStartCommand(app models.App, session *models.Session) {
	msg := messages.BuildStartMessage()
	app.SendMessage(msg, nil)

	if err := HandleAuthProcess(app, session); err != nil {
		app.SendMessage("Не удалось войти в систему. Отправьте /start заново.", nil)
		sl.Log.Error("error auth process", slog.Any("err", err))
		return
	}

	HandleMenuCommand(app, session)
}

func HandleHelpCommand(app models.App, session *models.Session) {
	msg := messages.BuildHelpMessage()

	app.SendMessage(msg, nil)
}
