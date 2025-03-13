package parser

import (
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func ParseCollectionName(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 3, 100, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.CollectionDetailState.Name = v })
}

func ParseCollectionDescription(app models.App, session *models.Session, retry, next func(models.App, *models.Session)) {
	ProcessInput(app, session, retry, next, 0, 500, utils.ParseMessageString, utils.IsValidStringLength, validator.HandleInvalidInputLength, func(s *models.Session, v string) { s.CollectionDetailState.Description = v })
}

func ParseCollectionFindName(app models.App, session *models.Session, next func(models.App, *models.Session)) {
	session.CollectionsState.Name = utils.ParseMessageString(app.Update)
	session.CollectionsState.CurrentPage = 1

	session.ClearState()
	next(app, session)
}
