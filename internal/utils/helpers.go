package utils

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

func GetItemID(index, currentPage, pageSize int) int {
	return (index + 1) + (currentPage-1)*pageSize
}

func IsBotMessage(update *tgbotapi.Update) bool {
	if update == nil {
		return false
	}
	return (update.Message != nil && update.Message.From.IsBot) ||
		(update.CallbackQuery != nil && update.CallbackQuery.From.IsBot)
}

func ParseMessageID(update *tgbotapi.Update) int {
	if update == nil {
		return -1
	}
	switch {
	case update.Message != nil:
		return update.Message.MessageID
	case update.CallbackQuery != nil && update.CallbackQuery.Message != nil:
		return update.CallbackQuery.Message.MessageID
	default:
		return -1
	}
}

func ParseTelegramID(update *tgbotapi.Update) int {
	if update == nil {
		return -1
	}
	switch {
	case update.Message != nil:
		return update.Message.From.ID
	case update.CallbackQuery != nil:
		return update.CallbackQuery.From.ID
	default:
		return -1
	}
}

func ParseTelegramName(update *tgbotapi.Update) string {
	if update == nil {
		return "Guest"
	}
	switch {
	case update.Message != nil:
		return update.Message.From.FirstName
	case update.CallbackQuery != nil:
		return update.CallbackQuery.From.FirstName
	default:
		return "Guest"
	}
}

func ParseTelegramUsername(update *tgbotapi.Update) string {
	if update == nil {
		return ""
	}
	switch {
	case update.Message != nil:
		return update.Message.From.UserName
	case update.CallbackQuery != nil:
		return update.CallbackQuery.From.UserName
	default:
		return ""
	}
}

func ParseLanguageCode(update *tgbotapi.Update) string {
	if update == nil {
		return "en"
	}
	switch {
	case update.Message != nil:
		return update.Message.From.LanguageCode
	case update.CallbackQuery != nil:
		return update.CallbackQuery.From.LanguageCode
	default:
		return "en"
	}
}

func ParseCallback(update *tgbotapi.Update) string {
	if update != nil && update.CallbackQuery != nil {
		return update.CallbackQuery.Data
	}
	return ""
}

func ParseMessageCommand(update *tgbotapi.Update) string {
	if update == nil {
		return ""
	}
	switch {
	case update.Message != nil:
		return update.Message.Command()
	case update.CallbackQuery != nil && update.CallbackQuery.Message != nil:
		return update.CallbackQuery.Message.Command()
	default:
		return ""
	}
}

func ParseMessageString(update *tgbotapi.Update) string {
	if update == nil {
		return ""
	}
	switch {
	case update.Message != nil:
		return update.Message.Text
	case update.CallbackQuery != nil && update.CallbackQuery.Message != nil:
		return update.CallbackQuery.Message.Text
	default:
		return ""
	}
}

func ParseMessageInt(update *tgbotapi.Update) int {
	num, err := strconv.Atoi(ParseMessageString(update))
	if err != nil {
		return -1
	}
	return num
}

func ParseMessageFloat(update *tgbotapi.Update) float64 {
	num, err := strconv.ParseFloat(ParseMessageString(update), 64)
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

	if !strings.Contains(parsedUrl.Host, "kinopoisk.ru") {
		return "", "", fmt.Errorf("URL is not from kinopoisk.ru")
	}

	parts := strings.Split(strings.Trim(parsedUrl.Path, "/"), "/")
	if len(parts) > 1 {
		return "id", parts[1], nil
	}

	query := parsedUrl.Query()
	if id, ok := query["rt"]; ok && len(id) > 0 {
		return "externalId.kpHD", id[0], nil
	}

	return "", "", fmt.Errorf("ID not found")
}

func SplitTextByLength(text string, maxLength int) (string, string) {
	runes := []rune(text)

	if len(runes) <= maxLength {
		return text, ""
	}

	splitPoint := maxLength
	if idx := LastIndexRune(runes, splitPoint, ' '); idx != -1 {
		splitPoint = idx
	}
	if idx := LastIndexRune(runes, splitPoint, '\n'); idx != -1 && idx > splitPoint {
		splitPoint = idx
	}

	firstPart := string(runes[:splitPoint])
	secondPart := string(runes[splitPoint:])
	if len(secondPart) > 0 {
		firstPart += "..."
	}

	return firstPart, secondPart
}

func LastIndexRune(runes []rune, maxLength int, target rune) int {
	for i := maxLength - 1; i >= 0; i-- {
		if runes[i] == target {
			return i
		}
	}
	return -1
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
		return "", fmt.Errorf("could not extract video ID")
	}

	return videoID, nil
}

func Round(v float64) float64 {
	rounded, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", v), 64)
	return rounded
}

func FormatTextDate(date string) string {
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return ""
	}
	return parsedDate.Format("02.01.2006 15:04")
}

func ParseISO8601Duration(isoDuration string) string {
	duration, err := time.ParseDuration(strings.ReplaceAll(strings.ToLower(isoDuration), "pt", ""))
	if err != nil {
		return ""
	}

	hours := int(duration.Hours())
	minutes := int(duration.Minutes()) % 60
	seconds := int(duration.Seconds()) % 60

	if hours > 0 {
		return fmt.Sprintf("%02d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%02d:%02d", minutes, seconds)
}

func ParseSupportedLanguages(dir string) ([]string, error) {
	files, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	var languages []string
	for _, file := range files {
		if !file.IsDir() && strings.HasSuffix(file.Name(), ".json") {
			languages = append(languages, strings.TrimSuffix(file.Name(), ".json"))
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

func ViewedToEmojiColored(value bool) string {
	if value {
		return "✅"
	}
	return "👀"
}

func BoolToEmojiColored(value bool) string {
	if value {
		return "✅"
	}
	return "❌"
}

func BoolToStar(value bool) string {
	if value {
		return "⭐"
	}
	return "☆"
}

func BoolToStarOrEmpty(value bool) string {
	if value {
		return "⭐"
	}
	return ""
}

func SortDirectionToEmoji(value string) string {
	if strings.HasPrefix(value, "-") {
		return "⬇️"
	}

	return "⬆️"
}

func NumberToEmoji(number int) string {
	emojis := []string{"0️⃣", "1️⃣", "2️⃣", "3️⃣", "4️⃣", "5️⃣", "6️⃣", "7️⃣", "8️⃣", "9️⃣"}
	if number < 10 {
		return emojis[number]
	}

	var result string
	for number > 0 {
		digit := number % 10
		result = emojis[digit] + result
		number /= 10
	}
	return result
}

func GetLogFilePath(userID int) (string, error) {
	logFile := filepath.Join("logs", fmt.Sprintf("user_%d.log", userID))
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return "", err
	}
	return logFile, nil
}

func CloseBody(Body io.ReadCloser) {
	if Body == nil {
		return
	} else if err := Body.Close(); err != nil {
		LogBodyCloseWarn(err)
	}
}

func CloseFile(file *os.File) {
	if file == nil {
		return
	} else if err := file.Close(); err != nil {
		LogFileCloseWarn(err)
	}
}

func RemoveFile(path string) {
	if path == "" {
		return
	} else if err := os.Remove(path); err != nil {
		LogRemoveFileWarn(err, path)
	}
}

func CalculateOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}
