package users

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleProfileCommand(app models.App, session *models.Session) {
	user, err := watchlist.GetUser(app, session)
	if err != nil {
		app.SendMessage(err.Error(), nil)
		return
	}

	msg := fmt.Sprintf("Вот ваш профиль, %s:\n", user.Username) +
		fmt.Sprintf("Ваш ID в системе API: %d\n", user.ID) +
		fmt.Sprintf("Ваш email: %s\n", user.Email) +
		fmt.Sprintf("Аккаунт был создан %v", user.CreatedAt)

	keyboard := builders.NewKeyboard(1).AddBack(states.CallbackProfileBack).Build()

	app.SendMessage(msg, keyboard)
}

func HandleProfileButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackProfileBack:
		general.HandleMenuCommand(app, session)
	}
}
