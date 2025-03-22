package utils

import (
	"fmt"
	"log/slog"
)

// LogRequestError logs an error that occurs during an HTTP request.
func LogRequestError(telegramID int, message string, err error, method, requestURL string) {
	slog.Error(
		message,
		slog.Any("error", err),
		slog.String("method", method),
		slog.String("url", requestURL),
		slog.Int("telegram_id", telegramID))
}

// LogRequestDebug logs a debug-level message for an outgoing HTTP request.
func LogRequestDebug(telegramID int, method, requestURL string) {
	slog.Debug(
		"send request",
		slog.String("method", method),
		slog.String("url", requestURL),
		slog.Int("telegram_id", telegramID))
}

// LogResponseError logs an error related to an HTTP response.
// Returns an error with a formatted message for further handling.
func LogResponseError(telegramID int, url, method string, expectedCode, code int, status string) error {
	slog.Error(
		"failed response",
		slog.String("method", method),
		slog.String("url", url),
		slog.Int("expected_code", expectedCode),
		slog.Int("code", code),
		slog.String("status", status),
		slog.Int("telegram_id", telegramID),
	)
	return fmt.Errorf("failed response from %s with code %d", url, code)
}

// LogParseJSONError logs an error that occurs while parsing JSON data.
func LogParseJSONError(telegramID int, err error, method, requestURL string) {
	slog.Error(
		"failed parsing",
		slog.Any("error", err),
		slog.String("method", method),
		slog.String("url", requestURL),
		slog.Int("telegram_id", telegramID))
}

// LogUpdateInfo logs information about a received update from Telegram.
func LogUpdateInfo(telegramID, messageID int, updateType string) {
	slog.Info(
		"received update",
		slog.Int("message_id", messageID),
		slog.String("type", updateType),
		slog.Int("telegram_id", telegramID))
}

// LogMessageInfo logs information about a message sent by the bot.
func LogMessageInfo(chatID int64, messageID int, hasText, hasImage, isPinned bool) {
	slog.Info(
		"bot sent message",
		slog.Int("message_id", messageID),
		slog.Bool("has_text", hasText),
		slog.Bool("has_image", hasImage),
		slog.Bool("is_pinned", isPinned),
		slog.Int64("telegram_id", chatID))
}

// LogMessageError logs an error that occurs while sending a message.
func LogMessageError(err error, chatID int64, messageID int) {
	slog.Error("error when sending message",
		slog.Any("error", err),
		slog.Int("message_id", messageID),
		slog.Int64("telegram_id", chatID))
}

// LogRemoveFileWarn logs a warning when a file cannot be removed.
func LogRemoveFileWarn(err error, path string) {
	slog.Warn(
		"failed to remove file",
		slog.Any("error", err),
		slog.String("path", path))
}

// LogBodyCloseWarn logs a warning when an HTTP response body cannot be closed.
func LogBodyCloseWarn(err error) {
	slog.Warn(
		"failed to close body",
		slog.Any("error", err))
}

// LogFileCloseWarn logs a warning when a file cannot be closed.
func LogFileCloseWarn(err error, filename string) {
	slog.Warn(
		"failed to close file",
		slog.Any("error", err),
		slog.String("filename", filename))
}

// LogParseSelectError logs an error that occurs while parsing a select value from a callback.
func LogParseSelectError(telegramID int, err error, callback string) {
	slog.Error(
		"failed to parse select value",
		slog.Any("error", err),
		slog.String("callback", callback),
		slog.Int("telegram_id", telegramID))
}

// LogEncryptError logs an error that occurs during data encryption.
func LogEncryptError(telegramID int, err error) {
	slog.Error(
		"failed to encrypt data",
		slog.Any("error", err),
		slog.Int("telegram_id", telegramID))
}

// LogDecryptError logs an error that occurs during data decryption.
func LogDecryptError(telegramID int, err error) {
	slog.Error(
		"failed to decrypt data",
		slog.Any("error", err),
		slog.Int("telegram_id", telegramID))
}

// LogParseLanguagesError logs an error that occurs while parsing supported languages from a directory.
func LogParseLanguagesError(err error, dir string) {
	slog.Error(
		"failed to parse supported languages",
		slog.Any("error", err),
		slog.String("dir", dir))
}

// LogParseFromURLError logs an error that occurs while parsing data from a URL.
func LogParseFromURLError(telegramID int, message string, err error, url string) {
	slog.Error(
		message,
		slog.Any("error", err),
		slog.String("url", url),
		slog.Int("telegram_id", telegramID))
}
