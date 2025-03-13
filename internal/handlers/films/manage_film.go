package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleManageFilmCommand(app models.App, session *models.Session) {
	app.SendImage(
		session.FilmDetailState.Film.ImageURL,
		messages.ManageFilm(session),
		keyboards.FilmManage(session),
	)
}

func HandleManageFilmButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallManageFilmBack:
		HandleFilmDetailCommand(app, session)

	case states.CallManageFilmUpdate:
		HandleUpdateFilmCommand(app, session)

	case states.CallManageFilmDelete:
		HandleDeleteFilmCommand(app, session)

	case states.CallManageFilmRemoveFromCollection:
		handleRemoveFilmFromCollection(app, session)
	}
}

func handleRemoveFilmFromCollection(app models.App, session *models.Session) {
	if err := watchlist.DeleteCollectionFilm(app, session); err != nil {
		app.SendMessage(messages.RemoveFilmFailure(session), keyboards.Back(session, states.CallFilmsManage))
		return
	}

	app.SendMessage(messages.RemoveFilmSuccess(session), nil)
	HandleFilmsCommand(app, session)
}
