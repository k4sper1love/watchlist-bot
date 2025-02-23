package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

var updateProfileButtons = []Button{
	{"", "Имя", states.CallbackUpdateProfileSelectUsername, "", true},
	{"", "Email", states.CallbackUpdateProfileSelectEmail, "", true},
}

func BuildProfileKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddProfileUpdate()

	keyboard.AddDelete(states.CallbackProfileSelectDelete)

	keyboard.AddBack("")

	return keyboard.Build(session.Lang)
}

func BuildUpdateProfileKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return NewKeyboard().
		AddButtons(updateProfileButtons...).
		AddBack(states.CallbackUpdateProfileSelectBack).
		Build(session.Lang)
}

func (k *Keyboard) AddProfileUpdate() *Keyboard {
	return k.AddButton("✏️", "edit", states.CallbackProfileSelectUpdate, "", true)
}

func (k *Keyboard) AddProfileDelete() *Keyboard {
	return k.AddButton("🗑️", "delete", states.CallbackProfileSelectDelete, "", true)
}
