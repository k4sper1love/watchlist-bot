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
		keyboards.BuildFilmManageKeyboard(session),
	)
}

func HandleManageFilmButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallbackManageFilmSelectBack:
		HandleFilmsDetailCommand(app, session)

	case states.CallbackManageFilmSelectUpdate:
		HandleUpdateFilmCommand(app, session)

	case states.CallbackManageFilmSelectDelete:
		HandleDeleteFilmCommand(app, session)

	case states.CallbackManageFilmSelectRemoveFromCollection:
		handleRemoveFilmFromCollection(app, session)
	}
}

func handleRemoveFilmFromCollection(app models.App, session *models.Session) {
	if err := watchlist.DeleteCollectionFilm(app, session); err != nil {
		app.SendMessage(messages.RemoveFilmFailure(session), keyboards.BuildKeyboardWithBack(session, states.CallbackFilmsManage))
		return
	}

	app.SendMessage(messages.RemoveFilmSuccess(session), nil)
	HandleFilmsCommand(app, session)
}
