package handlers

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/admin"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/collectionFilms"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/collections"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/profile"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"log/slog"
	"strings"
)

func HandleUpdates(app models.App) {
	logUpdate(app)

	session, err := postgres.GetSessionByTelegramID(app)
	if err != nil {
		app.SendMessage(messages.SessionError(utils.ParseLanguageCode(app.Update)), nil)
		return
	}

	if utils.ParseMessageCommand(app.Update) == "reset" {
		handleReset(app, session)
		return
	}

	if general.Auth(app, session) {
		routeUpdate(app, session)
	}

	postgres.SaveSessionWithDependencies(session)
}

func routeUpdate(app models.App, session *models.Session) {
	switch {
	case app.Update.CallbackQuery != nil:
		handleCallbackQuery(app, session)

	case session.State == "":
		handleCommands(app, session)

	default:
		handleInput(app, session)
	}
}

func handleReset(app models.App, session *models.Session) {
	session.Logout()
	postgres.SaveSessionWithDependencies(session)
	general.RequireAuth(app, session, general.HandleStartCommand)
}

func handleCommands(app models.App, session *models.Session) {
	command := utils.ParseMessageCommand(app.Update)
	callbackData := utils.ParseCallback(app.Update)

	switch {
	case command == "start":
		general.HandleStartCommand(app, session)

	case command == "help":
		general.HandleHelpCommand(app, session)

	case command == "menu":
		general.HandleMenuCommand(app, session)

	case command == "profile" || callbackData == states.CallMenuProfile:
		profile.HandleProfileCommand(app, session)

	case command == "logout" || callbackData == states.CallMenuLogout:
		general.HandleLogoutCommand(app, session)

	case command == "films" || callbackData == states.CallMenuFilms:
		session.FilmsState.CurrentPage = 1
		session.SetContext(states.CtxFilm)
		films.HandleFilmsCommand(app, session)

	case command == "collections" || callbackData == states.CallMenuCollections:
		session.CollectionsState.CurrentPage = 1
		collections.HandleCollectionsCommand(app, session)

	case command == "settings" || callbackData == states.CallMenuSettings:
		general.HandleSettingsCommand(app, session)

	case command == "feedback" || callbackData == states.CallMenuFeedback:
		general.HandleFeedbackCommand(app, session)

	case command == "admin" || callbackData == states.CallMenuAdmin:
		general.RequireRole(app, session, admin.HandleMenuCommand, roles.Helper)

	default:
		app.SendMessage(messages.UnknownCommand(session), nil)
	}
}

func handleInput(app models.App, session *models.Session) {
	switch {
	case strings.HasPrefix(session.State, states.LogoutAwait):
		general.HandleLogoutProcess(app, session)

	case strings.HasPrefix(session.State, states.UserDetailAwait):
		general.RequireRole(app, session, admin.HandleUserDetailProcess, roles.Admin)

	case strings.HasPrefix(session.State, states.EntitiesAwait):
		general.RequireRole(app, session, admin.HandleEntitiesProcess, roles.Admin)

	case strings.HasPrefix(session.State, states.BroadcastAwait):
		general.RequireRole(app, session, admin.HandleBroadcastProcess, roles.Admin)

	case strings.HasPrefix(session.State, states.FeedbackAwait):
		general.HandleFeedbackProcess(app, session)

	case strings.HasPrefix(session.State, states.UpdateProfileAwait):
		profile.HandleUpdateProfileProcess(app, session)

	case strings.HasPrefix(session.State, states.FilmFiltersAwait):
		films.HandleFilmFiltersProcess(app, session)

	case strings.HasPrefix(session.State, states.FilmSortingAwait):
		films.HandleSortingFilmsProcess(app, session)

	case strings.HasPrefix(session.State, states.FilmsAwait):
		films.HandleFilmsProcess(app, session)

	case strings.HasPrefix(session.State, states.NewFilmAwait):
		films.HandleNewFilmProcess(app, session)

	case strings.HasPrefix(session.State, states.UpdateFilmAwait):
		films.HandleUpdateFilmProcess(app, session)

	case strings.HasPrefix(session.State, states.ViewedFilmAwait):
		films.HandleViewedFilmProcess(app, session)

	case strings.HasPrefix(session.State, states.DeleteFilmAwait):
		films.HandleDeleteFilmProcess(app, session)

	case strings.HasPrefix(session.State, states.DeleteProfileAwait):
		profile.HandleDeleteProfileProcess(app, session)

	case strings.HasPrefix(session.State, states.CollectionSortingAwait):
		collections.HandleSortingCollectionsProcess(app, session)

	case strings.HasPrefix(session.State, states.CollectionsAwait):
		collections.HandleCollectionProcess(app, session)

	case strings.HasPrefix(session.State, states.NewCollectionAwait):
		collections.HandleNewCollectionProcess(app, session)

	case strings.HasPrefix(session.State, states.UpdateCollectionAwait):
		collections.HandleUpdateCollectionProcess(app, session)

	case strings.HasPrefix(session.State, states.DeleteCollectionAwait):
		collections.HandleDeleteCollectionProcess(app, session)

	case strings.HasPrefix(session.State, states.AddFilmToCollectionAwait):
		collectionFilms.HandleAddFilmToCollectionProcess(app, session)

	case strings.HasPrefix(session.State, states.AddCollectionToFilmAwait):
		collectionFilms.HandleAddCollectionToFilmProcess(app, session)

	case strings.HasPrefix(session.State, states.SettingsAwait):
		general.HandleSettingsProcess(app, session)

	default:
		app.SendMessage(messages.UnknownState(session), nil)
	}
}

func handleCallbackQuery(app models.App, session *models.Session) {
	callbackData := utils.ParseCallback(app.Update)

	switch {
	case callbackData == states.CallMainMenu:
		general.HandleMenuCommand(app, session)

	case strings.HasPrefix(callbackData, states.Menu):
		handleCommands(app, session)

	case strings.HasPrefix(callbackData, states.SelectStartLang):
		general.HandleLanguageButton(app, session)

	case strings.HasPrefix(callbackData, states.Settings) || strings.HasPrefix(callbackData, states.SelectLang):
		general.HandleSettingsButtons(app, session)

	case strings.HasPrefix(callbackData, states.Admin):
		general.RequireRole(app, session, admin.HandleMenuButton, roles.Helper)

	case strings.HasPrefix(callbackData, states.AdminDetail):
		general.RequireRole(app, session, admin.HandleAdminDetailButtons, roles.SuperAdmin)

	case strings.HasPrefix(callbackData, states.UserDetail):
		general.RequireRole(app, session, admin.HandleUserDetailButton, roles.Helper)

	case strings.HasPrefix(callbackData, states.Feedbacks) || strings.HasPrefix(callbackData, states.SelectFeedback):
		general.RequireRole(app, session, admin.HandleFeedbacksButtons, roles.Helper)

	case strings.HasPrefix(callbackData, states.Entities) || strings.HasPrefix(callbackData, states.SelectEntity):
		general.RequireRole(app, session, admin.HandleEntitiesButtons, roles.Helper)

	case strings.HasPrefix(callbackData, states.FeedbackDetail):
		general.RequireRole(app, session, admin.HandleFeedbackDetailButtons, roles.Helper)

	case strings.HasPrefix(callbackData, states.FeedbackCategory):
		general.HandleFeedbackButtons(app, session)

	case strings.HasPrefix(callbackData, states.Profile):
		profile.HandleProfileButtons(app, session)

	case strings.HasPrefix(callbackData, states.UpdateProfile):
		profile.HandleUpdateProfileButtons(app, session)

	case strings.HasPrefix(callbackData, states.Films) || strings.HasPrefix(callbackData, states.SelectFilm):
		if session.Context == states.CtxFilm {
			films.HandleFilmsButtons(app, session, general.HandleMenuCommand)
		} else if session.Context == states.CtxCollection {
			films.HandleFilmsButtons(app, session, collections.HandleCollectionsCommand)
		}

	case strings.HasPrefix(callbackData, states.FilmFilters):
		films.HandleFilmFiltersButtons(app, session)

	case strings.HasPrefix(callbackData, states.FilmSorting):
		films.HandleSortingFilmsButtons(app, session)

	case strings.HasPrefix(callbackData, states.FindFilms):
		films.HandleFindFilmsButtons(app, session)

	case strings.HasPrefix(callbackData, states.FindNewFilm) || strings.HasPrefix(callbackData, states.SelectNewFilm):
		films.HandleFindNewFilmButtons(app, session)

	case strings.HasPrefix(callbackData, states.NewFilm):
		films.HandleNewFilmButtons(app, session)

	case strings.HasPrefix(callbackData, states.ManageFilm):
		films.HandleManageFilmButtons(app, session)

	case strings.HasPrefix(callbackData, states.UpdateFilm):
		films.HandleUpdateFilmButtons(app, session)

	case strings.HasPrefix(callbackData, states.FilmDetail):
		films.HandleFilmDetailButtons(app, session)

	case strings.HasPrefix(callbackData, states.Collections) || strings.HasPrefix(callbackData, states.SelectCollection):
		collections.HandleCollectionsButtons(app, session)

	case strings.HasPrefix(callbackData, states.CollectionSorting):
		collections.HandleSortingCollectionsButtons(app, session)

	case strings.HasPrefix(callbackData, states.FindCollections):
		collections.HandleFindCollectionsButtons(app, session)

	case strings.HasPrefix(callbackData, states.ManageCollection):
		collections.HandleManageCollectionButtons(app, session)

	case strings.HasPrefix(callbackData, states.UpdateCollection):
		collections.HandleUpdateCollectionButtons(app, session)

	case strings.HasPrefix(callbackData, states.CollectionFilmsFrom):
		session.CollectionFilmsState.CurrentPage = 1
		collectionFilms.HandleCollectionFilmsButtons(app, session)

	case strings.HasPrefix(callbackData, states.FilmToCollectionOption):
		session.CollectionFilmsState.CurrentPage = 1
		collectionFilms.HandleOptionsFilmToCollectionButtons(app, session)

	case strings.HasPrefix(callbackData, states.AddCollectionToFilm) || strings.HasPrefix(callbackData, states.SelectCFCollection):
		collectionFilms.HandleAddCollectionToFilmButtons(app, session)

	case strings.HasPrefix(callbackData, states.AddFilmToCollection) || strings.HasPrefix(callbackData, states.SelectCFFilm):
		collectionFilms.HandleAddFilmToCollectionButtons(app, session)

	default:
		handleInput(app, session)
	}

	answerCallbackQuery(app)
}

func answerCallbackQuery(app models.App) {
	_, err := app.Bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
		CallbackQueryID: app.Update.CallbackQuery.ID,
	})
	if err != nil {
		sl.Log.Error("failed to answer callback", slog.Any("error", err), slog.String("callback_id", app.Update.CallbackQuery.ID))
	}
}

func logUpdate(app models.App) {
	telegramID := utils.ParseTelegramID(app.Update)
	messageID := utils.ParseMessageID(app.Update)
	input := fmt.Sprintf(" #%d: ", messageID)

	switch {
	case app.Update.Message != nil:
		utils.LogUpdateInfo(telegramID, messageID, "message")
		input += fmt.Sprintf("(message) %s", utils.ParseMessageString(app.Update))

	case app.Update.CallbackQuery != nil:
		utils.LogUpdateInfo(telegramID, messageID, "callback")
		input += fmt.Sprintf("(callback) %s", utils.ParseCallback(app.Update))

	default:
		utils.LogUpdateInfo(telegramID, messageID, "unknown")
		input += "(unknown)"
	}

	app.LogAsUser(telegramID).Print(input)
}
