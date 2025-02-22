package collectionFilms

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleCollectionFilmsButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallbackCollectionFilmsFromFilm:
		session.CollectionFilmsState.CurrentPage = 1
		session.CollectionsState.Name = ""
		HandleAddCollectionToFilmCommand(app, session)

	case states.CallbackCollectionFilmsFromCollection:
		HandleOptionsFilmToCollectionCommand(app, session)
	}
}

func addFilmToCollection(app models.App, session *models.Session) {
	collectionFilm, err := watchlist.AddCollectionFilm(app, session)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "createFilmFailure", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackMenuSelectFilms).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		return
	}

	msg := "➕ " + translator.Translate(session.Lang, "filmToCollectionSuccess", map[string]interface{}{
		"Film":       collectionFilm.Film.Title,
		"Collection": collectionFilm.Collection.Name,
	}, nil)

	app.SendMessage(msg, nil)

	session.ClearAllStates()

	if session.Context == states.ContextFilm {
		films.HandleFilmsDetailCommand(app, session)
	} else if session.Context == states.ContextCollection {
		films.HandleFilmsCommand(app, session)
	}
}
