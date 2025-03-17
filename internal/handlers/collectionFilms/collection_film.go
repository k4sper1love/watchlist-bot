package collectionFilms

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// HandleCollectionFilmsButtons handles button interactions related to adding films to collections or vice versa.
// Supports actions like switching between adding a film to a collection or a collection to a film.
func HandleCollectionFilmsButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallCollectionFilmsFromFilm:
		clearAndHandleCollectionToFilm(app, session)

	case states.CallCollectionFilmsFromCollection:
		HandleOptionsFilmToCollectionCommand(app, session)
	}
}

// handleFilmsWithContext determines the current context (film or collection) and delegates to the appropriate handler.
// Used to manage navigation between film and collection views.
func handleFilmsWithContext(app models.App, session *models.Session) {
	if session.Context == states.CtxFilm {
		films.HandleFilmDetailCommand(app, session)
	} else if session.Context == states.CtxCollection {
		films.HandleFilmsCommand(app, session)
	}
}

// addFilmToCollection adds a film to a collection using the Watchlist service.
// Sends a success message upon completion and clears the session states.
func addFilmToCollection(app models.App, session *models.Session) {
	collectionFilm, err := watchlist.AddCollectionFilm(app, session)
	if err != nil {
		app.SendMessage(messages.CreateFilmFailure(session), keyboards.Back(session, states.CallMenuFilms))
		return
	}

	app.SendMessage(messages.AddFilmToCollectionSuccess(session, collectionFilm), nil)
	session.ClearAllStates()
	handleFilmsWithContext(app, session)
}

// clearAndHandleCollectionToFilm resets the session state and navigates to the "add collection to film" command.
func clearAndHandleCollectionToFilm(app models.App, session *models.Session) {
	session.CollectionFilmsState.CurrentPage = 1
	session.CollectionsState.Name = ""
	HandleAddCollectionToFilmCommand(app, session)
}

// clearAndHandleFilmToCollection resets the session state and navigates to the "add film to collection" command.
func clearAndHandleFilmToCollection(app models.App, session *models.Session) {
	session.CollectionFilmsState.CurrentPage = 1
	session.FilmsState.Title = ""
	HandleAddFilmToCollectionCommand(app, session)
}
