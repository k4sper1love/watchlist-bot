package models

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
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
	BaseURL     string
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
	imagePath, _ := downloadImage(imageURL)

	msg := tgbotapi.NewPhotoUpload(app.GetChatID(), imagePath)
	msg.ParseMode = "HTML"

	if text != "" {
		msg.Caption = text
	}

	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}

	app.send(msg)
}

func downloadImage(imageURL string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	splitURL := strings.Split(imageURL, "/")
	filename := splitURL[len(splitURL)-1]
	path := filepath.Join("static", filename)

	if _, err := os.Stat("static"); os.IsNotExist(err) {
		err = os.Mkdir("static", os.ModePerm)
		if err != nil {
			return "", err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return "", err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", err
	}

	return path, nil
}
