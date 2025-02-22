package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/handlers"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/logger"
	"log/slog"
)

func Start(app *models.App) error {
	bot, err := tgbotapi.NewBotAPI(app.Config.BotToken)
	if err != nil {
		sl.Log.Error("failed to create bot", slog.Any("error", err))
		return err
	}

	bot.Debug = false
	app.Bot = bot
	sl.Log.Info("bot authorized", slog.String("username", bot.Self.UserName))

	updates, err := fetchUpdates(bot)
	if err != nil {
		return err
	}

	processUpdates(app, updates)
	return nil
}

func fetchUpdates(bot *tgbotapi.BotAPI) (tgbotapi.UpdatesChannel, error) {
	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		sl.Log.Error("failed to get updates", slog.Any("error", err))
		return nil, err
	}
	return updates, nil
}

func processUpdates(app *models.App, updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if utils.IsBotMessage(&update) {
			continue
		}

		updateAppContext(app, &update)
		go handlers.HandleUpdates(*app)
	}
}

func updateAppContext(app *models.App, update *tgbotapi.Update) {
	userID := utils.ParseTelegramID(update)

	if userID != app.Bot.Self.ID {
		app.FileLogger = logger.GetLogger(userID)
	}

	app.Upd = update
}
