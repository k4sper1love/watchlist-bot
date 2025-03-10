package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

var settingsButtons = []Button{
	{"🈳", "language", states.CallbackSettingsLanguage, "", true},
	{"🌐", "kinopoiskToken", states.CallbackSettingsKinopoiskToken, "", true},
	{"🔢", "collectionsPageSize", states.CallbackSettingsCollectionsPageSize, "", true},
	{"🔢", "filmsPageSize", states.CallbackSettingsFilmsPageSize, "", true},
	{"🔢", "objectsPageSize", states.CallbackSettingsObjectsPageSize, "", true},
}

func Settings(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddButtons(settingsButtons...).
		AddBack("").
		Build(session.Lang)
}

func SettingsLanguageSelect(session *models.Session, languages []string) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddLanguageSelect(languages, states.PrefixSelectLang).
		AddBack(states.CallbackSettingsBack).
		Build(session.Lang)
}
