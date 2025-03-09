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

func HandleStartCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.Start(app, session), nil)
	HandleLanguageCommand(app, session)
}

func HandleHelpCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.Help(session), nil)
}

func HandleLanguageCommand(app models.App, session *models.Session) {
	languages, err := utils.ParseSupportedLanguages(app.Config.LocalesDir)
	if err != nil {
		sl.Log.Error("failed to parse supported languages", slog.Any("error", err), slog.String("dir", app.Config.LocalesDir))
		app.SendMessage(messages.LanguagesFailure(session), keyboards.BuildKeyboardWithBack(session, ""))
		return
	}

	app.SendMessage(messages.Languages(languages), keyboards.BuildLanguageSelectKeyboard(languages))
}

func HandleLanguageButton(app models.App, session *models.Session) {
	session.Lang = strings.TrimPrefix(utils.ParseCallback(app.Update), states.PrefixSelectStartLang)
	HandleMenuCommand(app, session)
}

func HandleMenuCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.Menu(session), keyboards.BuildMenuKeyboard(session))
}
