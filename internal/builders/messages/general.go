package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"strings"
)

func BuildStartMessage(app models.App, session *models.Session) string {
	part1 := translator.Translate(session.Lang, "welcomeMessageGreeting", map[string]interface{}{
		"Name": utils.ParseTelegramName(app.Update),
	}, nil)

	part2 := translator.Translate(session.Lang, "welcomeMessageBody", map[string]interface{}{
		"Version": app.Config.Version,
	}, nil)

	part3 := translator.Translate(session.Lang, "welcomeMessageCallToAction", nil, nil)

	msg := fmt.Sprintf("👋 <b>%s</b>\n\n%s 🚀\n\n%s", part1, part2, part3)

	return msg
}

func BuildHelpMessage(session *models.Session) string {
	msg := translator.Translate(session.Lang, "helpMessage", nil, nil)

	return msg
}

func BuildMenuMessage(session *models.Session) string {
	part1 := translator.Translate(session.Lang, "mainMenu", nil, nil)
	part2 := translator.Translate(session.Lang, "choiceAction", nil, nil)

	msg := fmt.Sprintf("📋 <b>%s</b>\n\n%s", part1, part2)

	return msg
}

func BuildLanguagesMessage(languages []string) string {
	var res string
	for _, language := range languages {
		translatedMsg := translator.Translate(language, "choiceLanguage", nil, nil)
		upperLanguage := strings.ToUpper(language)
		res += fmt.Sprintf("<b>%s</b>: %s\n\n", upperLanguage, translatedMsg)
	}

	return res
}

func BuildLanguagesFailureMessage(session *models.Session) string {
	part1 := translator.Translate(session.Lang, "someError", nil, nil)
	part2 := translator.Translate(session.Lang, "setDefaultLanguage", nil, nil)

	return fmt.Sprintf("🚨 %s\n\n%s", part1, part2)
}

func BuildKinopoiskTokenMessage(session *models.Session) string {
	msg := "⚠️ " + translator.Translate(session.Lang, "tokenRequestInfo", nil, nil)

	if session.KinopoiskAPIToken != "" {
		part := translator.Translate(session.Lang, "currentToken", nil, nil)
		msg += fmt.Sprintf("\n\n%s: <code>%s</code>", part, session.KinopoiskAPIToken)
	}

	part := translator.Translate(session.Lang, "tokenRequest", nil, nil)
	msg += fmt.Sprintf("\n\n<b>%s</b>", part)

	return msg
}

func BuildKinopoiskTokenSuccessMessage(session *models.Session) string {
	return "✅ " + translator.Translate(session.Lang, "tokenSuccess", nil, nil)
}

func BuildCancelActionMessage(session *models.Session) string {
	return "🚫 " + translator.Translate(session.Lang, "cancelAction", nil, nil)
}

func BuildLastPageAlertMessage(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
}

func BuildFirstPageAlertMessage(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
}

func BuildImageFailureMessage(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "getImageFailure", nil, nil)
}

func BuildChoiceWayMessage(session *models.Session) string {
	return fmt.Sprintf("<b>%s</b>", translator.Translate(session.Lang, "choiceWay", nil, nil))
}

func BuildTokenCodeMessage(session *models.Session, code int) string {
	return "🚨 " + translator.Translate(session.Lang, fmt.Sprintf("token%d", code), nil, nil)
}

func toBold(text string) string {
	return fmt.Sprintf("<b>%s</b>", text)
}

func toCode(text string) string {
	return fmt.Sprintf("<code>%s</code>", text)
}
