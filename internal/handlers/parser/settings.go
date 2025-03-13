package parser

import (
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
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
