package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

var updateCollectionButtons = []Button{
	{"", "title", states.CallbackUpdateCollectionSelectName, "", true},
	{"", "description", states.CallbackUpdateCollectionSelectDescription, "", true},
}

func BuildCollectionsKeyboard(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	if len(session.CollectionsState.Collections) > 0 {
		keyboard.AddSearch(states.CallbackCollectionsFind)
	}

	keyboard.AddCollectionsSelect(session)

	keyboard.AddNavigation(
		currentPage,
		lastPage,
		states.CallbackCollectionsPrevPage,
		states.CallbackCollectionsNextPage,
		states.CallbackCollectionsFirstPage,
		states.CallbackCollectionsLastPage,
	)

	keyboard.AddCollectionFiltersAndSorting(session)

	keyboard.AddCollectionsNew()

	keyboard.AddBack("")

	return keyboard.Build(session.Lang)
}

func BuildCollectionUpdateKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddButtons(updateCollectionButtons...)

	keyboard.AddBack(states.CallbackUpdateCollectionSelectBack)

	return keyboard.Build(session.Lang)
}

func BuildFindCollectionsKeyboard(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddCollectionsSelect(session)

	keyboard.AddNavigation(
		currentPage,
		lastPage,
		states.CallbackFindCollectionsPrevPage,
		states.CallbackFindCollectionsNextPage,
		states.CallbackFindCollectionsFirstPage,
		states.CallbackFindCollectionsLastPage,
	)

	keyboard.AddBack(states.CallbackFindCollectionsBack)

	return keyboard.Build(session.Lang)
}

func BuildCollectionsSortingKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddButtons(parseSortingCollectionsButtons(session)...)

	keyboard.AddResetAllSorting(states.CallbackSortingCollectionsSelectAllReset)

	keyboard.AddBack(states.CallbackSortingCollectionsSelectBack)

	return keyboard.Build(session.Lang)
}

func (k *Keyboard) AddCollectionsSelect(session *models.Session) *Keyboard {
	var buttons []Button

	for i, collection := range session.CollectionsState.Collections {
		itemID := utils.GetItemID(i, session.CollectionsState.CurrentPage, session.CollectionsState.PageSize)

		buttons = append(buttons, Button{"", fmt.Sprintf("%s %s", utils.NumberToEmoji(itemID), collection.Name), fmt.Sprintf("select_collection_%d", collection.ID), "", false})
	}

	k.AddButtonsWithRowSize(2, buttons...)

	return k
}

func (k *Keyboard) AddCollectionsNew() *Keyboard {
	return k.AddButton("➕", "createCollection", states.CallbackCollectionsNew, "", true)
}

func (k *Keyboard) AddCollectionsDelete() *Keyboard {
	return k.AddButton("🗑️", "deleteCollection", states.CallbackManageCollectionSelectDelete, "", true)
}

func (k *Keyboard) AddCollectionsUpdate() *Keyboard {
	return k.AddButton("✏️", "updateCollection", states.CallbackManageCollectionSelectUpdate, "", true)
}

func (k *Keyboard) AddCollectionsManage() *Keyboard {
	return k.AddButton("⚙️", "manageCollection", states.CallbackCollectionsManage, "", true)
}

func (k *Keyboard) AddCollectionFiltersAndSorting(session *models.Session) *Keyboard {
	sortingEnable := session.CollectionsState.Sorting.IsSortingEnabled()

	return k.AddButtonsWithRowSize(2,
		Button{utils.BoolToEmoji(sortingEnable), "sorting", states.CallbackCollectionsSorting, "", true},
		//Button{utils.BoolToEmoji(filtersEnable), "filters", states.CallbackFilmsFilters, ""},
	)
}

func parseSortingCollectionsButtons(session *models.Session) []Button {
	var buttons []Button
	sorting := session.CollectionsState.Sorting

	sortingEnabled := sorting.IsSortingFieldEnabled("id")
	text := translator.Translate(session.Lang, "id", nil, nil)
	if sortingEnabled {
		text += fmt.Sprintf(": %s", utils.SortDirectionToEmoji(sorting.Sort))
	}
	buttons = append(buttons, Button{utils.BoolToEmoji(sortingEnabled), text, states.CallbackSortingCollectionsSelectID, "", true})

	sortingEnabled = sorting.IsSortingFieldEnabled("name")
	text = translator.Translate(session.Lang, "name", nil, nil)
	if sortingEnabled {
		text += fmt.Sprintf(": %s", utils.SortDirectionToEmoji(sorting.Sort))
	}
	buttons = append(buttons, Button{utils.BoolToEmoji(sortingEnabled), text, states.CallbackSortingCollectionsSelectName, "", true})

	sortingEnabled = sorting.IsSortingFieldEnabled("created_at")
	text = translator.Translate(session.Lang, "created_at", nil, nil)
	if sortingEnabled {
		text += fmt.Sprintf(": %s", utils.SortDirectionToEmoji(sorting.Sort))
	}
	buttons = append(buttons, Button{utils.BoolToEmoji(sortingEnabled), text, states.CallbackSortingCollectionsSelectCreatedAt, "", true})

	sortingEnabled = sorting.IsSortingFieldEnabled("total_films")
	text = translator.Translate(session.Lang, "total_films", nil, nil)
	if sortingEnabled {
		text += fmt.Sprintf(": %s", utils.SortDirectionToEmoji(sorting.Sort))
	}
	buttons = append(buttons, Button{utils.BoolToEmoji(sortingEnabled), text, states.CallbackSortingCollectionsSelectTotalFilms, "", true})

	return buttons
}
