package utils

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"log/slog"
)

// LogRequestError logs an error that occurs during an HTTP request.
// It includes details such as the error message, HTTP method, and request URL.
func LogRequestError(message string, err error, method, requestURL string) {
	sl.Log.Error(
		message,
		slog.Any("error", err),
		slog.String("method", method),
		slog.String("url", requestURL))
}

// LogResponseError logs an error related to an HTTP response.
// It includes details such as the response status code, HTTP method, and request URL.
// Returns an error with a formatted message for further handling.
func LogResponseError(url, method string, code int, status string) error {
	sl.Log.Error(
		"failed response",
		slog.String("method", method),
		slog.String("url", url),
		slog.Int("code", code),
		slog.String("status", status),
	)

	return fmt.Errorf("failed response from %s with code %d", url, code)
}

// LogParseJSONError logs an error that occurs while parsing JSON data.
// It includes details such as the error, HTTP method, and request URL.
func LogParseJSONError(err error, method, requestURL string) {
	sl.Log.Error(
		"failed parsing",
		slog.Any("error", err),
		slog.String("method", method),
		slog.String("url", requestURL))
}

// LogUpdateInfo logs information about a received update from Telegram.
// It includes details such as the Telegram user ID, message ID, and update type.
func LogUpdateInfo(telegramID, messageID int, updateType string) {
	sl.Log.Info(
		"received update",
		slog.Int("from", telegramID),
		slog.Int("message_id", messageID),
		slog.String("type", updateType))
}

// LogMessageInfo logs information about a message sent by the bot.
// It includes details such as the chat ID, message ID, and whether the message contains text, an image, or is pinned.
func LogMessageInfo(chatID int64, messageID int, hasText, hasImage, isPinned bool) {
	sl.Log.Info(
		"bot sent message",
		slog.Int64("to", chatID),
		slog.Int("message_id", messageID),
		slog.Bool("has_text", hasText),
		slog.Bool("has_image", hasImage),
		slog.Bool("is_pinned", isPinned))
}

// LogMessageError logs an error that occurs while sending a message.
// It includes optional details such as the error, chat ID, and message ID.
func LogMessageError(err error, chatID int64, messageID int) {
	sl.Log.Error("error when sending message",
		slog.Any("error", err),
		slog.Int64("chat_id", chatID),
		slog.Int("message_id", messageID),
	)
}

// LogRemoveFileWarn logs a warning when a file cannot be removed.
// It includes details such as the error and file path.
func LogRemoveFileWarn(err error, path string) {
	sl.Log.Warn(
		"failed to remove file",
		slog.Any("error", err),
		slog.String("path", path))
}

// LogBodyCloseWarn logs a warning when an HTTP response body cannot be closed.
func LogBodyCloseWarn(err error) {
	sl.Log.Warn(
		"failed to close body",
		slog.Any("error", err))
}

// LogFileCloseWarn logs a warning when a file cannot be closed.
func LogFileCloseWarn(err error) {
	sl.Log.Warn(
		"failed to close file",
		slog.Any("error", err))
}

// LogParseSelectError logs an error that occurs while parsing a select value from a callback.
// It includes details such as the error and callback data.
func LogParseSelectError(err error, callback string) {
	sl.Log.Error(
		"failed to parse select value",
		slog.Any("error", err),
		slog.String("callback", callback))
}

// LogEncryptError logs an error that occurs during data encryption.
func LogEncryptError(err error) {
	sl.Log.Error(
		"failed to encrypt data",
		slog.Any("error", err))
}

// LogDecryptError logs an error that occurs during data decryption.
func LogDecryptError(err error) {
	sl.Log.Error(
		"failed to decrypt data",
		slog.Any("error", err))
}
