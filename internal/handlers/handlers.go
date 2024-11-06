package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/collections"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/users"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log"
	"log/slog"
	"strings"
)

func HandleUpdates(app models.App) {
	telegramID := utils.ParseTelegramID(app.Upd)
	session, err := postgres.GetSessionByTelegramID(telegramID)
	if err != nil {
		app.SendMessage("Произошла ошибка при получении сессии", nil)
		log.Println(err)
		return
	}

	switch {
	case app.Upd.CallbackQuery != nil:
		handleCallbackQuery(app, session)

	case app.Upd.Message.Command() == "reset":
		session.ClearState()

	case session.State == "":
		handleCommands(app, session)

	default:
		handleUserInput(app, session)
	}

	postgres.SaveSessionWihDependencies(session)
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
		general.RequireAuth(app, session, users.HandleProfileCommand)

	case command == "logout" || callbackData == states.CallbackMenuSelectLogout:
		general.RequireAuth(app, session, general.HandleLogoutCommand)

	case command == "collections" || callbackData == states.CallbackMenuSelectCollections:
		session.CollectionsState.CurrentPage = 1
		general.RequireAuth(app, session, collections.HandleCollectionsCommand)

	case command == "settings" || callbackData == states.CallbackMenuSelectSettings:
		general.HandleSettingsCommand(app, session)

	default:
		app.SendMessage("Неизвестная команда. Введите /help", nil)
	}
}

func handleUserInput(app models.App, session *models.Session) {
	switch {
	case strings.HasPrefix(session.State, "logout_awaiting"):
		general.HandleLogoutProcess(app, session)

	case strings.HasPrefix(session.State, "new_collection_awaiting"):
		collections.HandleNewCollectionProcess(app, session)

	case strings.HasPrefix(session.State, "update_collection_awaiting"):
		collections.HandleUpdateCollectionProcess(app, session)

	case strings.HasPrefix(session.State, "delete_collection_awaiting"):
		collections.HandleDeleteCollectionProcess(app, session)

	case strings.HasPrefix(session.State, "new_collection_film_awaiting"):
		collections.HandleNewCollectionFilmProcess(app, session)

	case strings.HasPrefix(session.State, "update_collection_film_awaiting"):
		collections.HandleUpdateCollectionFilmProcess(app, session)

	case strings.HasPrefix(session.State, "delete_collection_film_awaiting"):
		collections.HandleDeleteCollectionFilmProcess(app, session)

	case strings.HasPrefix(session.State, "settings_"):
		general.HandleSettingsProcess(app, session)

	default:
		app.SendMessage("Неизвестное состояние. Введите /reset для сброса.", nil)
	}
}

func handleCallbackQuery(app models.App, session *models.Session) {
	callbackData := utils.ParseCallback(app.Upd)

	switch {
	case strings.HasPrefix(callbackData, "menu_select_"):
		handleCommands(app, session)

	case strings.HasPrefix(callbackData, "settings_"):
		general.HandleSettingsButton(app, session)

	case strings.HasPrefix(callbackData, "profile_"):
		general.RequireAuth(app, session, users.HandleProfileButtons)

	case strings.HasPrefix(callbackData, "collections_") || strings.HasPrefix(callbackData, "select_collection_"):
		general.RequireAuth(app, session, collections.HandleCollectionsButtons)

	case strings.HasPrefix(callbackData, "manage_collection_select"):
		general.RequireAuth(app, session, collections.HandleManageCollectionButtons)

	case strings.HasPrefix(callbackData, "update_collection_select"):
		general.RequireAuth(app, session, collections.HandleUpdateCollectionButtons)

	case strings.HasPrefix(callbackData, "collection_films_") || strings.HasPrefix(callbackData, "select_cf_"):
		general.RequireAuth(app, session, collections.HandleCollectionFilmsButtons)

	case strings.HasPrefix(callbackData, "manage_collection_film_select"):
		general.RequireAuth(app, session, collections.HandleManageCollectionFilmButtons)

	case strings.HasPrefix(callbackData, "update_collection_film_select"):
		general.RequireAuth(app, session, collections.HandleUpdateCollectionFilmButtons)

	case strings.HasPrefix(callbackData, "collection_film_detail"):
		general.RequireAuth(app, session, collections.HandleCollectionFilmsDetailButtons)

	default:
		handleUserInput(app, session)
	}

	answerCallbackQuery(app)
}

func answerCallbackQuery(app models.App) {
	_, err := app.Bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
		CallbackQueryID: app.Upd.CallbackQuery.ID,
		//Text:            "Обработка завершена",
		ShowAlert: false,
	})

	if err != nil {
		sl.Log.Error("answer callback", slog.Any("error", err))
		panic(err)
	}
}
