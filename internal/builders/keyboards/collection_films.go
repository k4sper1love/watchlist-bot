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
		AddBack(states.CallbackOptionsFilmToCollectionBack).
		Build(session.Lang)
}

func AddFilmToCollection(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	state := session.CollectionFilmsState
	return New().
		AddIf(session.FilmsState.Title == "", func(k *Keyboard) {
			k.AddSearch(states.CallbackAddFilmToCollectionFind)
		}).
		AddIf(session.FilmsState.Title != "", func(k *Keyboard) {
			k.AddReset(states.CallbackAddFilmToCollectionReset)
		}).
		AddCollectionFilmSelectFilm(session.FilmsState.Films).
		AddNavigation(state.CurrentPage, state.LastPage, states.PrefixAddFilmToCollectionPage, true).
		AddBack(states.CallbackAddFilmToCollectionBack).
		Build(session.Lang)
}

func FilmToCollectionNotFound(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(session.FilmsState.Title != "", func(k *Keyboard) {
			k.AddAgain(states.CallbackAddFilmToCollectionAgain)
		}).
		AddIf(session.FilmsState.Title == "", func(k *Keyboard) {
			k.AddBack(states.CallbackAddFilmToCollectionBack)
		}).
		Build(session.Lang)
}

func AddCollectionToFilm(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	state := session.CollectionFilmsState
	return New().
		AddIf(session.CollectionsState.Name == "", func(k *Keyboard) {
			k.AddSearch(states.CallbackAddCollectionToFilmFind)
		}).
		AddIf(session.CollectionsState.Name != "", func(k *Keyboard) {
			k.AddReset(states.CallbackAddCollectionToFilmReset)
		}).
		AddCollectionFilmSelectCollection(session.CollectionsState.Collections).
		AddNavigation(state.CurrentPage, state.LastPage, states.PrefixAddCollectionToFilmPage, true).
		AddBack(states.CallbackAddCollectionToFilmBack).
		Build(session.Lang)

}

func CollectionToFilmNotFound(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(session.CollectionsState.Name != "", func(k *Keyboard) {
			k.AddAgain(states.CallbackAddCollectionToFilmAgain)
		}).
		AddIf(session.CollectionsState.Name == "", func(k *Keyboard) {
			k.AddBack(states.CallbackAddCollectionToFilmBack)
		}).
		Build(session.Lang)
}
