package general

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func HandleMenuCommand(app models.App, session *models.Session) {
	keyboard := keyboards.BuildMenuKeyboard(session.IsAdmin)

	app.SendMessage("📋 <b>Главное меню:</b>\n\nВыберите одно из действий", keyboard)
}
