package general

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log/slog"
	"strconv"
	"strings"
)

func HandleSettingsCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildSettingsMessage(session), keyboards.BuildSettingsKeyboard(session))
}

func HandleSettingsButton(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallbackSettingsBack:
		HandleSettingsCommand(app, session)
	case states.CallbackSettingsLanguage:
		handleLanguage(app, session)
	case states.CallbackSettingsKinopoiskToken:
		handleKinopoiskToken(app, session)
	case states.CallbackSettingsCollectionsPageSize:
		handleCollectionsPageSize(app, session)
	case states.CallbackSettingsFilmsPageSize:
		handleFilmsPageSize(app, session)
	case states.CallbackSettingsObjectsPageSize:
		handleObjectsPageSize(app, session)
	default:
		if strings.HasPrefix(utils.ParseCallback(app.Update), states.PrefixSelectLang) {
			handleLanguageSelect(app, session)
		}
	}
}

func HandleSettingsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearState()
		HandleSettingsCommand(app, session)
		return
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
	languages, err := utils.ParseSupportedLanguages(app.Config.LocalesDir)
	if err != nil {
		sl.Log.Error("failed to parse supported languages", slog.Any("error", err), slog.String("dir", app.Config.LocalesDir))
		app.SendMessage(messages.BuildLanguagesFailureMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackMenuSelectSettings))
		return
	}
	app.SendMessage(messages.BuildSettingsLanguageMessage(session), keyboards.BuildSettingsLanguageSelectKeyboard(session, languages))
}

func handleLanguageSelect(app models.App, session *models.Session) {
	session.Lang = strings.TrimPrefix(utils.ParseCallback(app.Update), "select_lang_")
	app.SendMessage(messages.BuildSettingsLanguageSuccessMessage(session), nil)
	HandleSettingsCommand(app, session)
}

func handleKinopoiskToken(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildKinopoiskTokenMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessSettingsAwaitingKinopoiskToken)
}

func parseKinopoiskToken(app models.App, session *models.Session) {
	session.KinopoiskAPIToken = utils.ParseMessageString(app.Update)
	app.SendMessage(messages.BuildKinopoiskTokenSuccessMessage(session), nil)
	session.ClearState()
	HandleSettingsCommand(app, session)
}

func handleCollectionsPageSize(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildSettingsPageSizeMessage(session, session.CollectionsState.PageSize), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessSettingsCollectionsAwaitingPageSize)
}

func parseCollectionsPageSize(app models.App, session *models.Session) {
	pageSize, ok := parseAndValidatePageSize(app, session)
	if !ok {
		handleCollectionsPageSize(app, session)
		return
	}
	session.CollectionsState.PageSize = pageSize
	handlePageSizeSuccess(app, session, pageSize)
}

func handleFilmsPageSize(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildSettingsPageSizeMessage(session, session.FilmsState.PageSize), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessSettingsFilmsAwaitingPageSize)
}

func parseFilmsPageSize(app models.App, session *models.Session) {
	pageSize, ok := parseAndValidatePageSize(app, session)
	if !ok {
		handleFilmsPageSize(app, session)
		return
	}
	session.FilmsState.PageSize = pageSize
	handlePageSizeSuccess(app, session, pageSize)
}

func handleObjectsPageSize(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildSettingsPageSizeMessage(session, session.CollectionFilmsState.PageSize), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessSettingsObjectsAwaitingPageSize)
}

func parseObjectsPageSize(app models.App, session *models.Session) {
	pageSize, ok := parseAndValidatePageSize(app, session)
	if !ok {
		handleObjectsPageSize(app, session)
		return
	}
	session.CollectionFilmsState.PageSize = pageSize
	handlePageSizeSuccess(app, session, pageSize)
}

func parseAndValidatePageSize(app models.App, session *models.Session) (int, bool) {
	pageSize, _ := strconv.Atoi(utils.ParseMessageString(app.Update))
	if pageSize < 1 || pageSize > 10 {
		app.SendMessage(messages.BuildSettingsPageSizeFailureMessage(session, 1, 10), nil)
		return 0, false
	}
	return pageSize, true
}

func handlePageSizeSuccess(app models.App, session *models.Session, pageSize int) {
	app.SendMessage(messages.BuildSettingsPageSizeSuccessMessage(session, pageSize), nil)
	session.ClearState()
	HandleSettingsCommand(app, session)
}
