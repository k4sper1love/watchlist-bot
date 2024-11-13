package films

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleDeleteFilmCommand(app models.App, session *models.Session) {
	msg := fmt.Sprintf("Вы уверены, что хотите удалить фильм %q", session.FilmDetailState.Object.Title)

	keyboard := builders.NewKeyboard(1).AddSurvey().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessDeleteFilmAwaitingConfirm)
}

func HandleDeleteFilmProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.ProcessDeleteFilmAwaitingConfirm:
		parseDeleteFilmConfirm(app, session)
	}
}

func parseDeleteFilmConfirm(app models.App, session *models.Session) {
	session.ClearState()

	switch utils.IsAgree(app.Upd) {
	case true:
		if err := watchlist.DeleteFilm(app, session); err != nil {
			app.SendMessage("Не удалось удалить фильм", nil)
			HandleManageFilmCommand(app, session)
			break
		}
		app.SendMessage("Фильм удален успешно", nil)
		HandleFilmsCommand(app, session)

	case false:
		app.SendMessage("Действие отменено", nil)
		HandleManageFilmCommand(app, session)
	}
}
