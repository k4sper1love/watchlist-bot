package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"strings"
)

func BuildSettingsMessage(session *models.Session) string {
	part1 := translator.Translate(session.Lang, "settings", nil, nil)
	part2 := translator.Translate(session.Lang, "settingsChoice", nil, nil)

	return fmt.Sprintf("⚙️ <b>%s</b>\n\n%s", part1, part2)
}

func BuildSettingsLanguageMessage(session *models.Session) string {
	part1 := translator.Translate(session.Lang, "currentLanguage", nil, nil)
	part2 := strings.ToUpper(session.Lang)
	part3 := translator.Translate(session.Lang, "languageChoice", nil, nil)

	return fmt.Sprintf("🈳 <b>%s:</b> <code>%s</code>\n\n%s", part1, part2, part3)
}

func BuildSettingsLanguageSuccessMessage(session *models.Session) string {
	return "🔄 " + translator.Translate(session.Lang, "settingsLanguageSuccess", map[string]interface{}{
		"Language": strings.ToUpper(session.Lang),
	}, nil)
}

func BuildSettingsPageSizeMessage(session *models.Session, pageSize int) string {
	part1 := translator.Translate(session.Lang, "currentPageSize", nil, nil)
	part2 := translator.Translate(session.Lang, "settingsPageSizeChoice", nil, nil)

	return fmt.Sprintf("🔢 <b>%s</b>: <code>%d</code>\n\n%s", part1, pageSize, part2)
}

func BuildSettingsPageSizeFailureMessage(session *models.Session, min, max int) string {
	return "❗️" + translator.Translate(session.Lang, "invalidInputRange", map[string]interface{}{
		"Min": min,
		"Max": max,
	}, nil)
}

func BuildSettingsPageSizeSuccessMessage(session *models.Session, pageSize int) string {
	return "🔄 " + translator.Translate(session.Lang, "settingsPageSizeSuccess", map[string]interface{}{
		"Size": pageSize,
	}, nil)
}
