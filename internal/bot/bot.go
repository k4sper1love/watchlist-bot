package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/handlers"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func Run(app *models.App) error {
	bot, err := tgbotapi.NewBotAPI(app.Vars.BotToken)
	if err != nil {
		return err
	}
	bot.Debug = false

	app.Bot = bot

	sl.Log.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		return err
	}

	for update := range updates {
		app.Upd = &update
		go handlers.HandleUpdates(*app)
	}

	return nil
}
