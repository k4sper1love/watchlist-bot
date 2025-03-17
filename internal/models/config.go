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
	maxMessageLength = 3500 // Maximum length of a Telegram message.
	maxCaptionLength = 1000 // Maximum length of a caption for media (e.g., images).
)

// App represents the main application structure, encapsulating configuration, bot API, and logging.
type App struct {
	Config *Config          // Application configuration.
	Bot    *tgbotapi.BotAPI // Telegram bot API instance.
	Update *tgbotapi.Update // Incoming update from Telegram.
	Logger *logger.Wrapper  // Logger for recording application events.
}

// Config contains the application's configuration settings.
type Config struct {
	Version         string // Application version.
	BotToken        string // Telegram bot token.
	Environment     string // Environment (e.g., "local", "dev", "prod").
	DatabaseURL     string // Database connection URL.
	LocalesDir      string // Directory for localization files.
	APIHost         string // Host for external API requests.
	APISecret       string // Secret key for API verification.
	RootID          int    // Root user ID for admin purposes.
	YoutubeAPIToken string // YouTube API token.
	IMDBAPIToken    string // IMDB API token.
}

// MessageConfig defines the configuration for sending messages, including chat ID, message ID, text, and media.
type MessageConfig struct {
	ChatID    int64  // ID of the chat where the message will be sent.
	MessageID int    // ID of the message being sent or updated.
	NeedPin   bool   // Whether the message should be pinned.
	ImageURL  string // URL of the image to send.
	Text      string // Text content of the message.
	File      string // Path to a file to send.
}

// GetChatID retrieves the chat ID from the incoming update.
func (app App) GetChatID() int64 {
	if app.Update.Message != nil {
		return app.Update.Message.Chat.ID
	} else if app.Update.CallbackQuery != nil {
		return app.Update.CallbackQuery.Message.Chat.ID
	}
	return -1 // Return -1 if no valid chat ID is found.
}

// logWithPrefix sets a prefix for the logger and returns the updated logger.
func (app App) logWithPrefix(prefix string) *logger.Wrapper {
	app.Logger.SetPrefix(prefix)
	return app.Logger
}

// LogAsBot creates a logger with a bot-specific prefix.
func (app App) LogAsBot() *logger.Wrapper {
	return app.logWithPrefix(fmt.Sprintf("BOT %s (%s)", app.Bot.Self.UserName, app.Config.Version))
}

// LogAsUser creates a logger with a user-specific prefix.
func (app App) LogAsUser(id int) *logger.Wrapper {
	return app.logWithPrefix(fmt.Sprintf("USER %d", id))
}

// logMessage logs details about a sent message, including its text, image, and pin status.
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

// send sends a message using the Telegram bot API and logs the result.
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

// chunkTextAndSend splits long text into chunks and sends them as separate messages.
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

// sendMessageInternal sends a text message with optional keyboard markup.
func (app App) sendMessageInternal(config MessageConfig, keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewMessage(app.GetChatID(), config.Text)
	msg.ParseMode = "HTML"
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}
	app.send(msg, config)
}

// sendImageInternal sends an image with optional caption and keyboard markup.
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

// SendMessage sends a text message, splitting it into chunks if necessary.
func (app App) SendMessage(text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	app.chunkTextAndSend(text, keyboard)
}

// SendImage sends an image with optional caption and keyboard markup.
func (app App) SendImage(imageURL, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	imagePath, err := utils.DownloadImage(imageURL)
	if err != nil {
		app.handleDownloadImageError()
		return
	}
	defer utils.RemoveFile(imagePath)

	app.sendImageInternal(MessageConfig{Text: text, ImageURL: imageURL}, imagePath, keyboard)
}

// SendBroadcastMessage sends a broadcast text message to multiple users.
func (app App) SendBroadcastMessage(ids []int, needPin bool, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	for _, id := range ids {
		app.createTemp(id).sendMessageInternal(MessageConfig{Text: text, NeedPin: needPin}, keyboard)
	}
}

// SendBroadcastImage sends a broadcast image with optional caption to multiple users.
func (app App) SendBroadcastImage(ids []int, needPin bool, imageURL, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	imagePath, err := utils.DownloadImage(imageURL)
	if err != nil {
		app.handleDownloadImageError()
		return
	}
	defer utils.RemoveFile(imagePath)

	for _, id := range ids {
		app.createTemp(id).sendImageInternal(MessageConfig{NeedPin: needPin, Text: text, ImageURL: imageURL}, imagePath, keyboard)
	}
}

// SendMessageByID sends a text message to a specific user by their ID.
func (app App) SendMessageByID(id int, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	app.createTemp(id).SendMessage(text, keyboard)
}

// SendImageByID sends an image to a specific user by their ID.
func (app App) SendImageByID(id int, imageURL, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	imagePath, err := utils.DownloadImage(imageURL)
	if err != nil {
		app.handleDownloadImageError()
		return
	}
	defer utils.RemoveFile(imagePath)

	app.createTemp(id).SendImage(imagePath, text, keyboard)
}

// SendFile sends a file with optional caption and keyboard markup.
func (app App) SendFile(filepath string, text string, keyboard *tgbotapi.InlineKeyboardMarkup) {
	msg := tgbotapi.NewDocumentUpload(app.GetChatID(), filepath)
	msg.ParseMode = "HTML"
	msg.Caption = text
	if keyboard != nil {
		msg.ReplyMarkup = keyboard
	}

	app.send(msg, MessageConfig{File: filepath, Text: text})
}

// createTemp creates a temporary App instance for sending messages to a specific user.
func (app App) createTemp(id int) *App {
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

// handleDownloadImageError handles errors that occur when downloading an image.
func (app App) handleDownloadImageError() {
	lang := utils.ParseLanguageCode(app.Update)
	msg := "🚨 " + translator.Translate(lang, "getImageFailure", nil, nil)
	app.SendMessage(msg, nil)
}
