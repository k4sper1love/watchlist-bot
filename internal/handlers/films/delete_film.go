package films

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleDeleteFilmCommand(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "deleteFilmConfirm", map[string]interface{}{
		"Film": session.FilmDetailState.Title,
	}, nil)

	keyboard := keyboards.NewKeyboard().AddSurvey().Build(session.Lang)

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
			msg := translator.Translate(session.Lang, "deleteFilmFailure", map[string]interface{}{
				"Film": session.FilmDetailState.Title,
			}, nil)

			app.SendMessage(msg, nil)
			HandleManageFilmCommand(app, session)
			break
		}

		msg := translator.Translate(session.Lang, "deleteFilmSuccess", map[string]interface{}{
			"Film": session.FilmDetailState.Title,
		}, nil)

		app.SendMessage(msg, nil)
		HandleFilmsCommand(app, session)

	case false:
		msg := translator.Translate(session.Lang, "cancelAction", nil, nil)
		app.SendMessage(msg, nil)
		HandleManageFilmCommand(app, session)
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
