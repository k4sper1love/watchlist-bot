package collectionFilms

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleOptionsFilmToCollectionCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.CollectionChoiceAction(session), keyboards.FilmToCollectionOptions(session))
}

func HandleOptionsFilmToCollectionButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallFilmToCollectionOptionBack:
		films.HandleFilmsCommand(app, session)

	case states.CallFilmToCollectionOptionNew:
		films.HandleNewFilmCommand(app, session)

	case states.CallFilmToCollectionOptionExisting:
		clearAndHandleFilmToCollection(app, session)
	}
}
