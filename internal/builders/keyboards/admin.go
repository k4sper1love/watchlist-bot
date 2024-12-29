package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

var superAdminButtons = []Button{
	{"", "Администраторы", ""},
}

var adminButtons = []Button{
	{"", "Пользователи", ""},
	{"", "Рассылка", ""},
}

var helperButtons = []Button{
	{"", "Фидбек", ""},
	{"", "adminOptionUserCount", states.CallbackAdminSelectUserCount},
	{"", "adminOptionBroadcast", states.CallbackAdminSelectBroadcastMessage},
	{"", "adminOptionFeedback", states.CallbackAdminSelectFeedback},
	{"", "adminOptionUserList", states.CallbackAdminSelectUsers},
}

func BuildManageAdminKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddButtons(adminButtons...)

	return keyboard.Build(session.Lang)
}
