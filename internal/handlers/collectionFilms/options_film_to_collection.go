package collectionFilms

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleOptionsFilmToCollectionCommand(app models.App, session *models.Session) {
	msg := messages.BuildCollectionHeader(session)
	msg += "<b>" + translator.Translate(session.Lang, "choiceAction", nil, nil) + "</b>"

	keyboard := keyboards.BuildOptionsFilmToCollectionKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func HandleOptionsFilmToCollectionButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallbackOptionsFilmToCollectionBack:
		films.HandleFilmsCommand(app, session)

	case states.CallbackOptionsFilmToCollectionNew:
		films.HandleNewFilmCommand(app, session)

	case states.CallbackOptionsFilmToCollectionExisting:
		session.CollectionFilmsState.CurrentPage = 1
		session.FilmsState.Title = ""
		HandleAddFilmToCollectionCommand(app, session)
	}
}
