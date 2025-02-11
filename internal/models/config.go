package models

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/logger"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"log"
	"log/slog"
	"os"
	"unicode/utf8"
)

const (
	maxMessageLength = 3500
	maxCaptionLength = 1000
)

type App struct {
	Vars       *Vars
	Bot        *tgbotapi.BotAPI
	Upd        *tgbotapi.Update
	FileLogger *logger.Wrapper
}

type Vars struct {
	Version           string
	BotToken          string
	Environment       string
	DSN               string
	Host              string
	Secret            string
	RootID            int
	KinopoiskAPIToken string
	YoutubeAPIToken   string
	IMDBAPIToken      string
}

type LogConfig struct {
	NeedPin   bool
	ImageURL  string
	Text      string
	File      string
	MessageID int
}

func (app App) SetUserPrefix(userID int) *logger.Wrapper {
	app.FileLogger.SetPrefix(fmt.Sprintf("USER %d: ", userID))
	return app.FileLogger
}

func (app App) SetBotPrefix() *logger.Wrapper {
	prefix := fmt.Sprintf("BOT %s", app.Bot.Self.UserName)
	if app.Vars.Version != "" {
		prefix += fmt.Sprintf(" (%s)", app.Vars.Version)
	}
	prefix += ": "

	app.FileLogger.SetPrefix(prefix)

	return app.FileLogger
}

func (app App) send(msg tgbotapi.Chattable, config LogConfig) {
	if msg == nil {
		sl.Log.Error("error sending the message", slog.Any("error", fmt.Errorf("msg is nil")))
		return
	}

	sentMsg, err := app.Bot.Send(msg)
	if err != nil {
		sl.Log.Error("error sending the message", slog.Any("error", err))
		return
	}
	config.MessageID = sentMsg.MessageID

	if config.NeedPin {
		pinConfig := tgbotapi.PinChatMessageConfig{
			ChatID:    sentMsg.Chat.ID,
			MessageID: sentMsg.MessageID,
		}

		_, err = app.Bot.PinChatMessage(pinConfig)
		if err != nil {
			sl.Log.Error("error pinning the message", slog.Any("error", err))
		}
	}

	if app.FileLogger == nil {
		sl.Log.Error("app.FileLogger is nil")
		return
	}

	app.LogMessage(config)
}

func (app App) LogMessage(config LogConfig) {
	response := fmt.Sprintf("#%d", config.MessageID)
	if config.NeedPin {
		response += "\n[PINNED]"
	}
	if config.File != "" {
		response += fmt.Sprintf("\n[file] %s", config.File)
	}
	if config.ImageURL != "" {
		response += fmt.Sprintf("\n[image] %s", config.ImageURL)
	}
	if config.Text != "" {
		response += fmt.Sprintf("\n%s", config.Text)
	}

	if app.FileLogger == nil {
		sl.Log.Error("app.FileLogger is nil")
		return
	}

	app.SetBotPrefix().Printf(response)
}

func (app App) GetChatID() int64 {
	if app.Upd.Message != nil {
		return app.Upd.Message.Chat.ID
	} else if app.Upd.CallbackQuery != nil {
		return app.Upd.CallbackQuery.Message.Chat.ID
	}
	return -1
}

func (app App) chunkTextAndSend(text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	iterationLimit := 100
	for utf8.RuneCountInString(text) > maxMessageLength && iterationLimit > 0 {
		firstPart, secondPart := utils.SplitTextByLength(text, maxMessageLength)
		if len(firstPart) == 0 {
			errorMsg := translator.Translate("ru", "chunkTextError", nil, nil)
			log.Println(errorMsg)
			app.SendMessage(errorMsg, nil)
			return
		}
		app.createAndSendMessage(firstPart, nil)
		text = secondPart
		iterationLimit--
	}
	if len(text) > 0 {
		app.createAndSendMessage(text, keyboard)
	}
}

func (app App) createAndSendMessage(text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(app.GetChatID(), text)
	msg.ParseMode = "HTML"
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}

	app.send(msg, LogConfig{Text: text})
}

func (app App) SendMessage(text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	app.chunkTextAndSend(text, keyboard)
}

func (app App) sendImageInternal(config LogConfig, imagePath string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewPhotoUpload(app.GetChatID(), imagePath)
	msg.ParseMode = "HTML"
	runeLen := utf8.RuneCountInString(config.Text)
	if config.Text != "" && runeLen < maxCaptionLength {
		msg.Caption = config.Text
	}
	if keyboard != nil && runeLen < maxCaptionLength {
		msg.ReplyMarkup = keyboard
	}

	app.send(msg, LogConfig{NeedPin: config.NeedPin, Text: msg.Caption, ImageURL: config.ImageURL})

	if runeLen > maxCaptionLength {
		app.chunkTextAndSend(config.Text, keyboard)
	}
}

func (app App) SendImage(imageURL, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	imagePath, err := utils.DownloadImage(imageURL)
	if err != nil {
		log.Println("Error when sending the image")
		return
	}
	defer func() {
		if err := os.Remove(imagePath); err != nil {
			log.Println("Failed to remove temp file", slog.Any("error", err))
		}
	}()

	app.sendImageInternal(LogConfig{Text: text, ImageURL: imageURL}, imagePath, keyboard)
}

func (app App) SendBroadcastMessage(telegramIDs []int, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	for _, telegramID := range telegramIDs {
		msg := tgbotapi.NewMessage(int64(telegramID), text)
		if keyboard != nil {
			msg.ReplyMarkup = keyboard
		}
		app.FileLogger = logger.GetLogger(telegramID)
		app.send(msg, LogConfig{NeedPin: true, Text: text})
	}
}

func (app App) SendBroadcastImage(telegramIDs []int, imageURL, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	imagePath, err := utils.DownloadImage(imageURL)
	if err != nil {
		app.SendBroadcastMessage(telegramIDs, "Error when sending the image", nil)
		return
	}
	defer func() {
		if err := os.Remove(imagePath); err != nil {
			log.Println("Failed to remove temp file", slog.Any("error", err))
		}
	}()

	for _, telegramID := range telegramIDs {
		app.FileLogger = logger.GetLogger(telegramID)
		app.Upd = &tgbotapi.Update{Message: &tgbotapi.Message{Chat: &tgbotapi.Chat{ID: int64(telegramID)}}}
		app.sendImageInternal(LogConfig{NeedPin: true, Text: text, ImageURL: imageURL}, imagePath, keyboard)
	}
}

func (app App) SendMessageByID(id int, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	tempApp := App{
		Bot:  app.Bot,
		Vars: app.Vars,
		Upd: &tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: int64(id)},
			},
		},
		FileLogger: logger.GetLogger(id),
	}

	tempApp.SendMessage(text, keyboard)
}

func (app App) SendImageByID(id int, imageURL, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	imagePath, err := utils.DownloadImage(imageURL)
	if err != nil {
		app.SendMessageByID(id, "Error when sending the image", nil)
		return
	}
	defer func() {
		if err := os.Remove(imagePath); err != nil {
			log.Println("Failed to remove temp file", slog.Any("error", err))
		}
	}()

	tempApp := App{
		Bot:  app.Bot,
		Vars: app.Vars,
		Upd: &tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{ID: int64(id)},
			},
		},
		FileLogger: logger.GetLogger(id),
	}
	tempApp.SendImage(imagePath, text, keyboard)
}

func (app App) SendFile(filepath string, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewDocumentUpload(app.GetChatID(), filepath)
	msg.ParseMode = "HTML"

	if text != "" {
		msg.Caption = text
	}
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}

	app.send(msg, LogConfig{File: filepath, Text: text})
}
