package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func BuildOptionsFilmToCollectionKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddNewFilmToCollection()

	keyboard.AddExistingFilmToCollection()

	keyboard.AddBack(states.CallbackOptionsFilmToCollectionBack)

	return keyboard.Build(session.Lang)
}

func BuildAddFilmToCollectionKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	films := session.FilmsState.Films
	currentPage := session.CollectionFilmsState.CurrentPage
	lastPage := session.CollectionFilmsState.LastPage

	keyboard := NewKeyboard()

	if session.FilmsState.Title == "" {
		keyboard.AddSearch(states.CallbackAddFilmToCollectionFind)
	} else {
		keyboard.AddReset(states.CallbackAddFilmToCollectionReset)
	}

	keyboard.AddCollectionFilmSelectFilm(films)

	keyboard.AddNavigation(
		currentPage,
		lastPage,
		states.CallbackAddFilmToCollectionPrevPage,
		states.CallbackAddFilmToCollectionNextPage,
		states.CallbackAddFilmToCollectionFirstPage,
		states.CallbackAddFilmToCollectionLastPage,
	)

	keyboard.AddBack(states.CallbackAddFilmToCollectionBack)

	return keyboard.Build(session.Lang)
}

func BuildAddFilmToCollectionNotFoundKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	if session.FilmsState.Title != "" {
		keyboard.AddAgain(states.CallbackAddFilmToCollectionAgain)
	} else {
		keyboard.AddBack(states.CallbackAddFilmToCollectionBack)
	}

	return keyboard.Build(session.Lang)
}

func BuildAddCollectionToFilmKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	collections := session.CollectionsState.Collections
	currentPage := session.CollectionFilmsState.CurrentPage
	lastPage := session.CollectionFilmsState.LastPage

	keyboard := NewKeyboard()

	if session.CollectionsState.Name == "" {
		keyboard.AddSearch(states.CallbackAddCollectionToFilmFind)
	} else {
		keyboard.AddReset(states.CallbackAddCollectionToFilmReset)
	}

	keyboard.AddCollectionFilmSelectCollection(collections)

	keyboard.AddNavigation(
		currentPage,
		lastPage,
		states.CallbackAddCollectionToFilmPrevPage,
		states.CallbackAddCollectionToFilmNextPage,
		states.CallbackAddCollectionToFilmFirstPage,
		states.CallbackAddCollectionToFilmLastPage,
	)

	keyboard.AddBack(states.CallbackAddCollectionToFilmBack)

	return keyboard.Build(session.Lang)
}

func BuildAddCollectionToFilmNotFoundKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	if session.CollectionsState.Name != "" {
		keyboard.AddAgain(states.CallbackAddCollectionToFilmAgain)
	} else {
		keyboard.AddBack(states.CallbackAddCollectionToFilmBack)
	}

	return keyboard.Build(session.Lang)
}

func (k *Keyboard) AddCollectionFilmFromFilm() *Keyboard {
	return k.AddButton("➕", "addCollectionToFilm", states.CallbackCollectionFilmsFromFilm, "", true)
}

func (k *Keyboard) AddCollectionFilmFromCollection() *Keyboard {
	return k.AddButton("➕", "addFilmToCollection", states.CallbackCollectionFilmsFromCollection, "", true)
}

func (k *Keyboard) AddNewFilmToCollection() *Keyboard {
	return k.AddButton("🆕", "createFilm", states.CallbackOptionsFilmToCollectionNew, "", true)
}

func (k *Keyboard) AddExistingFilmToCollection() *Keyboard {
	return k.AddButton("\U0001F7F0", "choiceFromFilms", states.CallbackOptionsFilmToCollectionExisting, "", true)
}

func (k *Keyboard) AddCollectionFilmSelectFilm(films []apiModels.Film) *Keyboard {
	for _, film := range films {
		k.AddButton("",
			fmt.Sprintf("%s %s (%d)", utils.BoolToStarOrEmpty(film.IsFavorite), film.Title, film.ID),
			fmt.Sprintf("select_cf_film_%d", film.ID),
			"",
			false,
		)
	}

	return k
}

func (k *Keyboard) AddCollectionFilmSelectCollection(collections []apiModels.Collection) *Keyboard {
	for _, collection := range collections {
		k.AddButton(
			"",
			fmt.Sprintf("%s %s (%d)", utils.BoolToStarOrEmpty(collection.IsFavorite), collection.Name, collection.ID),
			fmt.Sprintf("select_cf_collection_%d", collection.ID),
			"",
			false)
	}

	return k
}
