package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func BuildSettingsLanguageSelectKeyboard(session *models.Session, languages []string) *tgbotapi.InlineKeyboardMarkup {
	return NewKeyboard().AddLanguageSelect(languages, states.PrefixSelectLang).AddBack(states.CallbackSettingsBack).Build(session.Lang)
}
