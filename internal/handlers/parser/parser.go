package parser

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// ProcessInput processes user input with validation and retries if the input is invalid.
// It supports generic types for input parsing and validation, allowing flexible handling of various data types.
// Parameters:
// - app: The application context.
// - session: The user session to store parsed and validated data.
// - retry: A function to retry the current step if validation fails.
// - next: A function to proceed to the next step after successful processing.
// - min, max: The range or constraints for validation.
// - parser: A function to parse the raw input into the desired type.
// - validatorFunc: A function to validate the parsed input against the constraints.
// - errorHandler: A function to handle validation errors and notify the user.
// - setter: A function to store the validated input in the session state.
func ProcessInput[T any, N int | float64](
	app models.App, session *models.Session,
	retry, next func(models.App, *models.Session),
	min, max N,
	parser func(*tgbotapi.Update) T,
	validatorFunc func(T, N, N) bool,
	errorHandler func(models.App, *models.Session, N, N),
	setter func(*models.Session, T),
) {
	if utils.IsSkip(app.Update) {
		next(app, session)
		return
	}

	input := parser(app.Update)
	if !validatorFunc(input, min, max) {
		errorHandler(app, session, min, max)
		retry(app, session)
		return
	}

	setter(session, input)
	next(app, session)
}

// UploadImageFromMessage uploads an image from a Telegram message.
// Parses the image from the message and uploads it using the Watchlist service.
// Returns the uploaded image URL or an error if parsing or uploading fails.
func UploadImageFromMessage(app models.App) (string, error) {
	image, err := utils.ParseImageFromMessage(app.Bot, app.Update)
	if err != nil {
		return "", err
	}
	return watchlist.UploadImage(app, image)
}

// UploadImageFromURL uploads an image from a URL.
// Parses the image from the URL and uploads it using the Watchlist service.
// Returns the uploaded image URL or an error if parsing or uploading fails.
func UploadImageFromURL(app models.App, imageURL string) (string, error) {
	image, err := utils.ParseImageFromURL(imageURL)
	if err != nil {
		return "", err
	}
	return watchlist.UploadImage(app, image)
}
