package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
)

var menuButtons = []Button{
	{"👤 Профиль", states.CallbackMenuSelectProfile},
	{"🎥 Фильмы", states.CallbackMenuSelectFilms},
	{"📚 Коллекции", states.CallbackMenuSelectCollections},
	{"⚙️ Настройки", states.CallbackMenuSelectSettings},
	{"💬 Обратная связь", states.CallbackMenuSelectFeedback},
	{"🚪 Выйти из системы", states.CallbackMenuSelectLogout},
}

var settingsButtons = []Button{
	{"🔢 Изменить количество коллекций на странице", states.CallbackSettingsCollectionsPageSize},
	{"🔢 Изменить количество фильмов на странице", states.CallbackSettingsFilmsPageSize},
	{"🔢 Изменить количество объектов на странице", states.CallbackSettingsObjectsPageSize},
}

var feedbackCategoryButtons = []Button{
	{"💡 Предложения", states.CallbackFeedbackCategorySuggestions},
	{"🐞 Ошибки", states.CallbackFeedbackCategoryBugs},
	{"❓ Другие вопросы", states.CallbackFeedbackCategoryOther},
}

func BuildMenuKeyboard(isAdmin bool) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	if isAdmin {
		keyboard.AddButton("🛠️ Админ-панель", states.CallbackMenuSelectAdmin)
	}

	keyboard.AddButtons(menuButtons...)

	return keyboard.Build()
}

func BuildSettingsKeyboard() *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddButtons(settingsButtons...)

	keyboard.AddBack("")

	return keyboard.Build()
}

func BuildFeedbackKeyboard() *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddButtons(feedbackCategoryButtons...)

	keyboard.AddBack("")

	return keyboard.Build()
}

func (k *Keyboard) AddProfileUpdate() *Keyboard {
	return k.AddButton("✏️ Редактировать", states.CallbackProfileSelectUpdate)
}

func (k *Keyboard) AddProfileDelete() *Keyboard {
	return k.AddButton("⚠️ Удалить", states.CallbackProfileSelectDelete)
}
