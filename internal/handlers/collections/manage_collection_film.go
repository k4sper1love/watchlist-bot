package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func HandleManageCollectionFilmCommand(app models.App, session *models.Session) {
	film := session.CollectionFilmState.Object

	msg := builders.BuildCollectionFilmDetailMessage(&film)
	msg += "Выберите действие"

	keyboard := builders.NewKeyboard(1).
		AddCollectionFilmsManage().
		AddBack(states.CallbackManageCollectionFilmSelectBack).
		Build()

	app.SendImage(film.ImageURL, msg, keyboard)
}

func HandleManageCollectionFilmButtons(app models.App, session *models.Session) {
	switch session.State {
	case states.CallbackManageCollectionFilmSelectBack:
		HandleCollectionFilmsDetailCommand(app, session)

	case states.CallbackManageCollectionFilmSelectUpdate:
		HandleUpdateCollectionFilmCommand(app, session)

	case states.CallbackManageCollectionFilmSelectDelete:
		HandleDeleteCollectionFilmCommand(app, session)
	}
}
