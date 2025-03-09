package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleDeleteCollectionCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.DeleteCollection(session), keyboards.BuildKeyboardWithSurvey(session))
	session.SetState(states.ProcessDeleteCollectionAwaitingConfirm)
}

func HandleDeleteCollectionProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.ProcessDeleteCollectionAwaitingConfirm:
		parseDeleteCollectionConfirm(app, session)
		session.ClearState()
	}
}

func parseDeleteCollectionConfirm(app models.App, session *models.Session) {
	if !utils.IsAgree(app.Update) {
		app.SendMessage(messages.CancelAction(session), nil)
		HandleManageCollectionCommand(app, session)
		return
	}

	if err := watchlist.DeleteCollection(app, session); err != nil {
		app.SendMessage(messages.DeleteCollectionFailure(session), keyboards.BuildKeyboardWithBack(session, states.CallbackCollectionsManage))
		return
	}

	app.SendMessage(messages.DeleteCollectionSuccess(session), nil)
	HandleCollectionsCommand(app, session)
}
