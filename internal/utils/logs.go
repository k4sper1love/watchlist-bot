package utils

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"log/slog"
)

func LogRequestError(message string, err error, method, requestURL string) {
	sl.Log.Error(
		message,
		slog.Any("error", err),
		slog.String("method", method),
		slog.String("url", requestURL),
	)
}

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

func LogParseJSONError(err error, method, requestURL string) {
	sl.Log.Error(
		"failed parsing",
		slog.Any("error", err),
		slog.String("method", method),
		slog.String("url", requestURL),
	)
}

func LogUpdateInfo(telegramID, messageID int, updateType string) {
	sl.Log.Info(
		"received update",
		slog.Int("from", telegramID),
		slog.Int("message_id", messageID),
		slog.String("type", updateType),
	)
}

func LogMessageInfo(chatID int64, messageID int, hasText, hasImage, isPinned bool) {
	sl.Log.Info(
		"bot sent message",
		slog.Int64("to", chatID),
		slog.Int("message_id", messageID),
		slog.Bool("has_text", hasText),
		slog.Bool("has_image", hasImage),
		slog.Bool("is_pinned", isPinned),
	)
}

func LogMessageError(err error, chatID int64, messageID int) {
	var args []slog.Attr

	if err != nil {
		args = append(args, slog.Any("error", err))
	}

	if chatID != -1 {
		args = append(args, slog.Int64("chat_id", chatID))
	}

	if messageID != -1 {
		args = append(args, slog.Int("message_id", messageID))
	}

	sl.Log.Error("error when sending message", args)
}

func LogRemoveFileWarn(err error, path string) {
	sl.Log.Warn(
		"failed to remove file",
		slog.Any("error", err),
		slog.String("path", path),
	)
}

func LogDownloadFileError(err error, url string) {
	sl.Log.Warn(
		"failed to download file",
		slog.Any("error", err),
		slog.String("url", url),
	)
}

func LogBodyCloseWarn(err error) {
	sl.Log.Warn(
		"failed to close body",
		slog.Any("error", err),
	)
}

func LogFileCloseWarn(err error) {
	sl.Log.Warn(
		"failed to close file",
		slog.Any("error", err),
	)
}
