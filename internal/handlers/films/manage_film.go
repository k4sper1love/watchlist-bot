package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// HandleManageFilmCommand handles the command for managing a specific film.
// Sends a message with options to update, delete, or remove the film from a collection.
func HandleManageFilmCommand(app models.App, session *models.Session) {
	app.SendImage(
		session.FilmDetailState.Film.ImageURL,
		messages.ManageFilm(session),
		keyboards.FilmManage(session),
	)
}

// HandleManageFilmButtons handles button interactions related to managing a film.
// Supports actions like going back, updating, deleting, or removing the film from a collection.
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

// handleRemoveFilmFromCollection removes the current film from its associated collection.
func handleRemoveFilmFromCollection(app models.App, session *models.Session) {
	if err := watchlist.DeleteCollectionFilm(app, session); err != nil {
		app.SendMessage(messages.RemoveFilmFailure(session), keyboards.Back(session, states.CallFilmsManage))
		return
	}

	app.SendMessage(messages.RemoveFilmSuccess(session), nil)
	HandleFilmsCommand(app, session)
}
