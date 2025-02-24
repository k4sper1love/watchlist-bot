package users

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
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

	session.User = *user
	app.SendMessage(messages.BuildProfileMessage(session), keyboards.BuildProfileKeyboard(session))
}

func HandleProfileButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallbackProfileSelectUpdate:
		HandleUpdateProfileCommand(app, session)
	case states.CallbackProfileSelectDelete:
		HandleDeleteProfileCommand(app, session)
	}
}
