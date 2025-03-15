package general

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log/slog"
	"strings"
)

func HandleSettingsCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.Settings(session), keyboards.Settings(session))
}

func HandleSettingsButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallSettingsBack:
		HandleSettingsCommand(app, session)

	case states.CallSettingsLanguage:
		requestLanguageSelect(app, session)

	case states.CallSettingsKinopoiskToken:
		requestKinopoiskToken(app, session)

	case states.CallSettingsFilmsPageSize:
		requestSettingsFilmsPageSize(app, session)

	case states.CallSettingsCollectionsPageSize:
		requestSettingsCollectionsPageSize(app, session)

	case states.CallSettingsObjectsPageSize:
		requestSettingsObjectsPageSize(app, session)

	default:
		if strings.HasPrefix(utils.ParseCallback(app.Update), states.SelectLang) {
			parseLanguageSelect(app, session)
		}
	}
}

func HandleSettingsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		returnToSettingsMenu(app, session)
		return
	}

	switch session.State {
	case states.AwaitSettingsKinopoiskToken:
		parser.ParseKinopoiskToken(app, session, HandleSettingsCommand)

	case states.AwaitSettingsFilmsPageSize:
		parser.ParseSettingsFilmsPageSize(app, session, requestSettingsFilmsPageSize, finishUpdatePageSize)

	case states.AwaitSettingsCollectionsPageSize:
		parser.ParseSettingsCollectionsPageSize(app, session, requestSettingsCollectionsPageSize, finishUpdatePageSize)

	case states.AwaitSettingsObjectsPageSize:
		parser.ParseSettingsObjectsPageSize(app, session, requestSettingsObjectsPageSize, finishUpdatePageSize)
	}
}

func requestLanguageSelect(app models.App, session *models.Session) {
	if languages, err := utils.ParseSupportedLanguages(app.Config.LocalesDir); err != nil {
		sl.Log.Error("failed to parse supported languages", slog.Any("error", err), slog.String("dir", app.Config.LocalesDir))
		app.SendMessage(messages.LanguagesFailure(session), keyboards.Back(session, states.CallMenuSettings))
	} else {
		app.SendMessage(messages.SettingsLanguage(session), keyboards.SettingsLanguageSelect(session, languages))
	}
}

func parseLanguageSelect(app models.App, session *models.Session) {
	session.Lang = strings.TrimPrefix(utils.ParseCallback(app.Update), states.SelectLang)
	app.SendMessage(messages.SettingsLanguageSuccess(session), nil)
	returnToSettingsMenu(app, session)
}

func requestKinopoiskToken(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestKinopoiskToken(session), keyboards.Cancel(session))
	session.SetState(states.AwaitSettingsKinopoiskToken)
}

func requestSettingsFilmsPageSize(app models.App, session *models.Session) {
	app.SendMessage(messages.SettingsPageSize(session, session.FilmsState.PageSize), keyboards.Cancel(session))
	session.SetState(states.AwaitSettingsFilmsPageSize)
}

func requestSettingsCollectionsPageSize(app models.App, session *models.Session) {
	app.SendMessage(messages.SettingsPageSize(session, session.CollectionsState.PageSize), keyboards.Cancel(session))
	session.SetState(states.AwaitSettingsCollectionsPageSize)
}
func requestSettingsObjectsPageSize(app models.App, session *models.Session) {
	app.SendMessage(messages.SettingsPageSize(session, session.CollectionFilmsState.PageSize), keyboards.Cancel(session))
	session.SetState(states.AwaitSettingsObjectsPageSize)
}

func finishUpdatePageSize(app models.App, session *models.Session) {
	app.SendMessage(messages.SettingsPageSizeSuccess(session), nil)
	returnToSettingsMenu(app, session)
}

func returnToSettingsMenu(app models.App, session *models.Session) {
	session.ClearState()
	HandleSettingsCommand(app, session)
}
