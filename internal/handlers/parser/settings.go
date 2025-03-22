package parser

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/security"
)

// ParseSettingsFilmsPageSize processes the input for the films page size setting.
// Validates the page size range and retries if the input is invalid.
// Stores the validated page size in the session's FilmsState.
func ParseSettingsFilmsPageSize(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 1, 10, utils.ParseMessageInt, utils.IsValidNumberRange[int], validator.HandleInvalidInputRange[int], func(s *models.Session, v int) { session.FilmsState.PageSize = v })
}

// ParseSettingsCollectionsPageSize processes the input for the collections page size setting.
// Validates the page size range and retries if the input is invalid.
// Stores the validated page size in the session's CollectionsState.
func ParseSettingsCollectionsPageSize(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 1, 10, utils.ParseMessageInt, utils.IsValidNumberRange[int], validator.HandleInvalidInputRange[int], func(s *models.Session, v int) { session.CollectionsState.PageSize = v })
}

// ParseSettingsObjectsPageSize processes the input for the collection objects page size setting.
// Validates the page size range and retries if the input is invalid.
// Stores the validated page size in the session's CollectionFilmsState.
func ParseSettingsObjectsPageSize(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 1, 10, utils.ParseMessageInt, utils.IsValidNumberRange[int], validator.HandleInvalidInputRange[int], func(s *models.Session, v int) { session.CollectionFilmsState.PageSize = v })
}

// ParseKinopoiskToken processes the input for the Kinopoisk API token.
// Encrypts the token using the security package and stores it in the session.
func ParseKinopoiskToken(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	if encryptedToken, err := security.Encrypt(utils.ParseMessageString(app.Update)); err != nil {
		utils.LogEncryptError(session.TelegramID, err)
		app.SendMessage(messages.SomeError(session), nil)
	} else {
		session.KinopoiskAPIToken = encryptedToken
		app.SendMessage(messages.KinopoiskTokenSuccess(session), nil)
	}

	session.ClearState()
	next(app, session)
}
