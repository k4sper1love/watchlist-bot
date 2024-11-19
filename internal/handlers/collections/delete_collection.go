package collections

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleDeleteCollectionCommand(app models.App, session *models.Session) {
	msg := fmt.Sprintf("Вы уверены, что хотите удалить коллекцию %q", session.CollectionDetailState.Collection.Name)

	keyboard := keyboards.NewKeyboard().AddSurvey().Build()

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
			app.SendMessage("Произошла ошибка при удалении", nil)
			HandleManageCollectionCommand(app, session)
			break
		}
		app.SendMessage("Успешно удалено", nil)
		HandleCollectionsCommand(app, session)

	case false:
		app.SendMessage("Действие отменено", nil)
		HandleManageCollectionCommand(app, session)
	}
}
