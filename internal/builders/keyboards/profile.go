package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

// Predefined buttons for updating profile details.
var updateProfileButtons = []Button{
	{"", "name", states.CallUpdateProfileUsername, "", true},
	{"", "email", states.CallUpdateProfileEmail, "", true},
}

// Profile creates an inline keyboard for managing the user's profile.
func Profile(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddProfileUpdate().
		AddDelete(states.CallProfileDelete).
		AddBack("").
		Build(session.Lang)
}

// UpdateProfile creates an inline keyboard for updating specific profile fields (e.g., name, email).
func UpdateProfile(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddButtons(updateProfileButtons...).
		AddBack(states.CallUpdateProfileBack).
		Build(session.Lang)
}
