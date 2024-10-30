package handlers

import (
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"log"
)

func handleCollectionFilmsDetailCommand(app config.App, session *models.Session) {
	index := session.CollectionFilmState.Index
	collectionFilms := session.CollectionState.CollectionFilms

	if index == -1 || index >= len(collectionFilms.Films) {
		sendMessage(app, "Неизвестный фильм в коллекции. Начните с начала")
		log.Println(index)
		return
	}

	if len(collectionFilms.Films) == 0 {
		sendMessage(app, "Фильмы не найдены. Начните с начала")
		return
	}

	msg := builders.BuildCollectionFilmsDetailMessage(collectionFilms, index)

	buttons := builders.BuildNavigationButtons(index+1, len(collectionFilms.Films), CallbackCollectionFilmsDetailPrevPage, CallbackCollectionFilmsDetailNextPage)

	keyboard := builders.BuildButtonKeyboard(buttons, 2)

	sendImageWithTextAndKeyboard(app, collectionFilms.Films[index].ImageURL, msg, keyboard)
}

func handleCollectionFilmsDetailButtons(app config.App, session *models.Session) {
	currentPage := session.CollectionFilmState.Index + 1
	lastPage := len(session.CollectionState.CollectionFilms.Films)

	switch {
	case session.State == CallbackCollectionFilmsDetailNextPage:
		if currentPage < lastPage {
			session.CollectionFilmState.Index++
			handleCollectionFilmsDetailCommand(app, session)
		} else {
			sendMessage(app, "Вы уже на последней странице")
		}
		resetState(session)

	case session.State == CallbackCollectionFilmsDetailPrevPage:
		if currentPage > 1 {
			session.CollectionFilmState.Index--
			handleCollectionFilmsDetailCommand(app, session)
		} else {
			sendMessage(app, "Вы уже на первой странице")
		}
		resetState(session)
	}
}
