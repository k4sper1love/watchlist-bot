package parser

import (
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func ParseProfileUsername(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 3, 20, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { session.ProfileState.Username = v })
}

func ParseProfileEmail(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 6, 254, utils.ParseMessageString, utils.IsValidEmail, validator.HandleInvalidInputEmail, func(s *models.Session, v string) { session.ProfileState.Email = v })
}
