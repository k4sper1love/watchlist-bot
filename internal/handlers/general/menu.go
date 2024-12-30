package general

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"log"
)

func HandleMenuCommand(app models.App, session *models.Session) {
	log.Println(session.Lang)
	keyboard := keyboards.BuildMenuKeyboard(session)

	msg := messages.BuildMenuMessage(session)

	app.SendMessage(msg, keyboard)
}
