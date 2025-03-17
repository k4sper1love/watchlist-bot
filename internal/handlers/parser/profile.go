package parser

import (
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// ParseProfileUsername processes the input for a user's profile username.
// Validates the username length and retries if the input is invalid.
// Stores the validated username in the session's ProfileState.
func ParseProfileUsername(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 3, 20, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { session.ProfileState.Username = v })
}

// ParseProfileEmail processes the input for a user's profile email.
// Validates the email format and retries if the input is invalid.
// Stores the validated email in the session's ProfileState.
func ParseProfileEmail(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 6, 254, utils.ParseMessageString, utils.IsValidEmail, validator.HandleInvalidInputEmail, func(s *models.Session, v string) { session.ProfileState.Email = v })
}
