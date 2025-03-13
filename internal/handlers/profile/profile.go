package profile

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleProfileCommand(app models.App, session *models.Session) {
	if user, err := watchlist.GetUser(app, session); err != nil {
		app.SendMessage(err.Error(), nil)
	} else {
		session.User = *user
		app.SendMessage(messages.Profile(session), keyboards.Profile(session))
	}
}

func HandleProfileButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallProfileUpdate:
		HandleUpdateProfileCommand(app, session)

	case states.CallProfileDelete:
		HandleDeleteProfileCommand(app, session)
	}
}
