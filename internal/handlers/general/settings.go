package general

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"strconv"
	"strings"
)

func HandleSettingsCommand(app models.App, session *models.Session) {
	part1 := translator.Translate(session.Lang, "settings", nil, nil)
	part2 := translator.Translate(session.Lang, "settingsChoice", nil, nil)

	msg := fmt.Sprintf("⚙️ <b>%s</b>\n\n%s", part1, part2)

	keyboard := keyboards.BuildSettingsKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func HandleSettingsButton(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch {
	case callback == states.CallbackSettingsBack:
		HandleSettingsCommand(app, session)

	case callback == states.CallbackSettingsLanguage:
		handleLanguage(app, session)

	case callback == states.CallbackSettingsKinopoiskToken:
		handleKinopoiskToken(app, session)

	case callback == states.CallbackSettingsCollectionsPageSize:
		handleCollectionsPageSize(app, session)

	case callback == states.CallbackSettingsFilmsPageSize:
		handleFilmsPageSize(app, session)

	case callback == states.CallbackSettingsObjectsPageSize:
		handleObjectsPageSize(app, session)

	case strings.HasPrefix(callback, "select_lang_"):
		handleLanguageSelect(app, session)
	}

}

func HandleSettingsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearState()
		HandleSettingsCommand(app, session)
	}

	switch session.State {
	case states.ProcessSettingsAwaitingKinopoiskToken:
		parseKinopoiskToken(app, session)
	case states.ProcessSettingsCollectionsAwaitingPageSize:
		parseCollectionsPageSize(app, session)
	case states.ProcessSettingsFilmsAwaitingPageSize:
		parseFilmsPageSize(app, session)
	case states.ProcessSettingsObjectsAwaitingPageSize:
		parseObjectsPageSize(app, session)
	}
}

func handleLanguage(app models.App, session *models.Session) {
	localesDir := "./locales"

	languages, err := utils.ParseSupportedLanguages(localesDir)
	if err != nil {
		msg := "❗️" + translator.Translate(session.Lang, "parseLanguageFailure", nil, nil)
		app.SendMessage(msg, nil)
		session.ClearAllStates()
		HandleMenuCommand(app, session)
		return
	}

	part1 := translator.Translate(session.Lang, "currentLanguage", nil, nil)
	part2 := strings.ToUpper(session.Lang)
	part3 := translator.Translate(session.Lang, "languageChoice", nil, nil)

	msg := fmt.Sprintf("🈳 <b>%s:</b> <code>%s</code>\n\n%s", part1, part2, part3)

	keyboard := keyboards.NewKeyboard().AddLanguageSelect(languages, "select_lang").AddBack(states.CallbackSettingsBack).Build(session.Lang)

	app.SendMessage(msg, keyboard)
}

func handleLanguageSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	lang := strings.TrimPrefix(callback, "select_lang_")

	session.Lang = lang

	msg := "🔄 " + translator.Translate(session.Lang, "settingsLanguageSuccess", map[string]interface{}{
		"Language": strings.ToUpper(lang),
	}, nil)

	app.SendMessage(msg, nil)

	HandleSettingsCommand(app, session)
}

func handleKinopoiskToken(app models.App, session *models.Session) {
	msg := messages.BuildKinopoiskTokenMessage(session)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessSettingsAwaitingKinopoiskToken)
}

func parseKinopoiskToken(app models.App, session *models.Session) {
	token := utils.ParseMessageString(app.Update)

	session.KinopoiskAPIToken = token

	msg := messages.BuildKinopoiskTokenSuccessMessage(session)
	app.SendMessage(msg, nil)

	session.ClearState()

	HandleSettingsCommand(app, session)
}

func handleCollectionsPageSize(app models.App, session *models.Session) {
	part1 := translator.Translate(session.Lang, "currentPageSize", nil, nil)
	part2 := translator.Translate(session.Lang, "settingsPageSizeChoice", nil, nil)

	msg := fmt.Sprintf("🔢 <b>%s</b>: <code>%d</code>\n\n%s", part1, session.CollectionsState.PageSize, part2)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessSettingsCollectionsAwaitingPageSize)
}

func parseCollectionsPageSize(app models.App, session *models.Session) {
	pageSize, err := strconv.Atoi(utils.ParseMessageString(app.Update))
	if err != nil || pageSize < 1 || pageSize > 100 {
		handleCollectionsPageSize(app, session)
		return
	}

	session.CollectionsState.PageSize = pageSize

	msg := "🔄 " + translator.Translate(session.Lang, "settingsPageSizeSuccess", map[string]interface{}{
		"Size": pageSize,
	}, nil)

	app.SendMessage(msg, nil)

	session.ClearState()

	HandleSettingsCommand(app, session)
}

func handleFilmsPageSize(app models.App, session *models.Session) {
	part1 := translator.Translate(session.Lang, "currentPageSize", nil, nil)
	part2 := translator.Translate(session.Lang, "settingsPageSizeChoice", nil, nil)

	msg := fmt.Sprintf("🔢 <b>%s</b>: <code>%d</code>\n\n%s", part1, session.CollectionsState.PageSize, part2)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessSettingsFilmsAwaitingPageSize)
}

func parseFilmsPageSize(app models.App, session *models.Session) {
	pageSize, err := strconv.Atoi(utils.ParseMessageString(app.Update))
	if err != nil || pageSize < 1 {
		handleCollectionsPageSize(app, session)
		return
	}

	session.FilmsState.PageSize = pageSize

	msg := "🔄 " + translator.Translate(session.Lang, "settingsPageSizeSuccess", map[string]interface{}{
		"Size": pageSize,
	}, nil)

	app.SendMessage(msg, nil)

	session.ClearState()

	HandleSettingsCommand(app, session)
}

func handleObjectsPageSize(app models.App, session *models.Session) {
	part1 := translator.Translate(session.Lang, "currentPageSize", nil, nil)
	part2 := translator.Translate(session.Lang, "settingsPageSizeChoice", nil, nil)

	msg := fmt.Sprintf("🔢 <b>%s</b>: <code>%d</code>\n\n%s", part1, session.CollectionsState.PageSize, part2)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessSettingsObjectsAwaitingPageSize)
}

func parseObjectsPageSize(app models.App, session *models.Session) {
	pageSize, err := strconv.Atoi(utils.ParseMessageString(app.Update))
	if err != nil || pageSize < 1 {
		handleCollectionsPageSize(app, session)
		return
	}

	session.CollectionFilmsState.PageSize = pageSize

	msg := "🔄 " + translator.Translate(session.Lang, "settingsPageSizeSuccess", map[string]interface{}{
		"Size": pageSize,
	}, nil)

	app.SendMessage(msg, nil)

	session.ClearState()

	HandleSettingsCommand(app, session)
}
