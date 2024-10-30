package handlers

import (
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"log"
	"strconv"
	"strings"
)

func handleCollectionFilmsCommand(app config.App, session *models.Session) {
	if session.CollectionState.ObjectID == -1 {
		sendMessage(app, "Неизвестный ID коллекции. Начните с начала")
		return
	}

	collectionFilmsResponse, err := watchlist.GetCollectionFilms(app, session)
	if err != nil {
		sendMessage(app, err.Error())
		return
	}

	if collectionFilmsResponse.Metadata.TotalRecords == 0 {
		sendMessage(app, "В этой коллекции нет фильмов")
		return
	}

	session.CollectionState.CollectionFilms = collectionFilmsResponse.CollectionFilms

	currentPage := collectionFilmsResponse.Metadata.CurrentPage
	lastPage := collectionFilmsResponse.Metadata.LastPage

	session.CollectionFilmState.CurrentPage = currentPage
	session.CollectionFilmState.LastPage = lastPage

	msg := builders.BuildCollectionFilmsMessage(collectionFilmsResponse)

	buttons := builders.BuildCollectionFilmsSelectButtons(collectionFilmsResponse)
	buttons = append(buttons, builders.BuildNavigationButtons(currentPage, lastPage, CallbackCollectionFilmsPrevPage, CallbackCollectionFilmsNextPage)...)

	keyboard := builders.BuildButtonKeyboard(buttons, 1)

	if lastPage != 1 {
		sendMessageWithKeyboard(app, msg, keyboard)
	} else {
		sendMessage(app, msg)
	}
}

func handleCollectionFilmsButtons(app config.App, session *models.Session) {
	switch {
	case session.State == CallbackCollectionFilmsNextPage:
		if session.CollectionFilmState.CurrentPage < session.CollectionFilmState.LastPage {
			session.CollectionFilmState.CurrentPage++
			handleCollectionFilmsCommand(app, session)
		} else {
			sendMessage(app, "Вы уже на последней странице")
		}
		resetState(session)

	case session.State == CallbackCollectionFilmsPrevPage:
		if session.CollectionFilmState.CurrentPage > 1 {
			session.CollectionFilmState.CurrentPage--
			handleCollectionFilmsCommand(app, session)
		} else {
			sendMessage(app, "Вы уже на первой странице")
		}
		resetState(session)
	case strings.HasPrefix(session.State, "select_cf_"):
		handleCollectionFilmSelect(app, session)
		resetState(session)
	}
}

func handleCollectionFilmSelect(app config.App, session *models.Session) {
	listIndexStr := strings.TrimPrefix(session.State, "select_cf_")
	listIndex, err := strconv.Atoi(listIndexStr)

	if err != nil {
		sendMessage(app, "Ошибка при получении ID коллекции.")
		log.Printf("error parsing collection film index: %v", err)
		return
	}

	session.CollectionFilmState.Index = listIndex
	handleCollectionFilmsDetailCommand(app, session)
}
