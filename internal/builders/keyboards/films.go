package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

var updateFilmURLButton = Button{"", "filmURL", states.CallbackUpdateFilmSelectURL, ""}

var updateFilmButtons = []Button{
	{"", "image", states.CallbackUpdateFilmSelectImage, ""},
	{"", "title", states.CallbackUpdateFilmSelectTitle, ""},
	{"", "description", states.CallbackUpdateFilmSelectDescription, ""},
	{"", "genre", states.CallbackUpdateFilmSelectGenre, ""},
	{"", "rating", states.CallbackUpdateFilmSelectRating, ""},
	{"", "yearOfRelease", states.CallbackUpdateFilmSelectYear, ""},
	{"", "comment", states.CallbackUpdateFilmSelectComment, ""},
	{"", "viewed", states.CallbackUpdateFilmSelectViewed, ""},
}

var updateFilmsAfterViewedButtons = []Button{
	{"", "userRating", states.CallbackUpdateFilmSelectUserRating, ""},
	{"", "Review", states.CallbackUpdateFilmSelectReview, ""},
}

func BuildFilmsKeyboard(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	if len(session.FilmsState.Films) > 0 {
		keyboard.AddFilmFind()
	}

	keyboard.AddFilmSelect(session)

	keyboard.AddNavigation(
		currentPage,
		lastPage,
		states.CallbackFilmsPrevPage,
		states.CallbackFilmsNextPage,
	)

	keyboard.AddFilmFiltersAndSorting(session)

	switch session.Context {
	case states.ContextFilm:
		keyboard.AddFilmNew()
		keyboard.AddBack("")

	case states.ContextCollection:
		keyboard.AddCollectionFilmFromCollection()
		keyboard.AddCollectionsManage()
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
	)

	keyboard.AddBack(states.CallbackFindFilmsBack)

	return keyboard.Build(session.Lang)
}

func BuildFilmDetailKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	index := session.FilmDetailState.Index
	itemID := utils.GetItemID(index, session.FilmsState.CurrentPage, session.FilmsState.PageSize)
	film := session.FilmDetailState.Film

	keyboard := NewKeyboard()

	if film.URL != "" {
		keyboard.AddURLButton("", "openInBrowser", film.URL)
	}

	if !film.IsViewed {
		keyboard.AddButton("✔️", "viewed", states.CallbackFilmDetailViewed, "")
	}

	keyboard.AddFilmManage()

	if session.Context == states.ContextFilm {
		keyboard.AddCollectionFilmFromFilm()
	}

	keyboard.AddNavigation(
		itemID,
		session.FilmsState.TotalRecords,
		states.CallbackFilmDetailPrevPage,
		states.CallbackFilmDetailNextPage)

	keyboard.AddBack(states.CallbackFilmDetailBack)

	return keyboard.Build(session.Lang)
}

func BuildFilmManageKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddFilmUpdate()

	if session.Context == states.ContextCollection {
		keyboard.AddFilmRemoveFromCollection()
	}

	keyboard.AddFilmDelete()

	keyboard.AddBack(states.CallbackManageFilmSelectBack)

	return keyboard.Build(session.Lang)
}

func BuildFilmNewKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddNewFilmManually()

	keyboard.AddNewFilmFromURL()

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
	keyboard := NewKeyboard()

	keyboard.AddButtons(parseFiltersFilmsButtons(session)...)

	keyboard.AddResetAllFilters()

	keyboard.AddBack(states.CallbackFiltersFilmsSelectBack)

	return keyboard.Build(session.Lang)
}

func BuildFilmsSortingKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddButtons(parseSortingFilmsButtons(session)...)

	keyboard.AddResetAllSorting()

	keyboard.AddBack(states.CallbackSortingFilmsSelectBack)

	return keyboard.Build(session.Lang)
}

func BuildFilmsSortingDirectionKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddSortingDirection()

	keyboard.AddResetSorting(session)

	keyboard.AddCancel()

	return keyboard.Build(session.Lang)
}

func (k *Keyboard) AddFilmSelect(session *models.Session) *Keyboard {
	var buttons []Button

	for i, film := range session.FilmsState.Films {
		itemID := utils.GetItemID(i, session.FilmsState.CurrentPage, session.FilmsState.PageSize)

		buttons = append(buttons, Button{"", fmt.Sprintf("%s (%d)", film.Title, itemID), fmt.Sprintf("select_film_%d", i), ""})
	}

	k.AddButtonsWithRowSize(2, buttons...)

	return k
}

func (k *Keyboard) AddFilmFind() *Keyboard {
	return k.AddButton("🔎", "findFilmByTitle", states.CallbackFilmsFind, "")
}

func (k *Keyboard) AddFilmNew() *Keyboard {
	return k.AddButton("➕", "createFilm", states.CallbackFilmsNew, "")
}

func (k *Keyboard) AddFilmFiltersAndSorting(session *models.Session) *Keyboard {
	filtersEnable := session.GetFilmsFiltersByContext().IsFiltersEnabled()
	sortingEnable := session.GetFilmsSortingByContext().IsSortingEnabled()

	return k.AddButtonsWithRowSize(2,
		Button{utils.BoolToEmoji(sortingEnable), "sorting", states.CallbackFilmsSorting, ""},
		Button{utils.BoolToEmoji(filtersEnable), "filters", states.CallbackFilmsFilters, ""},
	)
}

func (k *Keyboard) AddFilmDelete() *Keyboard {
	return k.AddButton("🗑️", "deleteFilm", states.CallbackManageFilmSelectDelete, "")
}

func (k *Keyboard) AddFilmUpdate() *Keyboard {
	return k.AddButton("✏️", "updateFilm", states.CallbackManageFilmSelectUpdate, "")
}

func (k *Keyboard) AddFilmManage() *Keyboard {
	return k.AddButton("⚙️", "manageFilm", states.CallbackFilmsManage, "")
}

func (k *Keyboard) AddFilmRemoveFromCollection() *Keyboard {
	return k.AddButton("❌", "removeFilmFromCollection", states.CallbackManageFilmSelectRemoveFromCollection, "")
}

func (k *Keyboard) AddNewFilmManually() *Keyboard {
	return k.AddButton("", "manually", states.CallbackNewFilmSelectManually, "")
}

func (k *Keyboard) AddNewFilmFromURL() *Keyboard {
	return k.AddButton("", "fromURL", states.CallbackNewFilmSelectFromURL, "")
}

func (k *Keyboard) AddAgain(callback string) *Keyboard {
	return k.AddButton("↻", "again", callback, "")
}

func (k *Keyboard) AddResetAllFilters() *Keyboard {
	return k.AddButton("", "resetFilters", states.CallbackFiltersFilmsSelectAllReset, "")
}

func (k *Keyboard) AddResetFilter(session *models.Session, filterType string) *Keyboard {
	filter := session.GetFilmsFiltersByContext()

	if filter.IsFilterEnabled(filterType) {
		return k.AddButton("", "reset", states.CallbackProcessReset, "")
	}

	return k
}

func (k *Keyboard) AddResetAllSorting() *Keyboard {
	return k.AddButton("", "resetSorting", states.CallbackSortingFilmsSelectAllReset, "")
}

func (k *Keyboard) AddResetSorting(session *models.Session) *Keyboard {
	sorting := session.GetFilmsSortingByContext()

	if sorting.IsSortingFieldEnabled(sorting.Field) {
		return k.AddButton("", "reset", states.CallbackProcessReset, "")
	}

	return k
}

func (k *Keyboard) AddSortingDirection() *Keyboard {
	return k.AddButtonsWithRowSize(2,
		Button{"⬆️", "increaseOrder", states.CallbackIncrease, ""},
		Button{"⬇️", "decreaseOrder", states.CallbacktDecrease, ""},
	)
}

func parseFiltersFilmsButtons(session *models.Session) []Button {
	var buttons []Button
	filters := session.GetFilmsFiltersByContext()

	filterEnabled := filters.IsFilterEnabled("minRating")
	text := translator.Translate(session.Lang, "minRating", nil, nil)
	if filterEnabled {
		text += fmt.Sprintf(": %.2f", filters.MinRating)
	}
	buttons = append(buttons, Button{
		utils.BoolToEmoji(filterEnabled),
		text,
		states.CallbackFiltersFilmsSelectMinRating,
		"",
	})

	filterEnabled = filters.IsFilterEnabled("maxRating")
	text = translator.Translate(session.Lang, "maxRating", nil, nil)
	if filterEnabled {
		text += fmt.Sprintf(": %.2f", filters.MaxRating)
	}
	buttons = append(buttons, Button{
		utils.BoolToEmoji(filterEnabled),
		text,
		states.CallbackFiltersFilmsSelectMaxRating,
		"",
	})

	return buttons
}

func parseSortingFilmsButtons(session *models.Session) []Button {
	var buttons []Button
	sorting := session.GetFilmsSortingByContext()

	sortingEnabled := sorting.IsSortingFieldEnabled("id")
	text := translator.Translate(session.Lang, "id", nil, nil)
	if sortingEnabled {
		text += fmt.Sprintf(": %s", utils.SortDirectionToEmoji(sorting.Sort))
	}
	buttons = append(buttons, Button{utils.BoolToEmoji(sortingEnabled), text, states.CallbackSortingFilmsSelectID, ""})

	sortingEnabled = sorting.IsSortingFieldEnabled("title")
	text = translator.Translate(session.Lang, "title", nil, nil)
	if sortingEnabled {
		text += fmt.Sprintf(": %s", utils.SortDirectionToEmoji(sorting.Sort))
	}
	buttons = append(buttons, Button{utils.BoolToEmoji(sortingEnabled), text, states.CallbackSortingFilmsSelectTitle, ""})

	sortingEnabled = sorting.IsSortingFieldEnabled("rating")
	text = translator.Translate(session.Lang, "rating", nil, nil)
	if sortingEnabled {
		text += fmt.Sprintf(": %s", utils.SortDirectionToEmoji(sorting.Sort))
	}
	buttons = append(buttons, Button{utils.BoolToEmoji(sortingEnabled), text, states.CallbackSortingFilmsSelectRating, ""})

	return buttons
}
