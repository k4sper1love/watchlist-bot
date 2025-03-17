package general

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log/slog"
	"strings"
)

// HandleStartCommand handles the "/start" command.
// Sends a welcome message and prompts the user to select a language.
func HandleStartCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.Start(app, session), nil)
	HandleLanguageCommand(app, session)
}

// HandleLanguageCommand handles the command for selecting a language.
// Retrieves supported languages from the locales directory and sends a message with options to choose one.
func HandleLanguageCommand(app models.App, session *models.Session) {
	if languages, err := utils.ParseSupportedLanguages(app.Config.LocalesDir); err != nil {
		sl.Log.Error("failed to parse supported languages", slog.Any("error", err), slog.String("dir", app.Config.LocalesDir))
		app.SendMessage(messages.LanguagesFailure(session), keyboards.Back(session, ""))
	} else {
		app.SendMessage(messages.Languages(languages), keyboards.LanguageSelect(languages))
	}
}

// HandleMenuCommand handles the command for displaying the main menu.
// Sends a message with options for navigating to different sections of the bot.
func HandleMenuCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.Menu(session), keyboards.Menu(session))
}

// HandleHelpCommand handles the "/help" command.
// Sends a message with instructions or information about using the bot.
func HandleHelpCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.Help(session), nil)
}

// HandleLanguageButton handles button interactions related to language selection.
// Updates the session's language and navigates to the main menu.
func HandleLanguageButton(app models.App, session *models.Session) {
	session.Lang = strings.TrimPrefix(utils.ParseCallback(app.Update), states.SelectStartLang)
	HandleMenuCommand(app, session)
}
