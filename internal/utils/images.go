package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"io"
	"net/http"
	"os"
)

var SupportedTypes = []string{
	"image/jpeg",
	"image/png",
	"image/gif",
	"application/octet-stream",
}

func processImageFromMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) ([]byte, error) {
	if update.Message == nil || update.Message.Photo == nil {
		return nil, fmt.Errorf("not found photo")
	}

	photo := (*update.Message.Photo)[len(*update.Message.Photo)-1]

	file, err := bot.GetFile(tgbotapi.FileConfig{FileID: photo.FileID})
	if err != nil {
		return nil, err
	}

	fileURL := fmt.Sprintf("https://api.telegram.org/file/bot%s/%s", bot.Token, file.FilePath)

	return processImageFromURL(fileURL)
}

func processImageFromURL(imageURL string) ([]byte, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return nil, fmt.Errorf("failed to fetch URL: %w", err)
	}
	defer resp.Body.Close()

	contentType := resp.Header.Get("Content-Type")
	if !isSupportedImageType(contentType) {
		return nil, fmt.Errorf("unsupported content type: %s", contentType)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return data, nil
}

func DownloadImage(imageURL string) (string, error) {
	resp, err := http.Get(imageURL)
	if err != nil {
		return "", fmt.Errorf("error uploading an image by URL:%v", err)
	}
	defer resp.Body.Close()

	file, err := os.CreateTemp("", "image_*.jpg")
	if err != nil {
		return "", fmt.Errorf("error creating a temporary file: %v", err)
	}
	defer file.Close()

	_, err = io.Copy(file, resp.Body)
	if err != nil {
		return "", fmt.Errorf("error copying image data to a file: %v", err)
	}

	return file.Name(), nil
}

func ParseImageFromMessage(bot *tgbotapi.BotAPI, update *tgbotapi.Update) ([]byte, error) {
	var data []byte
	var err error

	if update.Message.Photo != nil {
		data, err = processImageFromMessage(bot, update)
		if err != nil {
			return nil, err
		}
	} else {
		imageURL := ParseMessageString(update)

		data, err = processImageFromURL(imageURL)
		if err != nil {
			return nil, err
		}
	}

	return data, nil
}

func ParseImageFromURL(imageURL string) ([]byte, error) {
	return processImageFromURL(imageURL)
}

func isSupportedImageType(contentType string) bool {
	for _, t := range SupportedTypes {
		if contentType == t {
			return true
		}
	}
	return false
}
