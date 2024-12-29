package general

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"log/slog"
)

func HandleStartCommand(app models.App, session *models.Session) {
	msg := messages.BuildStartMessage(app, session)
	app.SendMessage(msg, nil)

	if err := HandleAuthProcess(app, session); err != nil {
		msg = translator.Translate(session.Lang, "authFailure", nil, nil)
		app.SendMessage(msg, nil)
		sl.Log.Error("error auth process", slog.Any("err", err))
		return
	}

	HandleMenuCommand(app, session)
}

func HandleHelpCommand(app models.App, session *models.Session) {
	msg := messages.BuildHelpMessage(session)

	app.SendMessage(msg, nil)
}
