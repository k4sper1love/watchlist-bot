package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

func sendMessage(app config.App, text string) {
	chatID := parseChatID(app.Upd)
	msg := tgbotapi.NewMessage(chatID, text)
	_, err := app.Bot.Send(msg)
	if err != nil {
		sl.Log.Error("error sending the message", slog.Any("chat_id", chatID))
		return
	}
}

func sendMessageWithKeyboard(app config.App, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	chatID := parseChatID(app.Upd)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard

	_, err := app.Bot.Send(msg)
	if err != nil {
		sl.Log.Error("error sending the message", slog.Any("chat_id", chatID), slog.Any("error", err))
		return
	}
}

func sendImage(app config.App, imageURL string) {
	chatID := parseChatID(app.Upd)

	image, err := createImageMessage(chatID, imageURL)
	if err != nil {
		sl.Log.Error("error creating the image", slog.Any("chat_id", chatID), slog.Any("error", err))
		return
	}

	_, err = app.Bot.Send(image)
	if err != nil {
		sl.Log.Error("error sending the image", slog.Any("chat_id", chatID), slog.Any("error", err))
		return
	}
}

func sendImageWithText(app config.App, imageURL, text string) {
	chatID := parseChatID(app.Upd)

	image, err := createImageMessage(chatID, imageURL)
	if err != nil {
		sl.Log.Error("error creating the image", slog.Any("chat_id", chatID), slog.Any("error", err))
		return
	}

	image.Caption = text

	_, err = app.Bot.Send(image)
	if err != nil {
		sl.Log.Error("error sending the image", slog.Any("chat_id", chatID), slog.Any("error", err))
		return
	}
}

func sendImageWithTextAndKeyboard(app config.App, imageURL, text string, keyboard tgbotapi.InlineKeyboardMarkup) {
	chatID := parseChatID(app.Upd)

	image, err := createImageMessage(chatID, imageURL)
	if err != nil {
		sl.Log.Error("error creating the image", slog.Any("chat_id", chatID), slog.Any("error", err))
		return
	}

	image.Caption = text
	image.ReplyMarkup = keyboard

	_, err = app.Bot.Send(image)
	if err != nil {
		sl.Log.Error("error sending the image", slog.Any("chat_id", chatID), slog.Any("error", err))
		return
	}
}

func parseTelegramID(update *tgbotapi.Update) int {
	if update.Message != nil {
		return update.Message.From.ID
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}
	return -1
}

func parseChatID(update *tgbotapi.Update) int64 {
	if update.Message != nil {
		return update.Message.Chat.ID
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Chat.ID
	}
	return -1
}

func parseMessageText(update *tgbotapi.Update) string {
	if update.Message == nil {
		return ""
	}

	return update.Message.Text
}

func parseCallbackData(update *tgbotapi.Update) string {
	if update.CallbackQuery == nil {
		return ""
	}

	return update.CallbackQuery.Data
}

func clearSession(session *models.Session) {
	session.UserID = -1
	session.AccessToken = ""
	session.RefreshToken = ""

	session.CollectionState.Name = ""
	session.CollectionState.Description = ""
	session.CollectionState.CurrentPage = 0
	session.CollectionState.LastPage = 0
	resetState(session)
}

func createImageMessage(chatID int64, imageURL string) (*tgbotapi.PhotoConfig, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	splitURL := strings.Split(imageURL, "/")
	filename := splitURL[len(splitURL)-1]
	path := filepath.Join("static", filename)

	if _, err := os.Stat("static"); os.IsNotExist(err) {
		err = os.Mkdir("static", os.ModePerm)
		if err != nil {
			return nil, err
		}
	}

	file, err := os.Create(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return nil, err
	}

	image := tgbotapi.NewPhotoUpload(chatID, filename)

	return &image, nil
}
