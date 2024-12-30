package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"strings"
)

var menuButtons = []Button{
	{"👤", "profile", states.CallbackMenuSelectProfile},
	{"🎥", "films", states.CallbackMenuSelectFilms},
	{"📚", "collections", states.CallbackMenuSelectCollections},
	{"⚙️", "settings", states.CallbackMenuSelectSettings},
	{"💬", "feedback", states.CallbackMenuSelectFeedback},
	{"🚪", "logout", states.CallbackMenuSelectLogout},
}

var settingsButtons = []Button{
	{"🈳", "settingsLanguage", states.CallbackSettingsLanguage},
	{"🔢", "settingsCollectionsPageSize", states.CallbackSettingsCollectionsPageSize},
	{"🔢", "settingsFilmsPageSize", states.CallbackSettingsFilmsPageSize},
	{"🔢", "settingsObjectsPageSize", states.CallbackSettingsObjectsPageSize},
}

var feedbackCategoryButtons = []Button{
	{"💡", "offers", states.CallbackFeedbackCategorySuggestions},
	{"🐞", "mistakes", states.CallbackFeedbackCategoryBugs},
	{"❓", "otherIssues", states.CallbackFeedbackCategoryOther},
}

func BuildMenuKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	if session.Role.HasAccess(roles.Helper) {
		keyboard.AddButton("🛠️", "adminPanel", states.CallbackMenuSelectAdmin)
	}

	keyboard.AddButtons(menuButtons...)

	return keyboard.Build(session.Lang)
}

func BuildSettingsKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddButtons(settingsButtons...)

	keyboard.AddBack("")

	return keyboard.Build(session.Lang)
}

func BuildFeedbackKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddButtons(feedbackCategoryButtons...)

	keyboard.AddBack("")

	return keyboard.Build(session.Lang)
}

func (k *Keyboard) AddProfileUpdate() *Keyboard {
	return k.AddButton("✏️", "edit", states.CallbackProfileSelectUpdate)
}

func (k *Keyboard) AddProfileDelete() *Keyboard {
	return k.AddButton("⚠️", "delete", states.CallbackProfileSelectDelete)
}

func (k *Keyboard) AddLanguageSelect(languages []string, callback string) *Keyboard {
	var buttons []Button

	for _, lang := range languages {
		buttons = append(buttons, Button{"", fmt.Sprintf(strings.ToUpper(lang)), fmt.Sprintf("%s_%s", callback, lang)})
	}

	k.AddButtonsWithRowSize(2, buttons...)

	return k
}
