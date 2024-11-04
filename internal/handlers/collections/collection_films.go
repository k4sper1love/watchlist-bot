package collections

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log"
	"strconv"
	"strings"
)

func HandleCollectionFilmsCommand(app models.App, session *models.Session) {
	collectionFilmsResponse, err := getCollectionFilms(app, session)
	if err != nil {
		app.SendMessage(err.Error(), nil)
		return
	}

	metadata := collectionFilmsResponse.Metadata

	msg := builders.BuildCollectionFilmsMessage(collectionFilmsResponse)

	keyboard := builders.NewKeyboard(1).
		AddCollectionFilmsSelect(collectionFilmsResponse).
		AddCollectionFilmsNew().
		AddCollectionsUpdate().
		AddCollectionsDelete().
		AddNavigation(metadata.CurrentPage, metadata.LastPage, states.CallbackCollectionFilmsPrevPage, states.CallbackCollectionFilmsNextPage).
		AddBack(states.CallbackCollectionFilmsBack).
		Build()

	app.SendMessage(msg, keyboard)
}

func HandleCollectionFilmsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)

	switch {
	case callback == states.CallbackCollectionFilmsNextPage:
		if session.CollectionDetailState.CurrentPage < session.CollectionDetailState.LastPage {
			session.CollectionDetailState.CurrentPage++
			HandleCollectionFilmsCommand(app, session)
		} else {
			app.SendMessage("Вы уже на последней странице", nil)
		}

	case callback == states.CallbackCollectionFilmsPrevPage:
		if session.CollectionDetailState.CurrentPage > 1 {
			session.CollectionDetailState.CurrentPage--
			HandleCollectionFilmsCommand(app, session)
		} else {
			app.SendMessage("Вы уже на первой странице", nil)
		}

	case callback == states.CallbackCollectionFilmsNew:
		HandleNewCollectionFilmCommand(app, session)

	case callback == states.CallbackCollectionFilmsBack:
		HandleCollectionsCommand(app, session)

	case callback == states.CallbackCollectionFilmsDelete:
		HandleDeleteCollectionFilmCommand(app, session)

	case callback == states.CallbackCollectionFilmsUpdate:
		HandleUpdateCollectionFilmCommand(app, session)

	case strings.HasPrefix(callback, "select_cf_"):
		handleCollectionFilmSelect(app, session)
	}
}

func handleCollectionFilmSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	indexStr := strings.TrimPrefix(callback, "select_cf_")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		app.SendMessage("Ошибка при получении ID коллекции.", nil)
		log.Printf("error parsing collection film index: %v", err)
		return
	}

	session.CollectionFilmState.Index = index
	HandleCollectionFilmsDetailCommand(app, session)
}

func getCollectionFilms(app models.App, session *models.Session) (*models.CollectionFilmsResponse, error) {
	if session.CollectionDetailState.ObjectID == -1 {
		return nil, fmt.Errorf("Неизвестная коллекция. Начните с начала")
	}

	collectionFilmsResponse, err := watchlist.GetCollectionFilms(app, session)
	if err != nil {
		return nil, err
	}

	session.CollectionDetailState.Object = collectionFilmsResponse.CollectionFilms
	session.CollectionDetailState.LastPage = collectionFilmsResponse.Metadata.LastPage
	session.CollectionDetailState.TotalRecords = collectionFilmsResponse.Metadata.TotalRecords

	return collectionFilmsResponse, nil
}

func updateCollectionFilmsList(app models.App, session *models.Session, next bool) error {
	currentPage := session.CollectionDetailState.CurrentPage
	lastPage := session.CollectionDetailState.LastPage

	switch next {
	case true:
		if currentPage < lastPage {
			session.CollectionDetailState.CurrentPage++
		} else {
			return fmt.Errorf("Вы уже на последней странице")
		}
	case false:
		if currentPage > 1 {
			session.CollectionDetailState.CurrentPage--
		} else {
			return fmt.Errorf("Вы уже на первой странице")
		}
	}

	_, err := getCollectionFilms(app, session)
	if err != nil {
		return err
	}

	return nil
}
