package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

// Predefined buttons for managing user settings.
var settingsButtons = []Button{
	{"🈳", "language", states.CallSettingsLanguage, "", true},
	{"🌐", "kinopoiskToken", states.CallSettingsKinopoiskToken, "", true},
	{"🔢", "collectionsPageSize", states.CallSettingsCollectionsPageSize, "", true},
	{"🔢", "filmsPageSize", states.CallSettingsFilmsPageSize, "", true},
	{"🔢", "objectsPageSize", states.CallSettingsObjectsPageSize, "", true},
}

// Settings creates an inline keyboard for managing user settings.
func Settings(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddButtons(settingsButtons...).
		AddBack("").
		Build(session.Lang)
}

// SettingsLanguageSelect creates an inline keyboard for selecting a language in the settings.
func SettingsLanguageSelect(session *models.Session, languages []string) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddLanguageSelect(languages, states.SelectLang).
		AddBack(states.CallSettingsBack).
		Build(session.Lang)
}
