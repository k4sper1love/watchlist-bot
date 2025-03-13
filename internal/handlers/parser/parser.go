package parser

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

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

func UploadImageFromMessage(app models.App) (string, error) {
	image, err := utils.ParseImageFromMessage(app.Bot, app.Update)
	if err != nil {
		return "", err
	}
	return watchlist.UploadImage(app, image)
}

func UploadImageFromURL(app models.App, imageURL string) (string, error) {
	image, err := utils.ParseImageFromURL(imageURL)
	if err != nil {
		return "", err
	}
	return watchlist.UploadImage(app, image)
}
