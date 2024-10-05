package handlers

import (
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/pkg/utils"
)

func handleCollectionsCommand(app config.App, session models.Session) {
	collectionsResponse, err := watchlist.GetCollections(app, session)
	if err != nil {
		utils.SendMessage(app.Bot, app.Upd, err.Error())
		return
	}

	if collectionsResponse.Metadata.TotalRecords == 0 {
		utils.SendMessage(app.Bot, app.Upd, "Не найдено коллекций")
	}
}
