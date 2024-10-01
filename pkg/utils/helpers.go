package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"log/slog"
)

func SendMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update, text string) {
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)

	_, err := bot.Send(msg)
	if err != nil {
		sl.Log.Error("error sending the message", slog.Any("chat_id", update.Message.Chat.ID))
	}
}
