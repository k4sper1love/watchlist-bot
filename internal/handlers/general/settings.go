package general

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
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

	msg := fmt.Sprintf("⚙️ <b>%s:</b>\n\n%s", part1, part2)

	keyboard := keyboards.BuildSettingsKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func HandleSettingsButton(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)

	switch {
	case callback == states.CallbackSettingsBack:
		HandleSettingsCommand(app, session)

	case callback == states.CallbackSettingsLanguage:
		handleLanguage(app, session)

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
	if utils.IsCancel(app.Upd) {
		session.ClearState()
		HandleSettingsCommand(app, session)
	}

	switch session.State {
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
		msg := translator.Translate(session.Lang, "parseLanguageFailure", nil, nil)
		app.SendMessage(msg, nil)
		session.ClearAllStates()
		HandleMenuCommand(app, session)
		return
	}

	part1 := translator.Translate(session.Lang, "currentLanguage", nil, nil)
	part2 := translator.Translate(session.Lang, "settingsLanguageChoice", nil, nil)

	msg := fmt.Sprintf("%s: %s\n%s", part1, session.Lang, part2)

	keyboard := keyboards.NewKeyboard().AddLanguageSelect(languages).AddBack(states.CallbackSettingsBack).Build("")

	app.SendMessage(msg, keyboard)
}

func handleLanguageSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)

	lang := strings.TrimPrefix(callback, "select_lang_")

	session.Lang = lang

	msg := translator.Translate(session.Lang, "settingsLanguageSuccess", map[string]interface{}{
		"Language": lang,
	}, nil)

	app.SendMessage(msg, nil)

	HandleSettingsCommand(app, session)
}

func handleCollectionsPageSize(app models.App, session *models.Session) {
	part1 := translator.Translate(session.Lang, "currentCollectionsPageSize", nil, nil)
	part2 := translator.Translate(session.Lang, "settingsPageSizeChoice", nil, nil)

	msg := fmt.Sprintf("%s:%d\n%s", part1, session.CollectionsState.PageSize, part2)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessSettingsCollectionsAwaitingPageSize)
}

func parseCollectionsPageSize(app models.App, session *models.Session) {
	pageSize, err := strconv.Atoi(utils.ParseMessageString(app.Upd))
	if err != nil || pageSize < 1 {
		msg := translator.Translate(session.Lang, "settingsPageSizeChoice", nil, nil)
		app.SendMessage(msg, nil)
		return
	}

	session.CollectionsState.PageSize = pageSize

	successMsg := translator.Translate(session.Lang, "settingsCollectionPageSizeSuccess", nil, nil)
	msg := fmt.Sprintf("%s:%d", successMsg, pageSize)

	app.SendMessage(msg, nil)

	session.ClearState()

	HandleSettingsCommand(app, session)
}

func handleFilmsPageSize(app models.App, session *models.Session) {
	part1 := translator.Translate(session.Lang, "currentFilmsPageSize", nil, nil)
	part2 := translator.Translate(session.Lang, "settingsPageSizeChoice", nil, nil)

	msg := fmt.Sprintf("%s: %d\n%s", part1, session.FilmsState.PageSize, part2)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessSettingsFilmsAwaitingPageSize)
}

func parseFilmsPageSize(app models.App, session *models.Session) {
	pageSize, err := strconv.Atoi(utils.ParseMessageString(app.Upd))
	if err != nil || pageSize < 1 {
		msg := translator.Translate(session.Lang, "settingsPageSizeChoice", nil, nil)
		app.SendMessage(msg, nil)
		return
	}

	session.FilmsState.PageSize = pageSize

	successMsg := translator.Translate(session.Lang, "settingsFilmsPageSizeSuccess", nil, nil)
	msg := fmt.Sprintf("%s: %d", successMsg, pageSize)

	app.SendMessage(msg, nil)

	session.ClearState()

	HandleSettingsCommand(app, session)
}

func handleObjectsPageSize(app models.App, session *models.Session) {
	part1 := translator.Translate(session.Lang, "currentObjectsPageSize", nil, nil)
	part2 := translator.Translate(session.Lang, "settingsPageSizeChoice", nil, nil)

	msg := fmt.Sprintf("%s: %d\n%s", part1, session.CollectionFilmsState.PageSize, part2)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessSettingsObjectsAwaitingPageSize)
}

func parseObjectsPageSize(app models.App, session *models.Session) {
	pageSize, err := strconv.Atoi(utils.ParseMessageString(app.Upd))
	if err != nil || pageSize < 1 {
		msg := translator.Translate(session.Lang, "settingsPageSizeChoice", nil, nil)
		app.SendMessage(msg, nil)
		return
	}

	session.CollectionFilmsState.PageSize = pageSize

	successMsg := translator.Translate(session.Lang, "settingsObjectsPageSizeSuccess", nil, nil)
	msg := fmt.Sprintf("%s: %d", successMsg, pageSize)

	app.SendMessage(msg, nil)

	session.ClearState()

	HandleSettingsCommand(app, session)
}
