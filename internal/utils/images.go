package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"log/slog"
	"net/http"
	"os"
)

// supportedTypes defines the list of supported image MIME types.
var supportedTypes = map[string]bool{
	"image/jpeg":                true,
	"image/png":                 true,
	"image/gif":                 true,
	"application/octet-stream":  true,
	"image/png; charset=UTF-8":  true,
	"image/jpeg; charset=UTF-8": true,
	"image/gif; charset=UTF-8":  true,
}

// ParseImageFromMessage extracts an image from a Telegram message.
// It handles both direct photo messages and URLs provided in the message text.
func ParseImageFromMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) ([]byte, error) {
	// Check if the message contains a photo.
	if update.Message == nil || update.Message.Photo == nil {
		return ParseImageFromURL(ParseMessageString(update))
	}

	// Get the largest available photo size.
	photo := (*update.Message.Photo)[len(*update.Message.Photo)-1]

	// Fetch the file metadata from Telegram's API.
	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: photo.FileID})
	if err != nil {
		slog.Error("failed to get file", slog.Any("error", err))
		return nil, err
	}

	// Construct the full URL for the file and parse the image.
	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)
	return ParseImageFromURL(fileURL)
}

// ParseImageFromURL fetches an image from a given URL.
// It checks the content type to ensure it is supported and reads the image data.
func ParseImageFromURL(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		slog.Error("failed to get image by URL", slog.Any("error", err), slog.String("url", imageURL))
		return nil, err
	}
	defer CloseBody(resp.Body) // Ensure the response body is closed after use.

	// Validate the content type of the response.
	if !isSupportedImageType(resp.Header.Get("Content-Type")) {
		slog.Warn(
			"image has unsupported content type",
			slog.String("content-type", resp.Header.Get("Content-Type")),
			slog.String("url", imageURL),
		)
		return nil, fmt.Errorf("unsupported content type: %s", resp.Header.Get("Content-Type"))
	}

	// Read the image data from the response body.
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		slog.Error("failed to read response body", slog.Any("error", err), slog.String("url", imageURL))
		return nil, err
	}

	return data, nil
}

// DownloadImage downloads an image from a URL and saves it as a temporary file.
// It returns the path to the temporary file.
func DownloadImage(imageURL string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		slog.Error("failed to get image by URL", slog.Any("error", err), slog.String("url", imageURL))
		return "", err
	}
	defer CloseBody(resp.Body) // Ensure the response body is closed after use.

	// Create a temporary file to store the image.
	file, err := os.CreateTemp("", "image_*.jpg")
	if err != nil {
		slog.Error("failed to create temporary file", slog.Any("error", err), slog.String("url", imageURL))
		return "", err
	}
	defer CloseFile(file) // Ensure the file is closed after use.

	// Copy the image data into the temporary file.
	_, err = io.Copy(file, resp.Body)
	if err != nil {
		slog.Error("failed to copy image data to file", slog.Any("error", err), slog.String("url", imageURL))
		return "", err
	}

	return file.Name(), nil
}

// isSupportedImageType checks if the given content type is supported.
func isSupportedImageType(contentType string) bool {
	return supportedTypes[contentType]
}
