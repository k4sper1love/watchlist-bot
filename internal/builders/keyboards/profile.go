package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

var updateProfileButtons = []Button{
	{"", "name", states.CallbackUpdateProfileSelectUsername, "", true},
	{"", "email", states.CallbackUpdateProfileSelectEmail, "", true},
}

func Profile(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddProfileUpdate().
		AddDelete(states.CallbackProfileSelectDelete).
		AddBack("").
		Build(session.Lang)
}

func UpdateProfile(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddButtons(updateProfileButtons...).
		AddBack(states.CallbackUpdateProfileSelectBack).
		Build(session.Lang)
}
