package general

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

var menuButtons = []builders.Button{
	{"Профиль", states.CallbackMenuSelectProfile},
	{"Фильмы", states.CallbackMenuSelectFilms},
	{"Коллекции", states.CallbackMenuSelectCollections},
	{"Настройки", states.CallbackMenuSelectSettings},
	{"Выйти из системы", states.CallbackMenuSelectLogout},
}

func HandleMenuCommand(app models.App, session *models.Session) {
	keyboard := builders.NewKeyboard(1)

	if session.IsAdmin {
		adminButton := builders.Button{"Админ-панель", states.CallbackMenuSelectAdmin}

		keyboard.Add(adminButton)
	}

	keyboard.AddSeveral(menuButtons)

	app.SendMessage("Выберите одно из действий", keyboard.Build())
}
