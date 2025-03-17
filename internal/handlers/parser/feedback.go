package parser

import (
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// ParseFeedbackMessage processes the input for a feedback message.
// Validates the message length and retries if the input is invalid.
// Stores the validated message in the session's FeedbackState.
func ParseFeedbackMessage(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 0, 3000, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { session.FeedbackState.Message = v })
}
