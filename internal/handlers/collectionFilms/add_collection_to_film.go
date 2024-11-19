package collectionFilms

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log"
	"strconv"
	"strings"
)

func HandleAddCollectionToFilmCommand(app models.App, session *models.Session) {
	if err := GetCollectionsExcludeFilm(app, session); err != nil {
		app.SendMessage(err.Error(), nil)
		return
	}

	msg := "Выберите, в какую коллекцию добавить фильм?"

	keyboard := keyboards.BuildAddCollectionToFilmKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func HandleAddCollectionToFilmButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	switch {
	case callback == states.CallbackAddCollectionToFilmBack:
		films.HandleFilmsDetailCommand(app, session)

	case strings.HasPrefix(callback, "select_cf_collection_"):
		HandleAddCollectionToFilmSelect(app, session)

	case callback == states.CallbackAddCollectionToFilmNextPage:
		if session.CollectionFilmsState.CurrentPage < session.CollectionFilmsState.LastPage {
			session.CollectionFilmsState.CurrentPage++
			HandleAddCollectionToFilmCommand(app, session)
		} else {
			app.SendMessage("Вы уже на последней странице", nil)
		}

	case callback == states.CallbackAddCollectionToFilmPrevPage:
		if session.CollectionFilmsState.CurrentPage > 1 {
			session.CollectionFilmsState.CurrentPage--
			HandleAddCollectionToFilmCommand(app, session)
		} else {
			app.SendMessage("Вы уже на первой странице", nil)
		}
	}
}
func HandleAddCollectionToFilmSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	idStr := strings.TrimPrefix(callback, "select_cf_collection_")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		app.SendMessage("Ошибка при получении ID коллекции.", nil)
		log.Printf("error parsing collection ID: %v", err)
		return
	}

	session.CollectionDetailState.Collection.ID = id

	addFilmToCollection(app, session)
}

func GetCollectionsExcludeFilm(app models.App, session *models.Session) error {
	collectionsResponse, err := watchlist.GetCollectionsExcludeFilm(app, session)
	if err != nil {
		return err
	}

	session.CollectionsState.Collections = collectionsResponse.Collections
	session.CollectionFilmsState.CurrentPage = collectionsResponse.Metadata.CurrentPage
	session.CollectionFilmsState.LastPage = collectionsResponse.Metadata.LastPage

	return nil
}
