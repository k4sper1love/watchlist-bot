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

// GetItemID calculates the global item ID based on the index, current page, and page size.
// It is used for pagination in lists of items (e.g., films or collections).
func GetItemID(index, currentPage, pageSize int) int {
	return (index + 1) + (currentPage-1)*pageSize
}

// IsBotMessage checks if the update originates from a bot.
// It returns true if the message or callback query is sent by a bot.
func IsBotMessage(update *tgbotapi.Update) bool {
	if update == nil {
		return false
	}
	return (update.Message != nil && update.Message.From.IsBot) ||
		(update.CallbackQuery != nil && update.CallbackQuery.From.IsBot)
}

// ParseMessageID extracts the message ID from the update.
// It handles both direct messages and callback queries with associated messages.
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

// ParseTelegramID extracts the Telegram user ID from the update.
// It handles both direct messages and callback queries.
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

// ParseTelegramName extracts the first name of the Telegram user from the update.
// If no name is available, it defaults to "Guest".
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

// ParseTelegramUsername extracts the username of the Telegram user from the update.
// If no username is available, it returns an empty string.
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

// ParseLanguageCode extracts the language code of the Telegram user from the update.
// If no language code is available, it defaults to "en" (English).
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

// ParseCallback extracts the callback data from the update.
// It returns an empty string if no callback query is present.
func ParseCallback(update *tgbotapi.Update) string {
	if update != nil && update.CallbackQuery != nil {
		return update.CallbackQuery.Data
	}
	return ""
}

// ParseMessageCommand extracts the command from the message in the update.
// It handles both direct messages and callback queries with associated messages.
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

// ParseMessageString extracts the text content of the message from the update.
// It handles both direct messages and callback queries with associated messages.
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

// ParseMessageInt parses the message text as an integer.
// If parsing fails, it returns -1.
func ParseMessageInt(update *tgbotapi.Update) int {
	num, err := strconv.Atoi(ParseMessageString(update))
	if err != nil {
		return -1
	}
	return num
}

// ParseMessageFloat parses the message text as a float64.
// If parsing fails, it returns -1.
func ParseMessageFloat(update *tgbotapi.Update) float64 {
	num, err := strconv.ParseFloat(ParseMessageString(update), 64)
	if err != nil {
		return -1
	}
	return num
}

// IsSkip checks if the callback query corresponds to a "skip" action.
func IsSkip(update *tgbotapi.Update) bool {
	return ParseCallback(update) == states.CallProcessSkip
}

// IsCancel checks if the callback query corresponds to a "cancel" action.
func IsCancel(update *tgbotapi.Update) bool {
	return ParseCallback(update) == states.CallProcessCancel
}

// IsReset checks if the callback query corresponds to a "reset" action.
func IsReset(update *tgbotapi.Update) bool {
	return ParseCallback(update) == states.CallProcessReset
}

// IsAgree checks if the callback query corresponds to an "agree" action.
func IsAgree(update *tgbotapi.Update) bool {
	return ParseCallback(update) == states.CallYes
}

// IsDecrease checks if the callback query corresponds to a "decrease" action.
func IsDecrease(update *tgbotapi.Update) bool {
	return ParseCallback(update) == states.CallDecrease
}

// ExtractKinopoiskQuery extracts the query key and ID from a Kinopoisk URL.
// It supports both path-based and query-parameter-based IDs.
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

// SplitTextByLength splits a string into two parts based on a maximum length.
// It attempts to split at spaces or newlines for better readability.
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

// LastIndexRune finds the last occurrence of a target rune in a slice of runes within a given range.
// It returns the index of the rune or -1 if the rune is not found.
func LastIndexRune(runes []rune, maxLength int, target rune) int {
	for i := maxLength - 1; i >= 0; i-- {
		if runes[i] == target {
			return i
		}
	}
	return -1
}

// ExtractYoutubeVideoID extracts the video ID from a YouTube URL.
// It supports both "youtu.be" short URLs and standard URLs with a "v" query parameter.
func ExtractYoutubeVideoID(rawUrl string) (string, error) {
	parsedURL, err := url.Parse(rawUrl)
	if err != nil {
		return "", err
	}

	// Handle short URLs like "https://youtu.be/<videoID>".
	if parsedURL.Host == "youtu.be" {
		return parsedURL.Path[1:], nil
	}

	// Handle standard URLs with a "v" query parameter.
	query := parsedURL.Query()
	videoID := query.Get("v")
	if videoID == "" {
		return "", fmt.Errorf("could not extract video ID")
	}

	return videoID, nil
}

// Round rounds a floating-point number to two decimal places.
func Round(v float64) float64 {
	rounded, _ := strconv.ParseFloat(fmt.Sprintf("%.2f", v), 64)
	return rounded
}

// FormatTextDate converts an ISO 8601 date string into a human-readable format ("DD.MM.YYYY HH:MM").
func FormatTextDate(date string) string {
	parsedDate, err := time.Parse(time.RFC3339, date)
	if err != nil {
		return ""
	}
	return parsedDate.Format("02.01.2006 15:04")
}

// ParseISO8601Duration parses an ISO 8601 duration string (e.g., "PT1H30M") into a formatted duration string (e.g., "01:30:00").
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

// ParseSupportedLanguages lists all supported languages by reading JSON files from a directory.
// Each file name (without the ".json" extension) represents a supported language code.
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

// BoolToEmoji converts a boolean value to a corresponding emoji ("✔️" for true, "✖️" for false).
func BoolToEmoji(value bool) string {
	if value {
		return "✔️"
	}
	return "✖️"
}

// ViewedToEmojiColored converts a boolean value to a colored emoji ("✅" for true, "👀" for false).
func ViewedToEmojiColored(value bool) string {
	if value {
		return "✅"
	}
	return "👀"
}

// BoolToEmojiColored converts a boolean value to a colored emoji ("✅" for true, "❌" for false).
func BoolToEmojiColored(value bool) string {
	if value {
		return "✅"
	}
	return "❌"
}

// BoolToStar converts a boolean value to a star emoji ("⭐" for true, "☆" for false).
func BoolToStar(value bool) string {
	if value {
		return "⭐"
	}
	return "☆"
}

// BoolToStarOrEmpty converts a boolean value to a star emoji or an empty string ("⭐" for true, "" for false).
func BoolToStarOrEmpty(value bool) string {
	if value {
		return "⭐"
	}
	return ""
}

// SortDirectionToEmoji converts a sort direction string to a corresponding emoji ("⬇️" for descending, "⬆️" for ascending).
func SortDirectionToEmoji(value string) string {
	if strings.HasPrefix(value, "-") {
		return "⬇️"
	}
	return "⬆️"
}

// NumberToEmoji converts a number into a sequence of emoji digits (e.g., "123" -> "1️⃣2️⃣3️⃣").
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

// GetLogFilePath constructs the log file path for a user based on their ID.
// It checks if the log file exists and returns an error if it does not.
func GetLogFilePath(userID int) (string, error) {
	logFile := filepath.Join("logs/users", fmt.Sprintf("user_%d.log", userID))
	if _, err := os.Stat(logFile); os.IsNotExist(err) {
		return "", err
	}
	return logFile, nil
}

// CloseBody ensures that an HTTP response body is closed safely.
// It logs a warning if an error occurs during closing.
func CloseBody(Body io.ReadCloser) {
	if Body == nil {
		return
	} else if err := Body.Close(); err != nil {
		LogBodyCloseWarn(err)
	}
}

// CloseFile ensures that a file is closed safely.
// It logs a warning if an error occurs during closing.
func CloseFile(file *os.File) {
	if file == nil {
		return
	} else if err := file.Close(); err != nil {
		LogFileCloseWarn(err)
	}
}

// RemoveFile removes a file at the specified path.
// It logs a warning if an error occurs during removal.
func RemoveFile(path string) {
	if path == "" {
		return
	} else if err := os.Remove(path); err != nil {
		LogRemoveFileWarn(err, path)
	}
}

// CalculateOffset calculates the offset for pagination based on the current page and page size.
func CalculateOffset(page, pageSize int) int {
	return (page - 1) * pageSize
}
