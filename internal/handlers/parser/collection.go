package parser

import (
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// ParseCollectionName processes the input for a collection name.
// Validates the name length and retries if the input is invalid.
// Stores the validated name in the session's CollectionDetailState.
func ParseCollectionName(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 3, 100, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.CollectionDetailState.Name = v })
}

// ParseCollectionDescription processes the input for a collection description.
// Validates the description length and retries if the input is invalid.
// Stores the validated description in the session's CollectionDetailState.
func ParseCollectionDescription(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 0, 500, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.CollectionDetailState.Description = v })
}

// ParseCollectionFindName processes the input for searching collections by name.
// Parses the input string, resets the current page, and clears the session state before proceeding.
func ParseCollectionFindName(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.CollectionsState.Name = utils.ParseMessageString(app.Update)
	session.CollectionsState.CurrentPage = 1

	session.ClearState()
	next(app, session)
}
