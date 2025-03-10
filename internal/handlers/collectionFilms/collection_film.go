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

func HandleCollectionFilmsButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallbackCollectionFilmsFromFilm:
		clearAndHandleCollectionToFilm(app, session)

	case states.CallbackCollectionFilmsFromCollection:
		HandleOptionsFilmToCollectionCommand(app, session)
	}
}

func handleFilmsWithContext(app models.App, session *models.Session) {
	if session.Context == states.ContextFilm {
		films.HandleFilmsDetailCommand(app, session)
	} else if session.Context == states.ContextCollection {
		films.HandleFilmsCommand(app, session)
	}
}

func addFilmToCollection(app models.App, session *models.Session) {
	collectionFilm, err := watchlist.AddCollectionFilm(app, session)
	if err != nil {
		app.SendMessage(messages.CreateFilmFailure(session), keyboards.Back(session, states.CallbackMenuSelectFilms))
		return
	}

	app.SendMessage(messages.AddFilmToCollectionSuccess(session, collectionFilm), nil)
	session.ClearAllStates()
	handleFilmsWithContext(app, session)
}

func clearAndHandleCollectionToFilm(app models.App, session *models.Session) {
	session.CollectionFilmsState.CurrentPage = 1
	session.CollectionsState.Name = ""
	HandleAddCollectionToFilmCommand(app, session)
}

func clearAndHandleFilmToCollection(app models.App, session *models.Session) {
	session.CollectionFilmsState.CurrentPage = 1
	session.FilmsState.Title = ""
	HandleAddFilmToCollectionCommand(app, session)
}
