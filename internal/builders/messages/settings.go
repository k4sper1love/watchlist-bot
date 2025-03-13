package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"strconv"
	"strings"
)

func Settings(session *models.Session) string {
	return fmt.Sprintf("⚙️ %s\n\n%s",
		toBold(translator.Translate(session.Lang, "settings", nil, nil)),
		translator.Translate(session.Lang, "settingsChoice", nil, nil))
}

func SettingsLanguage(session *models.Session) string {
	return fmt.Sprintf("🈳 %s: %s\n\n%s",
		toBold(translator.Translate(session.Lang, "currentLanguage", nil, nil)),
		toCode(strings.ToUpper(session.Lang)),
		translator.Translate(session.Lang, "languageChoice", nil, nil))
}

func SettingsLanguageSuccess(session *models.Session) string {
	return "🔄 " + translator.Translate(session.Lang, "settingsLanguageSuccess", map[string]interface{}{
		"Language": strings.ToUpper(session.Lang),
	}, nil)
}

func SettingsPageSize(session *models.Session, pageSize int) string {
	return fmt.Sprintf("🔢 %s: %s\n\n%s",
		toBold(translator.Translate(session.Lang, "currentPageSize", nil, nil)),
		toCode(strconv.Itoa(pageSize)),
		translator.Translate(session.Lang, "settingsPageSizeChoice", nil, nil))
}

func SettingsPageSizeSuccess(session *models.Session) string {
	return "🔄 " + translator.Translate(session.Lang, "settingsPageSizeSuccess", nil, nil)
}
