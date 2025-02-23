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
	{"👤", "profile", states.CallbackMenuSelectProfile, "", true},
	{"🎥", "films", states.CallbackMenuSelectFilms, "", true},
	{"📚", "collections", states.CallbackMenuSelectCollections, "", true},
	{"⚙️", "settings", states.CallbackMenuSelectSettings, "", true},
	{"💬", "feedback", states.CallbackMenuSelectFeedback, "", true},
	{"🚪", "logout", states.CallbackMenuSelectLogout, "", true},
}

var settingsButtons = []Button{
	{"🈳", "language", states.CallbackSettingsLanguage, "", true},
	{"🌐", "kinopoiskToken", states.CallbackSettingsKinopoiskToken, "", true},
	{"🔢", "collectionsPageSize", states.CallbackSettingsCollectionsPageSize, "", true},
	{"🔢", "filmsPageSize", states.CallbackSettingsFilmsPageSize, "", true},
	{"🔢", "objectsPageSize", states.CallbackSettingsObjectsPageSize, "", true},
}

var feedbackCategoryButtons = []Button{
	{"💡", "offers", states.CallbackFeedbackCategorySuggestions, "", true},
	{"🐞", "mistakes", states.CallbackFeedbackCategoryBugs, "", true},
	{"❓", "otherIssues", states.CallbackFeedbackCategoryOther, "", true},
}

func BuildKeyboardWithCancel(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return NewKeyboard().AddCancel().Build(session.Lang)
}

func BuildKeyboardWithBack(session *models.Session, callback string) *tgbotapi.InlineKeyboardMarkup {
	return NewKeyboard().AddBack(callback).Build(session.Lang)
}

func BuildKeyboardWithSurvey(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return NewKeyboard().AddSurvey().Build(session.Lang)
}

func BuildLanguageSelectKeyboard(languages []string) *tgbotapi.InlineKeyboardMarkup {
	return NewKeyboard().AddLanguageSelect(languages, states.PrefixSelectStartLang).Build("")
}

func BuildMenuKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	if session.Role.HasAccess(roles.Helper) {
		keyboard.AddButton("🛠️", "adminPanel", states.CallbackMenuSelectAdmin, "", true)
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

func (k *Keyboard) AddLanguageSelect(languages []string, callback string) *Keyboard {
	var buttons []Button

	for _, lang := range languages {
		buttons = append(buttons, Button{"", fmt.Sprintf(strings.ToUpper(lang)), callback + lang, "", false})
	}

	k.AddButtonsWithRowSize(2, buttons...)

	return k
}
