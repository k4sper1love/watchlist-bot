package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/security"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"strconv"
	"strings"
)

// Start generates a welcome message for the user when they start interacting with the bot.
// Includes a greeting and information about the bot version.
func Start(app models.App, session *models.Session) string {
	return fmt.Sprintf("👋 %s\n\n%s 🚀",
		toBold(translator.Translate(session.Lang, "welcomeMessageGreeting", map[string]interface{}{
			"Name": utils.ParseTelegramName(app.Update),
		}, nil)),
		translator.Translate(session.Lang, "welcomeMessageBody", map[string]interface{}{
			"Version": app.Config.Version,
		}, nil))
}

// Help generates a help message for the user.
func Help(session *models.Session) string {
	return translator.Translate(session.Lang, "helpMessage", nil, nil)
}

// Menu generates a message for the main menu.
func Menu(session *models.Session) string {
	return fmt.Sprintf("📋 %s\n\n%s",
		toBold(translator.Translate(session.Lang, "mainMenu", nil, nil)),
		translator.Translate(session.Lang, "choiceAction", nil, nil))
}

// Languages generates a message listing available languages for selection.
// Each language is displayed with its name and a localized description.
func Languages(languages []string) string {
	var msg strings.Builder
	for _, language := range languages {
		msg.WriteString(fmt.Sprintf("%s: %s\n\n",
			toBold(strings.ToUpper(language)),
			translator.Translate(language, "choiceLanguage", nil, nil)))
	}
	return msg.String()
}

// LanguagesFailure generates a message indicating an error occurred while setting the language.
// Includes a fallback instruction to set the default language.
func LanguagesFailure(session *models.Session) string {
	return fmt.Sprintf("🚨 %s\n\n%s",
		translator.Translate(session.Lang, "someError", nil, nil),
		translator.Translate(session.Lang, "setDefaultLanguage", nil, nil))
}

// RequestKinopoiskToken generates a message prompting the user to provide a Kinopoisk API token.
// Optionally includes the current token if available.
func RequestKinopoiskToken(session *models.Session) string {
	token, _ := security.Decrypt(session.KinopoiskAPIToken)
	return fmt.Sprintf("⚠️ %s\n\n%s%s",
		translator.Translate(session.Lang, "tokenRequestInfo", nil, nil),
		formatOptionalString(translator.Translate(session.Lang, "currentToken", nil, nil),
			toCode(token), "%s: %s\n\n"),
		toBold(translator.Translate(session.Lang, "tokenRequest", nil, nil)))
}

// KinopoiskTokenSuccess generates a success message after the user sets a valid Kinopoisk API token.
func KinopoiskTokenSuccess(session *models.Session) string {
	return "✅ " + translator.Translate(session.Lang, "tokenSuccess", nil, nil)
}

// UnknownCommand generates a message for unrecognized commands.
func UnknownCommand(session *models.Session) string {
	return "❗" + translator.Translate(session.Lang, "unknownCommand", nil, nil)
}

// UnknownState generates a message for unrecognized bot states.
func UnknownState(session *models.Session) string {
	return "❗" + translator.Translate(session.Lang, "unknownState", nil, nil)
}

// SessionError generates a message for session-related errors.
func SessionError(lang string) string {
	return translator.Translate(lang, "sessionError", nil, nil)
}

// CancelAction generates a message confirming that an action has been canceled.
func CancelAction(session *models.Session) string {
	return "🚫 " + translator.Translate(session.Lang, "cancelAction", nil, nil)
}

// LastPageAlert generates a message notifying the user that they are on the last page.
func LastPageAlert(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
}

// FirstPageAlert generates a message notifying the user that they are on the first page.
func FirstPageAlert(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
}

// ImageFailure generates a message indicating an error occurred while processing an image.
func ImageFailure(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "getImageFailure", nil, nil)
}

// ChoiceWay generates a message prompting the user to choose an action or path.
func ChoiceWay(session *models.Session) string {
	return toBold(translator.Translate(session.Lang, "choiceWay", nil, nil))
}

// KinopoiskFailureCode generates a message indicating a failure related to the Kinopoisk API.
// Includes the specific error code.
func KinopoiskFailureCode(session *models.Session, code int) string {
	return "🚨 " + translator.Translate(session.Lang, "tokenCodeError."+strconv.Itoa(code), nil, nil)
}

// SomeError generates a generic error message.
func SomeError(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "someError", nil, nil)
}

// NotFound generates a message indicating that the requested resource was not found.
func NotFound(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "notFound", nil, nil)
}

// RequestFailure generates a message indicating that a request failed.
func RequestFailure(session *models.Session) string {
	return "🚨" + translator.Translate(session.Lang, "requestFailure", nil, nil)
}
