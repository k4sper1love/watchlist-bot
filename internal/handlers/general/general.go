package general

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"log/slog"
	"strings"
)

func HandleStartCommand(app models.App, session *models.Session) {
	msg := messages.BuildStartMessage(app, session)
	app.SendMessage(msg, nil)

	if err := HandleAuthProcess(app, session); err != nil {
		msg = translator.Translate(session.Lang, "authFailure", nil, nil)
		app.SendMessage(msg, nil)
		sl.Log.Error("error auth process", slog.Any("err", err))
		return
	}

	HandleLanguageCommand(app, session)
}

func HandleHelpCommand(app models.App, session *models.Session) {
	msg := messages.BuildHelpMessage(session)

	app.SendMessage(msg, nil)
}

func HandleLanguageCommand(app models.App, session *models.Session) {
	msg, err := messages.BuildLanguageMessage()
	if err != nil {
		handleLanguageError(app, session)
		return
	}

	languages, err := utils.ParseSupportedLanguages("./locales")
	if err != nil {
		handleLanguageError(app, session)
		return
	}

	keyboard := keyboards.NewKeyboard().AddLanguageSelect(languages, "select_start_lang").Build("")

	app.SendMessage(msg, keyboard)
}

func HandleLanguageButton(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)

	lang := strings.TrimPrefix(callback, "select_start_lang_")

	session.Lang = lang

	HandleMenuCommand(app, session)
}

func handleLanguageError(app models.App, session *models.Session) {
	part1 := translator.Translate(session.Lang, "someError", nil, nil)
	part2 := translator.Translate(session.Lang, "setDefaultLanguage", nil, nil)

	msg := fmt.Sprintf("%s\n\n%s", part1, part2)
	keyboard := keyboards.NewKeyboard().AddBack("").Build(session.Lang)
	app.SendMessage(msg, keyboard)

}
