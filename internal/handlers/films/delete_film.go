package films

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleDeleteFilmCommand(app models.App, session *models.Session) {
	msg := fmt.Sprintf("Вы уверены, что хотите удалить фильм %q", session.FilmDetailState.Film.Title)

	keyboard := keyboards.NewKeyboard().AddSurvey().Build()

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
	session.ClearAllStates()

	switch utils.IsAgree(app.Upd) {
	case true:
		if err := DeleteFilm(app, session); err != nil {
			app.SendMessage("Не удалось удалить фильм", nil)
			HandleManageFilmCommand(app, session)
			break
		}
		app.SendMessage("Фильм удален успешно", nil)
		HandleFilmsDetailCommand(app, session)

	case false:
		app.SendMessage("Действие отменено", nil)
		HandleFilmsDetailCommand(app, session)
	}
}

func DeleteFilm(app models.App, session *models.Session) error {
	switch session.Context {
	case states.ContextFilm:
		return watchlist.DeleteFilm(app, session)

	case states.ContextCollection:
		return watchlist.DeleteCollectionFilm(app, session)

	default:
		return fmt.Errorf("unsupported session context: %v", session.Context)
	}
}
