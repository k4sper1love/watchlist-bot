package general

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strings"
)

// HandleSettingsCommand handles the command for displaying the settings menu.
// Sends a message with options to customize user preferences, such as language, page size, and Kinopoisk token.
func HandleSettingsCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.Settings(session), keyboards.Settings(session))
}

// HandleSettingsButtons handles button interactions related to the settings menu.
// Supports actions like going back, changing language, updating Kinopoisk token, or modifying page sizes.
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

// HandleSettingsProcess processes the workflow for updating settings.
// Handles states like awaiting input for Kinopoisk token or page size.
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

// requestLanguageSelect prompts the user to select a new language for the bot interface.
func requestLanguageSelect(app models.App, session *models.Session) {
	if languages, err := utils.ParseSupportedLanguages(app.Config.LocalesDir); err != nil {
		utils.LogParseLanguagesError(err, app.Config.LocalesDir)
		app.SendMessage(messages.LanguagesFailure(session), keyboards.Back(session, states.CallMenuSettings))
	} else {
		app.SendMessage(messages.SettingsLanguage(session), keyboards.SettingsLanguageSelect(session, languages))
	}
}

// parseLanguageSelect processes the user's selection of a new language.
func parseLanguageSelect(app models.App, session *models.Session) {
	session.Lang = strings.TrimPrefix(utils.ParseCallback(app.Update), states.SelectLang)
	app.SendMessage(messages.SettingsLanguageSuccess(session), nil)
	returnToSettingsMenu(app, session)
}

// requestKinopoiskToken prompts the user to enter their Kinopoisk API token.
func requestKinopoiskToken(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestKinopoiskToken(session), keyboards.Cancel(session))
	session.SetState(states.AwaitSettingsKinopoiskToken)
}

// requestSettingsFilmsPageSize prompts the user to update the page size for films.
func requestSettingsFilmsPageSize(app models.App, session *models.Session) {
	app.SendMessage(messages.SettingsPageSize(session, session.FilmsState.PageSize), keyboards.Cancel(session))
	session.SetState(states.AwaitSettingsFilmsPageSize)
}

// requestSettingsCollectionsPageSize prompts the user to update the page size for collections.
func requestSettingsCollectionsPageSize(app models.App, session *models.Session) {
	app.SendMessage(messages.SettingsPageSize(session, session.CollectionsState.PageSize), keyboards.Cancel(session))
	session.SetState(states.AwaitSettingsCollectionsPageSize)
}

// requestSettingsObjectsPageSize prompts the user to update the page size for collection objects.
func requestSettingsObjectsPageSize(app models.App, session *models.Session) {
	app.SendMessage(messages.SettingsPageSize(session, session.CollectionFilmsState.PageSize), keyboards.Cancel(session))
	session.SetState(states.AwaitSettingsObjectsPageSize)
}

// finishUpdatePageSize finalizes the process of updating the page size.
func finishUpdatePageSize(app models.App, session *models.Session) {
	app.SendMessage(messages.SettingsPageSizeSuccess(session), nil)
	returnToSettingsMenu(app, session)
}

// returnToSettingsMenu clears the session state and navigates back to the settings menu.
func returnToSettingsMenu(app models.App, session *models.Session) {
	session.ClearState()
	HandleSettingsCommand(app, session)
}
