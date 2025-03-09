package parser

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func ParseBroadcastMessage(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 0, 3000, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.AdminState.Message = v })
}

func ParseBroadcastImage(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	if utils.IsSkip(app.Update) {
		next(app, session)
		return
	}

	imageURL, err := ParseAndUploadImageFromMessage(app)
	if err != nil {
		app.SendMessage(messages.ImageFailure(session), nil)
	}

	session.AdminState.ImageURL = imageURL
	next(app, session)
}

func ParseBroadcastPin(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.AdminState.NeedFeedbackPin = utils.IsAgree(app.Update)
	next(app, session)
}
