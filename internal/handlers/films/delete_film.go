package films

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleDeleteFilmCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildDeleteFilmMessage(session), keyboards.BuildKeyboardWithSurvey(session))
	session.SetState(states.ProcessDeleteFilmAwaitingConfirm)
}

func HandleDeleteFilmProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.ProcessDeleteFilmAwaitingConfirm:
		parseDeleteFilmConfirm(app, session)
		session.ClearState()
	}
}

func parseDeleteFilmConfirm(app models.App, session *models.Session) {
	if !utils.IsAgree(app.Update) {
		app.SendMessage(messages.BuildCancelActionMessage(session), nil)
		HandleManageFilmCommand(app, session)
		return
	}

	if err := DeleteFilm(app, session); err != nil {
		app.SendMessage(messages.BuildDeleteFilmFailureMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackFilmsManage))
		return
	}

	app.SendMessage(messages.BuildDeleteFilmSuccessMessage(session), nil)
	HandleFilmsCommand(app, session)
}

func DeleteFilm(app models.App, session *models.Session) error {
	switch session.Context {
	case states.ContextFilm:
		return watchlist.DeleteFilm(app, session)
	case states.ContextCollection:
		return watchlist.DeleteCollectionFilm(app, session)
	default:
		return fmt.Errorf("unsupported session context: %s", session.Context)
	}
}
