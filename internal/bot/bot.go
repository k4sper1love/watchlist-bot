// Package bot provides functionality for initializing, running, and managing the Telegram bot.
//
// It is responsible for setting up the bot, handling updates, and integrating with
// external components such as the database, translator, and logger.
package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/config"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/logger"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"log/slog"
)

// Run initializes and starts the bot application.
// It loads the configuration, sets up logging, connects to the database,
// initializes the translator, and starts the bot.
func Run() error {
	slog.Info("starting application...")

	// Load application configuration
	app, err := config.Init()
	if err != nil {
		return err
	}
	slog.Info("application config loaded successfully")

	// Initialize structured logging
	sl.Init(app.Config.Environment)

	// Connect to the PostgreSQL database
	if err = postgres.ConnectDatabase(app.Config.DatabaseURL); err != nil {
		return err
	}
	slog.Info("database connection established successfully")

	// Initialize the translator with locale directory
	if err = translator.Init(app.Config.LocalesDir); err != nil {
		return err
	}
	slog.Info("translator initialized successfully")

	// Start the Telegram bot
	return startBot(app)
}

// startBot initializes the Telegram bot API and starts processing updates.
// It authorizes the bot, fetches updates, and processes them in a loop.
func startBot(app *models.App) error {
	bot, err := tgbotapi.NewBotAPI(app.Config.BotToken)
	if err != nil {
		slog.Error("failed to create bot", slog.Any("error", err))
		return err
	}

	bot.Debug = false
	app.Bot = bot
	slog.Info("bot authorized", slog.String("username", bot.Self.UserName))

	updates, err := fetchUpdates(bot)
	if err != nil {
		return err
	}

	processUpdates(app, updates)
	return nil
}

// fetchUpdates retrieves updates from the Telegram bot API.
// It configures the update channel with a timeout and returns the channel for processing.
func fetchUpdates(bot *tgbotapi.BotAPI) (tgbotapi.UpdatesChannel, error) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		slog.Error("failed to get updates", slog.Any("error", err))
		return nil, err
	}
	return updates, nil
}

// processUpdates processes incoming updates from the Telegram bot API.
// It filters out bot messages and delegates handling to the appropriate handlers.
func processUpdates(app *models.App, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if utils.IsBotMessage(&update) {
			continue
		}

		updateAppContext(app, &update)
		go handlers.HandleUpdates(*app)
	}
}

// updateAppContext updates the application context based on the current Telegram update.
// It sets the logger and updates the app's state with the latest update.
func updateAppContext(app *models.App, update *tgbotapi.Update) {
	userID := utils.ParseTelegramID(update)

	if userID != app.Bot.Self.ID {
		app.Logger = logger.Get(userID)
	}

	app.Update = update
}
