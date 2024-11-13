package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleManageFilmCommand(app models.App, session *models.Session) {
	film := session.FilmDetailState.Object

	msg := builders.BuildFilmDetailMessage(&film)
	msg += "Выберите действие"

	keyboard := builders.NewKeyboard(1).
		AddFilmsUpdate().
		AddFilmsDelete().
		AddBack(states.CallbackManageFilmSelectBack).
		Build()

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
