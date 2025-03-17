package profile

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// HandleProfileCommand handles the command for displaying the user's profile.
// Fetches the user's data from the Watchlist service and sends a message with their profile details.
func HandleProfileCommand(app models.App, session *models.Session) {
	if user, err := watchlist.GetUser(app, session); err != nil {
		app.SendMessage(err.Error(), nil)
	} else {
		session.User = *user
		app.SendMessage(messages.Profile(session), keyboards.Profile(session))
	}
}

// HandleProfileButtons handles button interactions related to the user's profile.
// Supports actions like updating the profile or deleting it.
func HandleProfileButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallProfileUpdate:
		HandleUpdateProfileCommand(app, session)

	case states.CallProfileDelete:
		HandleDeleteProfileCommand(app, session)
	}
}
