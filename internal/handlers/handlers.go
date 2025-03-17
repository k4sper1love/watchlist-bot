// Package handlers defines the core logic for processing user interactions in the Watchlist Telegram bot.
//
// It acts as the central hub for routing and handling updates, managing commands, and processing user input.
// This package ensures a seamless and structured user experience by centralizing all interaction logic.
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

// HandleUpdates processes incoming updates from the Telegram bot API.
// It logs the update, retrieves the user's session, and routes the update to the appropriate handler.
func HandleUpdates(app models.App) {
	logUpdate(app) // Log the incoming update for debugging purposes.

	session, err := postgres.GetSessionByTelegramID(app)
	if err != nil {
		// If there's an error retrieving the session, send an error message to the user.
		app.SendMessage(messages.SessionError(utils.ParseLanguageCode(app.Update)), nil)
		return
	}

	// Check if the user sent a "reset" command to reset their session state.
	if utils.ParseMessageCommand(app.Update) == "reset" {
		handleReset(app, session)
		return
	}

	// Ensure the user is authenticated before proceeding.
	general.RequireAuth(app, session, routeUpdate)

	// Save the session and its dependencies to the database.
	postgres.SaveSessionWithDependencies(session)
}

func routeUpdate(app models.App, session *models.Session) {
	switch {
	case app.Update.CallbackQuery != nil:
		// Handle callback queries (button interactions).
		handleCallbackQuery(app, session)

	case session.State == "":
		// Handle commands when no specific state is active.
		handleCommands(app, session)

	default:
		// Handle user input when a specific state is active.
		handleInput(app, session)
	}
}

// handleReset resets the user's session by logging them out and requiring re-authentication.
func handleReset(app models.App, session *models.Session) {
	session.Logout() // Clear the session data.
	postgres.SaveSessionWithDependencies(session)
	general.RequireAuth(app, session, general.HandleStartCommand) // Redirect the user to the start command.
}

// handleCommands processes user commands such as /start, /help, /menu, etc.
func handleCommands(app models.App, session *models.Session) {
	command := utils.ParseMessageCommand(app.Update) // Extract the command from the message.
	callbackData := utils.ParseCallback(app.Update)  // Extract the callback data if it's a button interaction.

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
		// Handle unknown commands by notifying the user.
		app.SendMessage(messages.UnknownCommand(session), nil)
	}
}

// handleInput processes user input based on the current session state.
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
		// Handle unknown states by notifying the user.
		app.SendMessage(messages.UnknownState(session), nil)
	}
}

// handleCallbackQuery processes callback queries (button interactions) from the Telegram bot.
func handleCallbackQuery(app models.App, session *models.Session) {
	callbackData := utils.ParseCallback(app.Update) // Extract the callback data from the update.

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
		// Handle film-related buttons based on the current session context (films or collections).
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
		// Handle unknown callback data as user input.
		handleInput(app, session)
	}

	// Acknowledge the callback query to prevent the "loading" indicator in Telegram.
	answerCallbackQuery(app)
}

// answerCallbackQuery sends an acknowledgment to Telegram for the callback query.
// This prevents the "loading" indicator from persisting in the Telegram interface.
func answerCallbackQuery(app models.App) {
	_, err := app.Bot.AnswerCallbackQuery(tgbotapi.CallbackConfig{
		CallbackQueryID: app.Update.CallbackQuery.ID,
	})
	if err != nil {
		// Log any errors that occur while answering the callback query.
		sl.Log.Error("failed to answer callback", slog.Any("error", err), slog.String("callback_id", app.Update.CallbackQuery.ID))
	}
}

// logUpdate logs incoming updates (messages or callback queries) for debugging and monitoring purposes.
func logUpdate(app models.App) {
	telegramID := utils.ParseTelegramID(app.Update) // Extract the Telegram user ID.
	messageID := utils.ParseMessageID(app.Update)   // Extract the message ID.
	input := fmt.Sprintf(" #%d: ", messageID)       // Format the log entry with the message ID.

	switch {
	case app.Update.Message != nil:
		// Log incoming messages.
		utils.LogUpdateInfo(telegramID, messageID, "message")
		input += fmt.Sprintf("(message) %s", utils.ParseMessageString(app.Update))

	case app.Update.CallbackQuery != nil:
		// Log incoming callback queries.
		utils.LogUpdateInfo(telegramID, messageID, "callback")
		input += fmt.Sprintf("(callback) %s", utils.ParseCallback(app.Update))

	default:
		// Log unknown update types.
		utils.LogUpdateInfo(telegramID, messageID, "unknown")
		input += "(unknown)"
	}

	// Log the update with the user's Telegram ID as the context.
	app.LogAsUser(telegramID).Print(input)
}
