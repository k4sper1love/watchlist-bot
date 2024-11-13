package models

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log/slog"
	"os"
)

type App struct {
	Vars *Vars
	Bot  *tgbotapi.BotAPI
	Upd  *tgbotapi.Update
}

type Vars struct {
	BotToken    string
	Environment string
	DSN         string
	Host        string
	Secret      string
}

func (app App) send(msg tgbotapi.Chattable) {
	if msg == nil {
		sl.Log.Error("error sending the message", slog.Any("error", fmt.Errorf("msg is nil")))
		return
	}
	_, err := app.Bot.Send(msg)
	if err != nil {
		sl.Log.Error("error sending the message", slog.Any("error", err))
		return
	}
}

func (app App) GetChatID() int64 {
	if app.Upd.Message != nil {
		return app.Upd.Message.Chat.ID
	} else if app.Upd.CallbackQuery != nil {
		return app.Upd.CallbackQuery.Message.Chat.ID
	}
	return -1
}

func (app App) SendMessage(text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(app.GetChatID(), text)
	msg.ParseMode = "HTML"

	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}

	app.send(msg)
}

func (app App) SendImage(imageURL, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	imagePath, err := utils.DownloadImage(imageURL)
	if err != nil {
		app.SendMessage("Ошибка при отправке изображения", nil)
		return
	}

	msg := tgbotapi.NewPhotoUpload(app.GetChatID(), imagePath)
	msg.ParseMode = "HTML"

	if text != "" {
		msg.Caption = text
	}

	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}

	app.send(msg)

	os.Remove(imagePath)
}

func (app App) SendBroadcastMessage(telegramIDs []int, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	for _, telegramID := range telegramIDs {
		msg := tgbotapi.NewMessage(int64(telegramID), text)

		if keyboard != nil {
			msg.ReplyMarkup = keyboard
		}

		app.send(msg)
	}
}
