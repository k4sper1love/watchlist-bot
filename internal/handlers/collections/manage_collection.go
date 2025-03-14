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
	app.SendMessage(messages.CollectionChoiceAction(session), keyboards.CollectionManage(session))
}

func HandleManageCollectionButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallManageCollectionBack:
		session.SetContext(states.CtxCollection)
		films.HandleFilmsCommand(app, session)

	case states.CallManageCollectionUpdate:
		HandleUpdateCollectionCommand(app, session)

	case states.CallManageCollectionDelete:
		HandleDeleteCollectionCommand(app, session)
	}
}
