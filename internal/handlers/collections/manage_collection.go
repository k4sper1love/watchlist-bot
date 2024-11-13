package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleManageCollectionCommand(app models.App, session *models.Session) {
	msg := builders.BuildCollectionDetailMessage(-1, &session.CollectionDetailState.Object.Collection)
	msg += "\nВыберите действие"

	keyboard := builders.NewKeyboard(1).
		AddCollectionsUpdate().
		AddCollectionsDelete().
		AddBack(states.CallbackManageCollectionSelectBack).
		Build()

	app.SendMessage(msg, keyboard)
}

func HandleManageCollectionButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackManageCollectionSelectBack:
		HandleCollectionFilmsCommand(app, session)

	case states.CallbackManageCollectionSelectUpdate:
		HandleUpdateCollectionCommand(app, session)

	case states.CallbackManageCollectionSelectDelete:
		HandleDeleteCollectionCommand(app, session)
	}
}
