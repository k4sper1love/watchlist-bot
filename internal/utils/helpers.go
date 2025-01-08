package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"
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

func ParseTelegramName(update *tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.From.FirstName
	} else if update.CallbackQuery.Message != nil {
		return update.CallbackQuery.Message.From.FirstName
	}

	return "Guest"
}

func ParseTelegramUsername(update *tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.From.UserName
	} else if update.CallbackQuery.Message != nil {
		return update.CallbackQuery.Message.From.UserName
	}

	return ""
}

func ParseLanguageCode(update *tgbotapi.Update) string {
	if update.Message != nil {
		return update.Message.From.LanguageCode
	} else if update.CallbackQuery != nil {
		return update.CallbackQuery.Message.From.LanguageCode
	}
	return "en"
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

func IsReset(update *tgbotapi.Update) bool {
	return ParseCallback(update) == states.CallbackProcessReset
}

func IsAgree(update *tgbotapi.Update) bool {
	return ParseCallback(update) == states.CallbackYes
}

func ExtractKinopoiskQuery(rawUrl string) (string, string, error) {
	parsedUrl, err := url.Parse(rawUrl)
	if err != nil {
		return "", "", err
	}

	host := parsedUrl.Host
	if !strings.Contains(host, "kinopoisk.ru") {
		return "", "", fmt.Errorf("URL is not from kinopoisk.ru")
	}

	if strings.HasPrefix(parsedUrl.Path, "/film/") || strings.HasPrefix(parsedUrl.Path, "/series/") {
		parts := strings.Split(strings.Trim(parsedUrl.Path, "/"), "/")
		if len(parts) > 1 {
			return "id", parts[1], nil
		}
		return "", "", fmt.Errorf("ID not found in URL path")
	}

	if strings.Contains(host, "hd.kinopoisk.ru") {
		query := parsedUrl.Query()
		if id, ok := query["rt"]; ok && len(id) > 0 {
			return "externalId.kpHD", id[0], nil
		}
		return "", "", fmt.Errorf("ID not found in URL query")
	}

	return "", "", fmt.Errorf("unsupported URL format")
}

func SplitTextByLength(text string, maxLength int) (string, string) {
	splitPoint := maxLength

	if idx := strings.LastIndex(text[:splitPoint], " "); idx != 1 {
		splitPoint = idx
	}

	firstPart := text[:splitPoint] + "..."
	secondPart := text[splitPoint:]

	return firstPart, secondPart
}

func CalculateNewElementPageAndIndex(totalRecords, pageSize int) (int, int) {
	newTotalRecords := totalRecords + 1

	newPage := (newTotalRecords + pageSize - 1) / pageSize

	newIndex := (newTotalRecords - 1) % pageSize

	return newPage, newIndex
}

func ExtractYoutubeVideoID(rawUrl string) (string, error) {
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	if parsedURL.Host == "youtu.be" {
		return parsedURL.Path[1:], nil
	}

	query := parsedURL.Query()

	videoID := query.Get("v")

	if videoID == "" {
		return "", fmt.Errorf("couldn't extract video ID")
	}

	return videoID, nil
}

func Round(v float64) (float64, error) {
	return strconv.ParseFloat(fmt.Sprintf("%.2f", v), 64)
}

func FormatTextDate(date string) string {
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return ""
	}
	return parsedDate.Format("02.01.2006 15:04")
}

func ParseISO8601Duration(isoDuration string) (string, error) {
	duration, err := time.ParseDuration(strings.ReplaceAll(strings.ToLower(isoDuration), "pt", ""))
	if err != nil {
		return "", nil
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds), nil
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds), nil
}

func ParseSupportedLanguages(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var languages []string

	for _, file := range files {
		if !file.IsDir() {
			fileName := file.Name()

			if strings.HasSuffix(fileName, ".json") {
				lang := strings.TrimSuffix(fileName, ".json")
				languages = append(languages, lang)
			}
		}
	}

	return languages, nil
}

func BoolToEmoji(value bool) string {
	if value {
		return "✔️"
	}
	return "✖️"
}

func SortDirectionToEmoji(value string) string {
	if strings.HasPrefix(value, "-") {
		return "⬇️"
	}

	return "⬆️"
}
