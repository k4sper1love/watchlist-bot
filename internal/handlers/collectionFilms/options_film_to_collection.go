package collectionFilms

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

// HandleOptionsFilmToCollectionCommand handles the command for choosing an action related to adding a film to a collection.
// Sends a message with options for creating a new film or selecting an existing one.
func HandleOptionsFilmToCollectionCommand(app models.App, session *models.Session) {
	app.SendMessage(messages.CollectionChoiceAction(session), keyboards.FilmToCollectionOptions(session))
}

// HandleOptionsFilmToCollectionButtons handles button interactions related to adding a film to a collection.
// Supports actions like going back, creating a new film, or selecting an existing film.
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
