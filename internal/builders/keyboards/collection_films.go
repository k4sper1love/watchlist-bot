package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func BuildOptionsFilmToCollectionKeyboard() *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddNewFilmToCollection()

	keyboard.AddExistingFilmToCollection()

	keyboard.AddBack(states.CallbackOptionsFilmToCollectionBack)

	return keyboard.Build()
}

func BuildAddFilmToCollectionKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	films := session.FilmsState.Films
	currentPage := session.CollectionFilmsState.CurrentPage
	lastPage := session.CollectionFilmsState.LastPage

	keyboard := NewKeyboard()

	keyboard.AddCollectionFilmSelectFilm(films)

	keyboard.AddNavigation(
		currentPage,
		lastPage,
		states.CallbackAddFilmToCollectionPrevPage,
		states.CallbackAddFilmToCollectionNextPage,
	)

	keyboard.AddBack(states.CallbackAddFilmToCollectionBack)

	return keyboard.Build()
}

func BuildAddCollectionToFilmKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	collections := session.CollectionsState.Collections
	currentPage := session.CollectionFilmsState.CurrentPage
	lastPage := session.CollectionFilmsState.LastPage

	keyboard := NewKeyboard()

	keyboard.AddCollectionFilmSelectCollection(collections)

	keyboard.AddNavigation(
		currentPage,
		lastPage,
		states.CallbackAddCollectionToFilmPrevPage,
		states.CallbackAddCollectionToFilmNextPage,
	)

	keyboard.AddBack(states.CallbackAddCollectionToFilmBack)

	return keyboard.Build()
}

func (k *Keyboard) AddCollectionFilmFromFilm() *Keyboard {
	return k.AddButton("➕ Добавить фильм в коллекцию", states.CallbackCollectionFilmsFromFilm)
}

func (k *Keyboard) AddCollectionFilmFromCollection() *Keyboard {
	return k.AddButton("➕ Добавить фильм", states.CallbackCollectionFilmsFromCollection)
}

func (k *Keyboard) AddNewFilmToCollection() *Keyboard {
	return k.AddButton("🆕 Создать фильм", states.CallbackOptionsFilmToCollectionNew)
}

func (k *Keyboard) AddExistingFilmToCollection() *Keyboard {
	return k.AddButton("\U0001F7F0 Добавить существующий фильм", states.CallbackOptionsFilmToCollectionExisting)
}

func (k *Keyboard) AddCollectionFilmSelectFilm(films []apiModels.Film) *Keyboard {
	for _, film := range films {
		k.AddButton(fmt.Sprintf("%s (%d)", film.Title, film.ID), fmt.Sprintf("select_cf_film_%d", film.ID))
	}

	return k
}

func (k *Keyboard) AddCollectionFilmSelectCollection(collections []apiModels.Collection) *Keyboard {
	for _, collection := range collections {
		k.AddButton(fmt.Sprintf("%s (%d)", collection.Name, collection.ID), fmt.Sprintf("select_cf_collection_%d", collection.ID))
	}

	return k
}
