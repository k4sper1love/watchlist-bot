package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/admin"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/collectionFilms"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/collections"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/users"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"log"
	"log/slog"
	"strings"
)

func HandleUpdates(app models.App) {
	lang := utils.ParseLanguageCode(app.Upd)
	session, err := postgres.GetSessionByTelegramID(app)
	if err != nil {
		msg := translator.Translate(lang, "session_error", nil, nil)
		app.SendMessage(msg, nil)
		log.Println(err)
		return
	}

	if utils.ParseMessageCommand(app.Upd) == "reset" {
		session.Logout()
		general.RequireAuth(app, session, general.HandleStartCommand)
		return
	}

	if ok := general.Auth(app, session); !ok {
		postgres.SaveSessionWithDependencies(session)
		return
	}

	switch {
	case app.Upd.CallbackQuery != nil:
		handleCallbackQuery(app, session)

	case session.State == "":
		handleCommands(app, session)

	default:
		handleUserInput(app, session)
	}

	postgres.SaveSessionWithDependencies(session)
}

func handleCommands(app models.App, session *models.Session) {
	command := utils.ParseMessageCommand(app.Upd)
	callbackData := utils.ParseCallback(app.Upd)

	switch {
	case command == "start":
		general.HandleStartCommand(app, session)

	case command == "help":
		general.HandleHelpCommand(app, session)

	case command == "menu":
		general.HandleMenuCommand(app, session)

	case command == "profile" || callbackData == states.CallbackMenuSelectProfile:
		users.HandleProfileCommand(app, session)

	case command == "logout" || callbackData == states.CallbackMenuSelectLogout:
		general.HandleLogoutCommand(app, session)

	case command == "films" || callbackData == states.CallbackMenuSelectFilms:
		session.FilmsState.CurrentPage = 1
		session.SetContext(states.ContextFilm)
		films.HandleFilmsCommand(app, session)

	case command == "collections" || callbackData == states.CallbackMenuSelectCollections:
		session.CollectionsState.CurrentPage = 1
		collections.HandleCollectionsCommand(app, session)

	case command == "settings" || callbackData == states.CallbackMenuSelectSettings:
		general.HandleSettingsCommand(app, session)

	case command == "feedback" || callbackData == states.CallbackMenuSelectFeedback:
		general.HandleFeedbackCommand(app, session)

	case command == "admin" || callbackData == states.CallbackMenuSelectAdmin:
		general.RequireRole(app, session, admin.HandleMenuCommand, roles.Helper)

	default:
		msg := translator.Translate(session.Lang, "unknownCommand", nil, nil)
		app.SendMessage(msg, nil)
	}
}

func handleUserInput(app models.App, session *models.Session) {
	switch {
	case strings.HasPrefix(session.State, "logout_awaiting"):
		general.HandleLogoutProcess(app, session)

	case strings.HasPrefix(session.State, "admin_manage_users_awaiting"):
		general.RequireRole(app, session, admin.HandleUsersProcess, roles.Admin)

	case strings.HasPrefix(session.State, "admin_user_detail_awaiting"):
		general.RequireRole(app, session, admin.HandleUserDetailProcess, roles.Admin)

	case strings.HasPrefix(session.State, "admin_list_awaiting"):
		general.RequireRole(app, session, admin.HandleAdminsProcess, roles.Admin)

	case strings.HasPrefix(session.State, "admin_broadcast_awaiting_"):
		general.RequireRole(app, session, admin.HandleBroadcastProcess, roles.Admin)

	case strings.HasPrefix(session.State, "feedback_awaiting"):
		general.HandleFeedbackProcess(app, session)

	case strings.HasPrefix(session.State, "update_profile_awaiting"):
		users.HandleUpdateProfileProcess(app, session)

	case strings.HasPrefix(session.State, "filters_films_awaiting"):
		films.HandleFiltersFilmsProcess(app, session)

	case strings.HasPrefix(session.State, "sorting_films_awaiting"):
		films.HandleSortingFilmsProcess(app, session)

	case strings.HasPrefix(session.State, "find_films_awaiting"):
		films.HandleFilmsProcess(app, session)

	case strings.HasPrefix(session.State, "new_film_awaiting"):
		films.HandleNewFilmProcess(app, session)

	case strings.HasPrefix(session.State, "update_film_awaiting"):
		films.HandleUpdateFilmProcess(app, session)

	case strings.HasPrefix(session.State, "viewed_film_awaiting"):
		films.HandleViewedFilmProcess(app, session)

	case strings.HasPrefix(session.State, "delete_film_awaiting"):
		films.HandleDeleteFilmProcess(app, session)

	case strings.HasPrefix(session.State, "delete_profile_awaiting"):
		users.HandleDeleteProfileProcess(app, session)

	case strings.HasPrefix(session.State, "sorting_collections_awaiting"):
		collections.HandleSortingCollectionsProcess(app, session)

	case strings.HasPrefix(session.State, "find_collections_awaiting"):
		collections.HandleCollectionProcess(app, session)

	case strings.HasPrefix(session.State, "new_collection_awaiting"):
		collections.HandleNewCollectionProcess(app, session)

	case strings.HasPrefix(session.State, "update_collection_awaiting"):
		collections.HandleUpdateCollectionProcess(app, session)

	case strings.HasPrefix(session.State, "delete_collection_awaiting"):
		collections.HandleDeleteCollectionProcess(app, session)

	case strings.HasPrefix(session.State, "settings_"):
		general.HandleSettingsProcess(app, session)

	default:
		msg := translator.Translate(session.Lang, "unknownState", nil, nil)
		app.SendMessage(msg, nil)
	}
}

func handleCallbackQuery(app models.App, session *models.Session) {
	callbackData := utils.ParseCallback(app.Upd)

	switch {
	case callbackData == states.CallbackMainMenu:
		general.HandleMenuCommand(app, session)

	case strings.HasPrefix(callbackData, "menu_select_"):
		handleCommands(app, session)

	case strings.HasPrefix(callbackData, "select_start_lang_"):
		general.HandleLanguageButton(app, session)

	case strings.HasPrefix(callbackData, "settings_") || strings.HasPrefix(callbackData, "select_lang_"):
		general.HandleSettingsButton(app, session)

	case strings.HasPrefix(callbackData, "admin_select"):
		general.RequireRole(app, session, admin.HandleMenuButton, roles.Helper)

	case strings.HasPrefix(callbackData, "admin_detail_"):
		general.RequireRole(app, session, admin.HandleAdminDetailButtons, roles.SuperAdmin)

	case strings.HasPrefix(callbackData, "admin_manage_users_select") || strings.HasPrefix(callbackData, "select_admin_user_") ||
		strings.HasPrefix(callbackData, "admin_users_list_"):
		general.RequireRole(app, session, admin.HandleUsersButton, roles.Helper)

	case strings.HasPrefix(callbackData, "admin_user_detail_") || strings.HasPrefix(callbackData, "admin_user_role_select_"):
		general.RequireRole(app, session, admin.HandleUserDetailButton, roles.Helper)

	case strings.HasPrefix(callbackData, "admin_feedback_list_") || strings.HasPrefix(callbackData, "select_admin_feedback_"):
		general.RequireRole(app, session, admin.HandleFeedbacksButtons, roles.Helper)

	case strings.HasPrefix(callbackData, "admin_list_") || strings.HasPrefix(callbackData, "select_admin_"):
		general.RequireRole(app, session, admin.HandleAdminsButtons, roles.SuperAdmin)

	case strings.HasPrefix(callbackData, "admin_feedback_detail_"):
		general.RequireRole(app, session, admin.HandleFeedbackDetailButtons, roles.Helper)

	case strings.HasPrefix(callbackData, "feedback_category"):
		general.HandleFeedbackButtons(app, session)

	case strings.HasPrefix(callbackData, "profile_"):
		users.HandleProfileButtons(app, session)

	case strings.HasPrefix(callbackData, "update_profile_select"):
		users.HandleUpdateProfileButtons(app, session)

	case strings.HasPrefix(callbackData, "films_") || strings.HasPrefix(callbackData, "select_film_"):
		if session.Context == states.ContextFilm {
			films.HandleFilmsButtons(app, session, general.HandleMenuCommand)
		} else if session.Context == states.ContextCollection {
			films.HandleFilmsButtons(app, session, collections.HandleCollectionsCommand)
		}

	case strings.HasPrefix(callbackData, "filters_films_select"):
		films.HandleFiltersFilmsButtons(app, session)

	case strings.HasPrefix(callbackData, "sorting_films_select"):
		films.HandleSortingFilmsButtons(app, session)

	case strings.HasPrefix(callbackData, "find_films"):
		films.HandleFindFilmsButtons(app, session)

	case strings.HasPrefix(callbackData, "new_film_select"):
		films.HandleNewFilmButtons(app, session)

	case strings.HasPrefix(callbackData, "manage_film_select"):
		films.HandleManageFilmButtons(app, session)

	case strings.HasPrefix(callbackData, "update_film_select"):
		films.HandleUpdateFilmButtons(app, session)

	case strings.HasPrefix(callbackData, "film_detail"):
		films.HandleFilmsDetailButtons(app, session)

	case strings.HasPrefix(callbackData, "collections_") || strings.HasPrefix(callbackData, "select_collection_"):
		collections.HandleCollectionsButtons(app, session)

	case strings.HasPrefix(callbackData, "sorting_collections_select"):
		collections.HandleSortingCollectionsButtons(app, session)

	case strings.HasPrefix(callbackData, "find_collections"):
		collections.HandleFindCollectionsButtons(app, session)

	case strings.HasPrefix(callbackData, "manage_collection_select"):
		collections.HandleManageCollectionButtons(app, session)

	case strings.HasPrefix(callbackData, "update_collection_select"):
		collections.HandleUpdateCollectionButtons(app, session)

	case strings.HasPrefix(callbackData, "collection_films_"):
		session.CollectionFilmsState.CurrentPage = 1
		collectionFilms.HandleCollectionFilmsButtons(app, session)

	case strings.HasPrefix(callbackData, "options_film_to_collection_"):
		session.CollectionFilmsState.CurrentPage = 1
		collectionFilms.HandleOptionsFilmToCollectionButtons(app, session)

	case strings.HasPrefix(callbackData, "add_collection_to_film_") || strings.HasPrefix(callbackData, "select_cf_collection_"):
		collectionFilms.HandleAddCollectionToFilmButtons(app, session)

	case strings.HasPrefix(callbackData, "add_film_to_collection_") || strings.HasPrefix(callbackData, "select_cf_film_"):
		collectionFilms.HandleAddFilmToCollectionButtons(app, session)

	default:
		handleUserInput(app, session)
	}

	answerCallbackQuery(app, session)
}

func answerCallbackQuery(app models.App, session *models.Session) {
	_, err := app.Bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
		CallbackQueryID: app.Upd.CallbackQuery.ID,
		//Text:            "Обработка завершена",
		ShowAlert: false,
	})

	if err != nil {
		sl.Log.Error("Failed to answer callback", slog.Any("error", err))
		general.HandleMenuCommand(app, session)
		return
	}
}
