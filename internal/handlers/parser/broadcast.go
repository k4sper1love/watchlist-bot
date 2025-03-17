package parser

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// ParseBroadcastMessage processes the input for a broadcast message.
// Validates the message length and retries if the input is invalid.
// Stores the validated message in the session's AdminState.
func ParseBroadcastMessage(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 0, 3000, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.AdminState.Message = v })
}

// ParseBroadcastImage processes the input for a broadcast image.
// If the user skips the image, proceeds to the next step; otherwise, uploads the image and stores its URL.
func ParseBroadcastImage(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	if utils.IsSkip(app.Update) {
		next(app, session)
		return
	}

	imageURL, err := UploadImageFromMessage(app)
	if err != nil {
		app.SendMessage(messages.ImageFailure(session), nil)
	}

	session.AdminState.ImageURL = imageURL
	next(app, session)
}

// ParseBroadcastPin processes the input for pinning a broadcast message.
// Sets the NeedPin flag in the session's AdminState based on the user's choice.
func ParseBroadcastPin(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.AdminState.NeedPin = utils.IsAgree(app.Update)
	next(app, session)
}
