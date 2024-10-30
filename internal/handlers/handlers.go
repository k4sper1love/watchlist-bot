package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"log/slog"
	"strings"
)

func HandleUpdates(app config.App) {
	telegramID := parseTelegramID(app.Upd)
	session, err := postgres.GetSessionByTelegramID(telegramID)
	if err != nil {
		sendMessage(app, "Произошла ошибка при получении сессии")
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

	postgres.SaveSessionWihDependencies(session)
}

func handleCommands(app config.App, session *models.Session) {
	switch app.Upd.Message.Command() {
	case "start":
		handleStartCommand(app, session)

	case "help":
		handleHelpCommand(app, session)

	case "profile":
		requireAuth(app, session, handleProfileCommand)

	case "logout":
		requireAuth(app, session, handleLogoutCommand)

	case "collections":
		session.CollectionState.CurrentPage = 1
		requireAuth(app, session, handleCollectionsCommand)

	case "new_collection":
		requireAuth(app, session, handleNewCollectionCommand)

	case "settings":
		handleSettingCommand(app, session)

	default:
		sendMessage(app, "Неизвестная команда. Введите /help")
	}
}

func handleUserInput(app config.App, session *models.Session) {
	switch {
	case strings.HasPrefix(session.State, "logout_"):
		handleLogoutProcess(app, session)

	case strings.HasPrefix(session.State, "new_collection_"):
		handleNewCollectionProcess(app, session)

	case strings.HasPrefix(session.State, "settings_"):
		handleSettingProcess(app, session)

	default:
		sendMessage(app, "Неизвестное состояние")
	}
}

func handleCallbackQuery(app config.App, session *models.Session) {
	callbackData := parseCallbackData(app.Upd)

	switch {
	case callbackData == CallbackSettingsCollectionsPageSize:
		setState(session, callbackData)
		handleSettingButton(app, session)

	case callbackData == CallbackCollectionsNextPage || callbackData == CallbackCollectionsPrevPage || strings.HasPrefix(callbackData, "select_collection_"):
		setState(session, callbackData)
		requireAuth(app, session, handleCollectionsButtons)

	case callbackData == CallbackCollectionFilmsNextPage || callbackData == CallbackCollectionFilmsPrevPage || strings.HasPrefix(callbackData, "select_cf_"):
		setState(session, callbackData)
		requireAuth(app, session, handleCollectionFilmsButtons)

	case callbackData == CallbackCollectionFilmsDetailNextPage || callbackData == CallbackCollectionFilmsDetailPrevPage:
		setState(session, callbackData)
		requireAuth(app, session, handleCollectionFilmsDetailButtons)

	default:
		sendMessage(app, "Неизвестная команда")
		return
	}

	answerCallbackQuery(app)
}

func answerCallbackQuery(app config.App) {
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
