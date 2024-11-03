package general

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

var menuButtons = []builders.Button{
	{"Профиль", states.CallbackMenuSelectProfile},
	{"Коллекции", states.CallbackMenuSelectCollections},
	{"Настройки", states.CallbackMenuSelectSettings},
	{"Выйти из системы", states.CallbackMenuSelectLogout},
}

func HandleMenuCommand(app models.App, session *models.Session) {
	keyboard := builders.NewKeyboard(1).AddSeveral(menuButtons).Build()

	app.SendMessage("Выберите одно из действий", keyboard)
}
