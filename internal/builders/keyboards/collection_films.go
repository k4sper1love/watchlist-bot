package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func FilmToCollectionOptions(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddNewFilmToCollection().
		AddExistingFilmToCollection().
		AddBack(states.CallFilmToCollectionOptionBack).
		Build(session.Lang)
}

func AddFilmToCollection(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	state := session.CollectionFilmsState
	return New().
		AddIf(session.FilmsState.Title == "", func(k *Keyboard) {
			k.AddSearch(states.CallAddFilmToCollectionFind)
		}).
		AddIf(session.FilmsState.Title != "", func(k *Keyboard) {
			k.AddReset(states.CallAddFilmToCollectionReset)
		}).
		AddCollectionFilmSelectFilm(session.FilmsState.Films).
		AddNavigation(state.CurrentPage, state.LastPage, states.AddFilmToCollectionPage, true).
		AddBack(states.CallAddFilmToCollectionBack).
		Build(session.Lang)
}

func FilmToCollectionNotFound(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(session.FilmsState.Title != "", func(k *Keyboard) {
			k.AddAgain(states.CallAddFilmToCollectionAgain)
		}).
		AddIf(session.FilmsState.Title == "", func(k *Keyboard) {
			k.AddBack(states.CallAddFilmToCollectionBack)
		}).
		Build(session.Lang)
}

func AddCollectionToFilm(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	state := session.CollectionFilmsState
	return New().
		AddIf(session.CollectionsState.Name == "", func(k *Keyboard) {
			k.AddSearch(states.CallAddCollectionToFilmFind)
		}).
		AddIf(session.CollectionsState.Name != "", func(k *Keyboard) {
			k.AddReset(states.CallAddCollectionToFilmReset)
		}).
		AddCollectionFilmSelectCollection(session.CollectionsState.Collections).
		AddNavigation(state.CurrentPage, state.LastPage, states.AddCollectionToFilmPage, true).
		AddBack(states.CallAddCollectionToFilmBack).
		Build(session.Lang)

}

func CollectionToFilmNotFound(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(session.CollectionsState.Name != "", func(k *Keyboard) {
			k.AddAgain(states.CallAddCollectionToFilmAgain)
		}).
		AddIf(session.CollectionsState.Name == "", func(k *Keyboard) {
			k.AddBack(states.CallAddCollectionToFilmBack)
		}).
		Build(session.Lang)
}
