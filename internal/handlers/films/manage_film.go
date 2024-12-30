package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleManageFilmCommand(app models.App, session *models.Session) {
	film := session.FilmDetailState.Film

	msg := messages.BuildFilmDetailMessage(session, &film)
	msg += translator.Translate(session.Lang, "choiceAction", nil, nil)

	keyboard := keyboards.BuildFilmManageKeyboard(session)

	app.SendImage(film.ImageURL, msg, keyboard)
}

func HandleManageFilmButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackManageFilmSelectBack:
		HandleFilmsDetailCommand(app, session)

	case states.CallbackManageFilmSelectUpdate:
		HandleUpdateFilmCommand(app, session)

	case states.CallbackManageFilmSelectDelete:
		HandleDeleteFilmCommand(app, session)
	}
}
