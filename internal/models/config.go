package models

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/logger"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"unicode/utf8"
)

const (
	maxMessageLength = 3500
	maxCaptionLength = 1000
)

type App struct {
	Config *Config
	Bot    *tgbotapi.BotAPI
	Update *tgbotapi.Update
	Logger *logger.Wrapper
}

type Config struct {
	Version         string
	BotToken        string
	Environment     string
	DatabaseURL     string
	APIHost         string
	APISecret       string
	RootID          int
	YoutubeAPIToken string
	IMDBAPIToken    string
}

type MessageConfig struct {
	ChatID    int64
	MessageID int
	NeedPin   bool
	ImageURL  string
	Text      string
	File      string
}

func (app App) GetChatID() int64 {
	if app.Update.Message != nil {
		return app.Update.Message.Chat.ID
	} else if app.Update.CallbackQuery != nil {
		return app.Update.CallbackQuery.Message.Chat.ID
	}
	return -1
}

func (app App) logWithPrefix(prefix string) *logger.Wrapper {
	app.Logger.SetPrefix(prefix)
	return app.Logger
}

func (app App) LogAsBot() *logger.Wrapper {
	return app.logWithPrefix(fmt.Sprintf("BOT %s (%s)", app.Bot.Self.UserName, app.Config.Version))
}

func (app App) LogAsUser(id int) *logger.Wrapper {
	return app.logWithPrefix(fmt.Sprintf("USER %d", id))
}

func (app App) logMessage(config MessageConfig) {
	utils.LogMessageInfo(config.ChatID, config.MessageID, config.Text != "", config.ImageURL != "", config.NeedPin)

	if app.Logger == nil {
		utils.LogMessageError(fmt.Errorf("logger is empty"), config.ChatID, config.MessageID)
		return
	}

	logStr := fmt.Sprintf(" #%d: ", config.MessageID)
	if config.NeedPin {
		logStr += "\n[PINNED]"
	}
	if config.File != "" {
		logStr += fmt.Sprintf("\n[file] %s", config.File)
	}
	if config.ImageURL != "" {
		logStr += fmt.Sprintf("\n[image] %s", config.ImageURL)
	}
	if config.Text != "" {
		logStr += fmt.Sprintf("\n%s", config.Text)
	}

	app.LogAsBot().Print(logStr)
}

func (app App) send(msg tgbotapi.Chattable, config MessageConfig) {
	if msg == nil {
		utils.LogMessageError(fmt.Errorf("message is empty"), app.GetChatID(), -1)
		return
	}

	sentMsg, err := app.Bot.Send(msg)
	if err != nil {
		utils.LogMessageError(err, app.GetChatID(), -1)
		return
	}

	config.ChatID = sentMsg.Chat.ID
	config.MessageID = sentMsg.MessageID

	if config.NeedPin {
		_, err = app.Bot.PinChatMessage(tgbotapi.PinChatMessageConfig{
			ChatID:    config.ChatID,
			MessageID: config.MessageID,
		})
		if err != nil {
			utils.LogMessageError(fmt.Errorf("failed to pin message: %v", err), config.ChatID, config.MessageID)
		}
	}

	app.logMessage(config)
}

func (app App) chunkTextAndSend(text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	iterationLimit := 100
	for utf8.RuneCountInString(text) > maxMessageLength && iterationLimit > 0 {
		firstPart, secondPart := utils.SplitTextByLength(text, maxMessageLength)
		if firstPart == "" {
			utils.LogMessageError(fmt.Errorf("failed to chuck text: first part is empty"), app.GetChatID(), -1)
			msg := "🚨 " + translator.Translate("en", "chunkTextError", nil, nil)
			app.SendMessage(msg, nil)
			return
		}
		app.sendMessageInternal(MessageConfig{Text: firstPart}, nil)
		text = secondPart
		iterationLimit--
	}
	if text != "" {
		app.sendMessageInternal(MessageConfig{Text: text}, keyboard)
	}
}

func (app App) sendMessageInternal(config MessageConfig, keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(app.GetChatID(), config.Text)
	msg.ParseMode = "HTML"
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	app.send(msg, config)
}

func (app App) sendImageInternal(config MessageConfig, imagePath string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewPhotoUpload(app.GetChatID(), imagePath)
	msg.ParseMode = "HTML"
	runeLen := utf8.RuneCountInString(config.Text)

	if runeLen <= maxCaptionLength {
		msg.Caption = config.Text
		if keyboard != nil {
			msg.ReplyMarkup = keyboard
		}
	}
	text := config.Text
	config.Text = msg.Caption

	app.send(msg, config)

	if runeLen > maxCaptionLength {
		app.chunkTextAndSend(text, keyboard)
	}
}

func (app App) SendMessage(text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	app.chunkTextAndSend(text, keyboard)
}

func (app App) SendImage(imageURL, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	imagePath, err := utils.DownloadImage(imageURL)
	if err != nil {
		app.handleDownloadImageError()
		return
	}
	defer utils.RemoveFile(imagePath)

	app.sendImageInternal(MessageConfig{Text: text, ImageURL: imageURL}, imagePath, keyboard)
}

func (app App) SendBroadcastMessage(ids []int, needPin bool, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	for _, id := range ids {
		app.createTempApp(id).sendMessageInternal(MessageConfig{Text: text, NeedPin: needPin}, keyboard)
	}
}

func (app App) SendBroadcastImage(ids []int, needPin bool, imageURL, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	imagePath, err := utils.DownloadImage(imageURL)
	if err != nil {
		app.handleDownloadImageError()
		return
	}
	defer utils.RemoveFile(imagePath)

	for _, id := range ids {
		app.createTempApp(id).sendImageInternal(MessageConfig{NeedPin: needPin, Text: text, ImageURL: imageURL}, imagePath, keyboard)
	}
}

func (app App) SendMessageByID(id int, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	app.createTempApp(id).SendMessage(text, keyboard)
}

func (app App) SendImageByID(id int, imageURL, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	imagePath, err := utils.DownloadImage(imageURL)
	if err != nil {
		app.handleDownloadImageError()
		return
	}
	defer utils.RemoveFile(imagePath)

	app.createTempApp(id).SendImage(imagePath, text, keyboard)
}

func (app App) SendFile(filepath string, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewDocumentUpload(app.GetChatID(), filepath)
	msg.ParseMode = "HTML"
	msg.Caption = text
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}

	app.send(msg, MessageConfig{File: filepath, Text: text})
}

func (app App) createTempApp(id int) *App {
	return &App{
		Bot:    app.Bot,
		Config: app.Config,
		Update: &tgbotapi.Update{
			Message: &tgbotapi.Message{
				Chat: &tgbotapi.Chat{
					ID: int64(id)},
			},
		},
		Logger: logger.GetLogger(id),
	}
}

func (app App) handleDownloadImageError() {
	lang := utils.ParseLanguageCode(app.Update)
	msg := "🚨 " + translator.Translate(lang, "getImageFailure", nil, nil)
	app.SendMessage(msg, nil)
}
