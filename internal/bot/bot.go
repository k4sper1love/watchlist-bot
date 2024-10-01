package bot

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/handlers"
)

func Run(cfg *config.Config) error {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	if err != nil {
		return err
	}
	bot.Debug = true

	sl.Log.Info(fmt.Sprintf("Authorized on account %s", bot.Self.UserName))

	updateConfig := tgbotapi.NewUpdate(0)
	updateConfig.Timeout = 60

	updates, err := bot.GetUpdatesChan(updateConfig)
	if err != nil {
		return err
	}

	for update := range updates {
		switch update.Message.Command() {
		case "start":
			go handlers.StartHandler(bot, &update)
		case "help":
			go handlers.HelpHandler(bot, &update)
		case "profile":
			go handlers.ProfileHandler(bot, &update)
		default:
			msg := tgbotapi.NewMessage(update.Message.Chat.ID, "The bot is currently under development")
			msg.ReplyToMessageID = update.Message.MessageID
			bot.Send(msg)
		}
	}

	return nil
}
