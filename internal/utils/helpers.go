package utils

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"strconv"
)

func GetItemID(index, currentPage, pageSize int) int {
	return (index + 1) + ((currentPage - 1) * pageSize)
}

func ParseTelegramID(update *tgbotapi.Update) int {
	if update.Message != nil {
		return update.Message.From.ID
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.From.ID
	}
	return -1
}

func ParseCallback(update *tgbotapi.Update) string {
	if update.CallbackQuery == nil {
		return ""
	}

	return update.CallbackQuery.Data
}

func ParseMessageCommand(update *tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.Command()
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Command()
	}
	return ""
}

func ParseMessageString(update *tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.Text
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.Text
	}

	return ""
}

func ParseMessageInt(update *tgbotapi.Update) int {
	numStr := ParseMessageString(update)
	num, err := strconv.Atoi(numStr)
	if err != nil {
		return -1
	}

	return num
}

func ParseMessageFloat(update *tgbotapi.Update) float64 {
	numStr := ParseMessageString(update)
	num, err := strconv.ParseFloat(numStr, 64)
	if err != nil {
		return -1
	}

	return num
}

func IsSkip(update *tgbotapi.Update) bool {
	return ParseCallback(update) == states.CallbackProcessSkip
}

func IsCancel(update *tgbotapi.Update) bool {
	return ParseCallback(update) == states.CallbackProcessCancel
}

func IsAgree(update *tgbotapi.Update) bool {
	return ParseCallback(update) == states.CallbackYes
}
