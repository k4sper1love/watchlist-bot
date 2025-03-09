package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"strings"
)

func Start(app models.App, session *models.Session) string {
	return fmt.Sprintf("👋%s\n\n%s🚀\n\n%s",
		toBold(translator.Translate(session.Lang, "welcomeMessageGreeting", map[string]interface{}{
			"Name": utils.ParseTelegramName(app.Update),
		}, nil)),
		translator.Translate(session.Lang, "welcomeMessageBody", map[string]interface{}{
			"Version": app.Config.Version,
		}, nil),
		translator.Translate(session.Lang, "welcomeMessageCallToAction", nil, nil))
}

func Help(session *models.Session) string {
	return translator.Translate(session.Lang, "helpMessage", nil, nil)
}

func Menu(session *models.Session) string {
	return fmt.Sprintf("📋 %s\n\n%s",
		toBold(translator.Translate(session.Lang, "mainMenu", nil, nil)),
		translator.Translate(session.Lang, "choiceAction", nil, nil))
}

func Languages(languages []string) string {
	var msg strings.Builder
	for _, language := range languages {
		msg.WriteString(fmt.Sprintf("%s: %s\n\n",
			toBold(strings.ToUpper(language)),
			translator.Translate(language, "choiceLanguage", nil, nil)))
	}
	return msg.String()
}

func LanguagesFailure(session *models.Session) string {
	return fmt.Sprintf("🚨 %s\n\n%s",
		translator.Translate(session.Lang, "someError", nil, nil),
		translator.Translate(session.Lang, "setDefaultLanguage", nil, nil))
}

func KinopoiskToken(session *models.Session) string {
	return fmt.Sprintf("⚠️ %s\n\n%s%s",
		translator.Translate(session.Lang, "tokenRequestInfo", nil, nil),
		formatOptionalString(translator.Translate(session.Lang, "currentToken", nil, nil),
			toCode(session.KinopoiskAPIToken), "%s: %s\n\n"),
		toBold(translator.Translate(session.Lang, "tokenRequest", nil, nil)))
}

func KinopoiskTokenSuccess(session *models.Session) string {
	return "✅ " + translator.Translate(session.Lang, "tokenSuccess", nil, nil)
}

func CancelAction(session *models.Session) string {
	return "🚫 " + translator.Translate(session.Lang, "cancelAction", nil, nil)
}

func LastPageAlert(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
}

func FirstPageAlert(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
}

func ImageFailure(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "getImageFailure", nil, nil)
}

func ChoiceWay(session *models.Session) string {
	return toBold(translator.Translate(session.Lang, "choiceWay", nil, nil))
}

func KinopoiskFailureCode(session *models.Session, code int) string {
	return "🚨 " + translator.Translate(session.Lang, fmt.Sprintf("token%d", code), nil, nil)
}

func SomeError(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
}

func NotFound(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "notFound", nil, nil)
}

func RequestFailure(session *models.Session) string {
	return "🚨" + translator.Translate(session.Lang, "requestFailure", nil, nil)
}
