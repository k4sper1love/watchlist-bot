package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleDeleteCollectionCommand(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "deleteCollectionConfirm", map[string]interface{}{
		"Collection": session.CollectionDetailState.Collection.Name,
	}, nil)

	keyboard := keyboards.NewKeyboard().AddSurvey().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessDeleteCollectionAwaitingConfirm)
}

func HandleDeleteCollectionProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.ProcessDeleteCollectionAwaitingConfirm:
		parseDeleteCollectionConfirm(app, session)
	}
}

func parseDeleteCollectionConfirm(app models.App, session *models.Session) {
	session.ClearState()

	switch utils.IsAgree(app.Upd) {
	case true:
		if err := watchlist.DeleteCollection(app, session); err != nil {
			msg := translator.Translate(session.Lang, "deleteCollectionFailure", map[string]interface{}{
				"Collection": session.CollectionDetailState.Collection.Name,
			}, nil)

			app.SendMessage(msg, nil)
			HandleManageCollectionCommand(app, session)
			break
		}

		msg := translator.Translate(session.Lang, "deleteCollectionSuccess", map[string]interface{}{
			"Collection": session.CollectionDetailState.Collection.Name,
		}, nil)

		app.SendMessage(msg, nil)
		HandleCollectionsCommand(app, session)

	case false:
		msg := translator.Translate(session.Lang, "cancelAction", nil, nil)
		app.SendMessage(msg, nil)
		HandleManageCollectionCommand(app, session)
	}
}
