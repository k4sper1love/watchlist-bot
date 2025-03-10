package general

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log/slog"
	"strconv"
	"strings"
)

func HandleSettingsCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.Settings(session), keyboards.Settings(session))
}

func HandleSettingsButton(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallbackSettingsBack:
		HandleSettingsCommand(app, session)
	case states.CallbackSettingsLanguage:
		handleLanguage(app, session)
	case states.CallbackSettingsKinopoiskToken:
		handleKinopoiskToken(app, session)
	case states.CallbackSettingsCollectionsPageSize,
		states.CallbackSettingsFilmsPageSize,
		states.CallbackSettingsObjectsPageSize:
		handlePageSizeSetting(app, session)
	default:
		if strings.HasPrefix(utils.ParseCallback(app.Update), states.PrefixSelectLang) {
			handleLanguageSelect(app, session)
		}
	}
}

func HandleSettingsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		returnToSettingsMenu(app, session)
		return
	}

	switch session.State {
	case states.ProcessSettingsAwaitingKinopoiskToken:
		parseKinopoiskToken(app, session)
	case states.ProcessSettingsCollectionsAwaitingPageSize,
		states.ProcessSettingsFilmsAwaitingPageSize,
		states.ProcessSettingsObjectsAwaitingPageSize:
		parsePageSizeSetting(app, session)
	}
}

func handleLanguage(app models.App, session *models.Session) {
	if languages, err := utils.ParseSupportedLanguages(app.Config.LocalesDir); err != nil {
		sl.Log.Error("failed to parse supported languages", slog.Any("error", err), slog.String("dir", app.Config.LocalesDir))
		app.SendMessage(messages.LanguagesFailure(session), keyboards.Back(session, states.CallbackMenuSelectSettings))
	} else {
		app.SendMessage(messages.SettingsLanguage(session), keyboards.SettingsLanguageSelect(session, languages))
	}
}

func handleLanguageSelect(app models.App, session *models.Session) {
	session.Lang = strings.TrimPrefix(utils.ParseCallback(app.Update), states.PrefixSelectLang)
	app.SendMessage(messages.SettingsLanguageSuccess(session), nil)
	returnToSettingsMenu(app, session)
}

func handleKinopoiskToken(app models.App, session *models.Session) {
	app.SendMessage(messages.KinopoiskToken(session), keyboards.Cancel(session))
	session.SetState(states.ProcessSettingsAwaitingKinopoiskToken)
}

func parseKinopoiskToken(app models.App, session *models.Session) {
	session.KinopoiskAPIToken = utils.ParseMessageString(app.Update)
	app.SendMessage(messages.KinopoiskTokenSuccess(session), nil)
	returnToSettingsMenu(app, session)
}

func handlePageSizeSetting(app models.App, session *models.Session) {
	var pageSize int
	switch utils.ParseCallback(app.Update) {
	case states.CallbackSettingsFilmsPageSize:
		pageSize = session.FilmsState.PageSize
		session.SetState(states.ProcessSettingsFilmsAwaitingPageSize)
	case states.CallbackSettingsCollectionsPageSize:
		pageSize = session.CollectionsState.PageSize
		session.SetState(states.ProcessSettingsCollectionsAwaitingPageSize)
	case states.CallbackSettingsObjectsPageSize:
		pageSize = session.CollectionFilmsState.PageSize
		session.SetState(states.ProcessSettingsObjectsAwaitingPageSize)
	}

	app.SendMessage(messages.SettingsPageSize(session, pageSize), keyboards.Cancel(session))
}

func parsePageSizeSetting(app models.App, session *models.Session) {
	pageSize, ok := parseAndValidatePageSize(app, session)
	if !ok {
		handlePageSizeSetting(app, session)
		return
	}

	switch session.State {
	case states.ProcessSettingsFilmsAwaitingPageSize:
		session.FilmsState.PageSize = pageSize
	case states.ProcessSettingsCollectionsAwaitingPageSize:
		session.CollectionsState.PageSize = pageSize
	case states.ProcessSettingsObjectsAwaitingPageSize:
		session.CollectionFilmsState.PageSize = pageSize
	}

	handlePageSizeSuccess(app, session, pageSize)
}

func parseAndValidatePageSize(app models.App, session *models.Session) (int, bool) {
	pageSize, _ := strconv.Atoi(utils.ParseMessageString(app.Update))
	if pageSize < 1 || pageSize > 10 {
		validator.HandleInvalidInputRange(app, session, 1, 10)
		return 0, false
	}
	return pageSize, true
}

func handlePageSizeSuccess(app models.App, session *models.Session, pageSize int) {
	app.SendMessage(messages.SettingsPageSizeSuccess(session, pageSize), nil)
	returnToSettingsMenu(app, session)
}

func returnToSettingsMenu(app models.App, session *models.Session) {
	session.ClearState()
	HandleSettingsCommand(app, session)
}
