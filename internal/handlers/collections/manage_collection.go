package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleManageCollectionCommand(app models.App, session *models.Session) {
	msg := messages.BuildCollectionDetailMessage(&session.CollectionDetailState.Collection)
	msg += "⚡️ <b>Выберите действие:</b> из меню ниже."

	keyboard := keyboards.NewKeyboard().
		AddCollectionsUpdate().
		AddCollectionsDelete().
		AddBack(states.CallbackManageCollectionSelectBack).
		Build()

	app.SendMessage(msg, keyboard)
}

func HandleManageCollectionButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackManageCollectionSelectBack:
		session.SetContext(states.ContextCollection)
		films.HandleFilmsCommand(app, session)

	case states.CallbackManageCollectionSelectUpdate:
		HandleUpdateCollectionCommand(app, session)

	case states.CallbackManageCollectionSelectDelete:
		HandleDeleteCollectionCommand(app, session)
	}
}
