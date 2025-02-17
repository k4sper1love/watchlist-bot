package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

var updateFilmURLButton = Button{"", "filmURL", states.CallbackUpdateFilmSelectURL, "", true}

var updateFilmButtons = []Button{
	{"", "image", states.CallbackUpdateFilmSelectImage, "", true},
	{"", "title", states.CallbackUpdateFilmSelectTitle, "", true},
	{"", "description", states.CallbackUpdateFilmSelectDescription, "", true},
	{"", "genre", states.CallbackUpdateFilmSelectGenre, "", true},
	{"", "rating", states.CallbackUpdateFilmSelectRating, "", true},
	{"", "yearOfRelease", states.CallbackUpdateFilmSelectYear, "", true},
	{"", "comment", states.CallbackUpdateFilmSelectComment, "", true},
	{"", "viewed", states.CallbackUpdateFilmSelectViewed, "", true},
}

var updateFilmsAfterViewedButtons = []Button{
	{"", "userRating", states.CallbackUpdateFilmSelectUserRating, "", true},
	{"", "review", states.CallbackUpdateFilmSelectReview, "", true},
}

func BuildFilmsKeyboard(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	if len(session.FilmsState.Films) > 0 {
		keyboard.AddSearch(states.CallbackFilmsFind)
	}

	keyboard.AddFilmSelect(session)

	keyboard.AddNavigation(
		currentPage,
		lastPage,
		states.CallbackFilmsPrevPage,
		states.CallbackFilmsNextPage,
		states.CallbackFilmsFirstPage,
		states.CallbackFilmsLastPage,
	)

	keyboard.AddFilmFiltersAndSorting(session)

	switch session.Context {
	case states.ContextFilm:
		keyboard.AddFilmNew()
		keyboard.AddBack("")

	case states.ContextCollection:
		keyboard.AddCollectionFilmFromCollection()
		keyboard.AddFavorite(session.CollectionDetailState.Collection.IsFavorite, states.CallbackCollectionsFavorite)
		keyboard.AddManage(states.CallbackCollectionsManage)
		keyboard.AddBack(states.CallbackFilmsBack)
	}

	return keyboard.Build(session.Lang)
}

func BuildFindFilmsKeyboard(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddFilmSelect(session)

	keyboard.AddNavigation(
		currentPage,
		lastPage,
		states.CallbackFindFilmsPrevPage,
		states.CallbackFindFilmsNextPage,
		states.CallbackFindFilmsFirstPage,
		states.CallbackFindFilmsLastPage,
	)

	keyboard.AddBack(states.CallbackFindFilmsBack)

	return keyboard.Build(session.Lang)
}

func BuildFindNewFilmKeyboard(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddFindNewFilmSelect(session)

	keyboard.AddNavigation(
		currentPage,
		lastPage,
		states.CallbackFindNewFilmPrevPage,
		states.CallbackFindNewFilmNextPage,
		states.CallbackFindNewFilmFirstPage,
		states.CallbackFindNewFilmLastPage,
	)

	keyboard.AddBack(states.CallbackFindNewFilmBack)

	return keyboard.Build(session.Lang)
}

func BuildFilmDetailKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	state := session.FilmDetailState
	film := state.Film

	keyboard := NewKeyboard()

	keyboard.AddFavorite(film.IsFavorite, states.CallbackFilmDetailFavorite)

	if film.URL != "" {
		keyboard.AddURLButton("", "openInBrowser", film.URL, true)
	}

	if !film.IsViewed {
		keyboard.AddButton("✔️", "viewed", states.CallbackFilmDetailViewed, "", true)
	}

	keyboard.AddManage(states.CallbackFilmsManage)

	if session.Context == states.ContextFilm {
		keyboard.AddCollectionFilmFromFilm()
	}

	if state.HasIndex() {
		itemID := utils.GetItemID(state.Index, session.FilmsState.CurrentPage, session.FilmsState.PageSize)

		keyboard.AddNavigation(
			itemID,
			session.FilmsState.TotalRecords,
			states.CallbackFilmDetailPrevPage,
			states.CallbackFilmDetailNextPage,
			"",
			"",
		)
	}

	keyboard.AddBack(states.CallbackFilmDetailBack)

	return keyboard.Build(session.Lang)
}

func BuildFilmManageKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	if session.Context == states.ContextCollection {
		keyboard.AddFilmRemoveFromCollection()
	}

	keyboard.AddUpdate(states.CallbackManageFilmSelectUpdate)

	keyboard.AddDelete(states.CallbackManageFilmSelectDelete)

	keyboard.AddBack(states.CallbackManageFilmSelectBack)

	return keyboard.Build(session.Lang)
}

func BuildFilmNewKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddNewFilmManually()

	keyboard.AddNewFilmFromURL()

	keyboard.AddNewFilmFind()

	keyboard.AddBack(states.CallbackNewFilmSelectBack)

	return keyboard.Build(session.Lang)
}

func BuildFilmUpdateKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	film := session.FilmDetailState.Film
	keyboard := NewKeyboard()

	keyboard.AddRow(updateFilmURLButton)

	keyboard.AddButtonsWithRowSize(2, updateFilmButtons...)

	if film.IsViewed {
		keyboard.AddButtonsWithRowSize(2, updateFilmsAfterViewedButtons...)
	}

	keyboard.AddBack(states.CallbackUpdateFilmSelectBack)

	return keyboard.Build(session.Lang)
}

func BuildFilmViewedKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddSkip()

	keyboard.AddCancel()

	return keyboard.Build(session.Lang)
}

func BuildFilmsFilterKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	filter := session.GetFilmsFiltersByContext()

	keyboard := NewKeyboard()

	keyboard.AddButtons(parseFiltersFilmsButtons(filter, session.Lang)...)

	if filter.IsFiltersEnabled() {
		keyboard.AddResetAllFilmsFilters()
	}

	keyboard.AddBack(states.CallbackFiltersFilmsSelectBack)

	return keyboard.Build(session.Lang)
}

func BuildFilmsSortingKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	sorting := session.GetFilmsSortingByContext()

	keyboard := NewKeyboard()

	keyboard.AddButtons(parseSortingFilmsButtons(sorting, session.Lang)...)

	if sorting.IsSortingEnabled() {
		keyboard.AddResetAllSorting(states.CallbackSortingFilmsSelectAllReset)
	}

	keyboard.AddBack(states.CallbackSortingFilmsSelectBack)

	return keyboard.Build(session.Lang)
}

func (k *Keyboard) AddFilmSelect(session *models.Session) *Keyboard {
	var buttons []Button

	for i, film := range session.FilmsState.Films {
		itemID := utils.GetItemID(i, session.FilmsState.CurrentPage, session.FilmsState.PageSize)

		buttons = append(buttons, Button{"", fmt.Sprintf("%s %s", utils.NumberToEmoji(itemID), film.Title), fmt.Sprintf("select_film_%d", i), "", false})
	}

	k.AddButtonsWithRowSize(2, buttons...)

	return k
}

func (k *Keyboard) AddFindNewFilmSelect(session *models.Session) *Keyboard {
	var buttons []Button

	for i, film := range session.FilmsState.Films {
		itemID := utils.GetItemID(i, session.FilmsState.CurrentPage, session.FilmsState.PageSize)

		buttons = append(buttons, Button{"", fmt.Sprintf("%s %s", utils.NumberToEmoji(itemID), film.Title), fmt.Sprintf("select_find_new_film_%d", i), "", false})
	}

	k.AddButtonsWithRowSize(2, buttons...)

	return k
}

func (k *Keyboard) AddFilmFind() *Keyboard {
	return k.AddButton("🔎", "findFilmByTitle", states.CallbackFilmsFind, "", true)
}

func (k *Keyboard) AddFilmNew() *Keyboard {
	return k.AddButton("➕", "createFilm", states.CallbackFilmsNew, "", true)
}

func (k *Keyboard) AddFilmFiltersAndSorting(session *models.Session) *Keyboard {
	filtersEnable := session.GetFilmsFiltersByContext().IsFiltersEnabled()
	sortingEnable := session.GetFilmsSortingByContext().IsSortingEnabled()

	return k.AddButtonsWithRowSize(2,
		Button{utils.BoolToEmoji(sortingEnable), "sorting", states.CallbackFilmsSorting, "", true},
		Button{utils.BoolToEmoji(filtersEnable), "filters", states.CallbackFilmsFilters, "", true},
	)
}

func (k *Keyboard) AddFilmDelete() *Keyboard {
	return k.AddButton("🗑️", "deleteFilm", states.CallbackManageFilmSelectDelete, "", true)
}

func (k *Keyboard) AddFilmUpdate() *Keyboard {
	return k.AddButton("✏️", "updateFilm", states.CallbackManageFilmSelectUpdate, "", true)
}

func (k *Keyboard) AddFilmManage() *Keyboard {
	return k.AddButton("⚙️", "manageFilm", states.CallbackFilmsManage, "", true)
}

func (k *Keyboard) AddFilmRemoveFromCollection() *Keyboard {
	return k.AddButton("🧹󠁝", "removeFromCollection", states.CallbackManageFilmSelectRemoveFromCollection, "", true)
}

func (k *Keyboard) AddNewFilmManually() *Keyboard {
	return k.AddButton("", "manually", states.CallbackNewFilmSelectManually, "", true)
}

func (k *Keyboard) AddNewFilmFromURL() *Keyboard {
	return k.AddButton("", "fromURL", states.CallbackNewFilmSelectFromURL, "", true)
}

func (k *Keyboard) AddNewFilmFind() *Keyboard {
	return k.AddButton("", "findFilm", states.CallbackNewFilmSelectFind, "", true)
}

func (k *Keyboard) AddAgain(callback string) *Keyboard {
	return k.AddButton("↻", "again", callback, "", true)
}

func (k *Keyboard) AddResetAllFilmsFilters() *Keyboard {
	return k.AddButton("🔄", "resetFilters", states.CallbackFiltersFilmsSelectAllReset, "", true)
}

func (k *Keyboard) AddResetFilmsFilter(session *models.Session, filterType string) *Keyboard {
	filter := session.GetFilmsFiltersByContext()

	if filter.IsFilterEnabled(filterType) {
		return k.AddButton("🔄", "reset", states.CallbackProcessReset, "", true)
	}

	return k
}

func (k *Keyboard) AddFavorite(isFavorite bool, callback string) *Keyboard {
	var messageCode string
	if isFavorite {
		messageCode = "removeFavorite"
	} else {
		messageCode = "makeFavorite"
	}

	return k.AddButton(utils.BoolToStar(!isFavorite), messageCode, callback, "", true)
}

func (k *Keyboard) AddChangeToken() *Keyboard {
	return k.AddButton("🔄", "changeToken", states.CallbackNewFilmSelectChangeKinopoiskToken, "", true)
}

func parseFiltersFilmsButtons(filter *models.FiltersFilm, lang string) []Button {
	var buttons []Button

	buttons = addFiltersFilmsButton(buttons, filter, lang, "isFavorite", states.CallbackFiltersFilmsSelectIsFavorite, true)

	buttons = addFiltersFilmsButton(buttons, filter, lang, "isViewed", states.CallbackFiltersFilmsSelectIsViewed, true)

	buttons = addFiltersFilmsButton(buttons, filter, lang, "year", states.CallbackFiltersFilmsSelectYear, false)

	buttons = addFiltersFilmsButton(buttons, filter, lang, "rating", states.CallbackFiltersFilmsSelectRating, false)

	buttons = addFiltersFilmsButton(buttons, filter, lang, "userRating", states.CallbackFiltersFilmsSelectUserRating, false)

	buttons = addFiltersFilmsButton(buttons, filter, lang, "hasURL", states.CallbackFiltersFilmsSelectHasURL, true)

	return buttons
}

func parseSortingFilmsButtons(sorting *models.Sorting, lang string) []Button {
	var buttons []Button

	buttons = addSortingButton(buttons, sorting, lang, "is_favorite", states.CallbackSortingFilmsSelectIsFavorite)

	buttons = addSortingButton(buttons, sorting, lang, "is_viewed", states.CallbackSortingFilmsSelectIsViewed)

	buttons = addSortingButton(buttons, sorting, lang, "title", states.CallbackSortingFilmsSelectTitle)

	buttons = addSortingButton(buttons, sorting, lang, "year", states.CallbackSortingFilmsSelectYear)

	buttons = addSortingButton(buttons, sorting, lang, "rating", states.CallbackSortingFilmsSelectRating)

	buttons = addSortingButton(buttons, sorting, lang, "user_rating", states.CallbackSortingFilmsSelectUserRating)

	buttons = addSortingButton(buttons, sorting, lang, "created_at", states.CallbackSortingFilmsSelectCreatedAt)

	return buttons
}
