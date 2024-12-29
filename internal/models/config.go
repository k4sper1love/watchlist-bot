package models

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log"
	"log/slog"
	"os"
)

const (
	maxMessageLength = 4000
	maxCaptionLength = 1900
)

type App struct {
	Vars *Vars
	Bot  *tgbotapi.BotAPI
	Upd  *tgbotapi.Update
}

type Vars struct {
	Version           string
	BotToken          string
	Environment       string
	DSN               string
	Host              string
	Secret            string
	AdminID           int
	KinopoiskAPIToken string
	YoutubeAPIToken   string
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
	for len(text) > maxMessageLength {
		firstPart, secondPart := utils.SplitTextByLength(text, maxMessageLength)

		msg := tgbotapi.NewMessage(app.GetChatID(), firstPart)
		msg.ParseMode = "HTML"
		app.send(msg)

		text = secondPart
	}

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
		app.SendMessage("Error when sending the image", nil)
		return
	}

	defer func() {
		if err := os.Remove(imagePath); err != nil {
			log.Println("Failed to remove temp file", slog.Any("error", err))
		}
	}()

	msg := tgbotapi.NewPhotoUpload(app.GetChatID(), imagePath)
	msg.ParseMode = "HTML"

	if text != "" && len(text) < maxCaptionLength {
		msg.Caption = text
	}

	if keyboard != nil && len(text) < maxCaptionLength {
		msg.ReplyMarkup = keyboard
	}

	app.send(msg)

	if len(text) > maxCaptionLength {
		app.SendMessage(text, keyboard)
	}
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
