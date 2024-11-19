package collectionFilms

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleCollectionFilmsButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackCollectionFilmsFromFilm:
		session.CollectionFilmsState.CurrentPage = 1
		HandleAddCollectionToFilmCommand(app, session)

	case states.CallbackCollectionFilmsFromCollection:
		HandleOptionsFilmToCollectionCommand(app, session)
	}
}

func addFilmToCollection(app models.App, session *models.Session) {
	collectionFilm, err := watchlist.AddCollectionFilm(app, session)
	if err != nil {
		app.SendMessage("Не удалось создать фильм", nil)
		films.HandleFilmsCommand(app, session)
		return
	}

	msg := fmt.Sprintf("Фильм %q успешно добавлен в коллекцию %q\n", collectionFilm.Film.Title, collectionFilm.Collection.Name)
	msg += messages.BuildFilmDetailMessage(&collectionFilm.Film)

	imageURL := collectionFilm.Film.ImageURL
	app.SendImage(imageURL, msg, nil)

	session.ClearAllStates()
	films.HandleFilmsCommand(app, session)
}
