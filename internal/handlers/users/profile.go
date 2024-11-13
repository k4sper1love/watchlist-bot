package users

import (
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

	msg := builders.BuildProfileMessage(user)

	keyboard := builders.NewKeyboard(1).
		AddProfileUpdate().
		AddProfileDelete().
		AddBack(states.CallbackProfileSelectBack).
		Build()

	app.SendMessage(msg, keyboard)
}

func HandleProfileButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)

	switch {
	case callback == states.CallbackProfileSelectBack:
		general.HandleMenuCommand(app, session)

	case callback == states.CallbackProfileSelectUpdate:
		HandleUpdateProfileCommand(app, session)

	case callback == states.CallbackProfileSelectDelete:
		HandleDeleteProfileCommand(app, session)
	}
}
