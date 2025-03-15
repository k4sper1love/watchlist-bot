package parser

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/security"
)

func ParseSettingsFilmsPageSize(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 1, 10, utils.ParseMessageInt, utils.IsValidNumberRange[int], validator.HandleInvalidInputRange[int], func(s *models.Session, v int) { session.FilmsState.PageSize = v })
}

func ParseSettingsCollectionsPageSize(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 1, 10, utils.ParseMessageInt, utils.IsValidNumberRange[int], validator.HandleInvalidInputRange[int], func(s *models.Session, v int) { session.CollectionsState.PageSize = v })
}

func ParseSettingsObjectsPageSize(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 1, 10, utils.ParseMessageInt, utils.IsValidNumberRange[int], validator.HandleInvalidInputRange[int], func(s *models.Session, v int) { session.CollectionFilmsState.PageSize = v })
}

func ParseKinopoiskToken(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	if encryptedToken, err := security.Encrypt(utils.ParseMessageString(app.Update)); err != nil {
		utils.LogEncryptError(err)
		app.SendMessage(messages.SomeError(session), nil)
	} else {
		session.KinopoiskAPIToken = encryptedToken
		app.SendMessage(messages.KinopoiskTokenSuccess(session), nil)
	}

	session.ClearState()
	next(app, session)
}
