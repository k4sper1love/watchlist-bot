package collections

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleDeleteCollectionFilmCommand(app models.App, session *models.Session) {
	msg := fmt.Sprintf("Вы уверены, что хотите удалить фильм %q", session.CollectionFilmState.Object.Title)

	keyboard := builders.NewKeyboard(1).AddSurvey().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessDeleteCollectionFilmAwaitingConfirm)
}

func HandleDeleteCollectionFilmProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.ProcessDeleteCollectionFilmAwaitingConfirm:
		parseDeleteCollectionFilmConfirm(app, session)
	}
}

func parseDeleteCollectionFilmConfirm(app models.App, session *models.Session) {
	session.ClearState()

	switch utils.IsAgree(app.Upd) {
	case true:
		if err := watchlist.DeleteCollectionFilm(app, session); err != nil {
			app.SendMessage("Не удалось удалить фильм из коллекции", nil)
			HandleManageCollectionFilmCommand(app, session)
			break
		}
		app.SendMessage("Фильм удален из коллекции успешно", nil)
		HandleCollectionFilmsCommand(app, session)

	case false:
		app.SendMessage("Действие отменено", nil)
		HandleManageCollectionFilmCommand(app, session)
	}
}
