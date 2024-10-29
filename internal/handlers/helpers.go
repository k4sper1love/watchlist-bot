package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"log/slog"
)

func sendMessage(app config.App, text string) {
	chatID := parseChatID(app.Upd)
	msg := tgbotapi.NewMessage(chatID, text)

	_, err := app.Bot.Send(msg)
	if err != nil {
		sl.Log.Error("error sending the message", slog.Any("chat_id", chatID))
	}
}

func sendMessageWithKeyboard(app config.App, keyboard tgbotapi.InlineKeyboardMarkup, text string) {

	chatID := parseChatID(app.Upd)

	msg := tgbotapi.NewMessage(chatID, text)
	msg.ReplyMarkup = keyboard

	_, err := app.Bot.Send(msg)
	if err != nil {
		sl.Log.Error("error sending the message", slog.Any("chat_id", chatID), slog.Any("error", err))
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
