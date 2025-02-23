package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"io"
	"log/slog"
	"net/http"
	"os"
)

var supportedTypes = map[string]bool{
	"image/jpeg":               true,
	"image/png":                true,
	"image/gif":                true,
	"application/octet-stream": true,
}

func ParseImageFromMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) ([]byte, error) {
	if update.Message == nil || update.Message.Photo == nil {
		return ParseImageFromURL(ParseMessageString(update))
	}

	photo := (*update.Message.Photo)[len(*update.Message.Photo)-1]

	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: photo.FileID})
	if err != nil {
		sl.Log.Error("failed to get file", slog.Any("error", err))
		return nil, err
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
	return ParseImageFromURL(fileURL)
}

func ParseImageFromURL(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		sl.Log.Error("failed to get image by URL", slog.Any("error", err), slog.String("url", imageURL))
		return nil, err
	}
	defer CloseBody(resp.Body)

	if !isSupportedImageType(resp.Header.Get("Content-Type")) {
		sl.Log.Warn(
			"image has unsupported content type",
			slog.String("content-type", resp.Header.Get("Content-Type")),
			slog.String("url", imageURL),
		)
		return nil, fmt.Errorf("unsupported content type: %s", resp.Header.Get("Content-Type"))
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		sl.Log.Error("failed to read response body", slog.Any("error", err), slog.String("url", imageURL))
		return nil, err
	}

	return data, nil
}

func DownloadImage(imageURL string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		sl.Log.Error("failed to get image by URL", slog.Any("error", err), slog.String("url", imageURL))
		return "", err
	}
	defer CloseBody(resp.Body)

	file, err := os.CreateTemp("", "image_*.jpg")
	if err != nil {
		sl.Log.Error("failed to create temporary file", slog.Any("error", err), slog.String("url", imageURL))
		return "", err
	}
	defer CloseFile(file)

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		sl.Log.Error("failed to copy image data to file", slog.Any("error", err), slog.String("url", imageURL))
		return "", err
	}

	return file.Name(), nil
}

func isSupportedImageType(contentType string) bool {
	return supportedTypes[contentType]
}
