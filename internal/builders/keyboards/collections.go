package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

var updateCollectionButtons = []Button{
	{"", "title", states.CallbackUpdateCollectionSelectName, "", true},
	{"", "description", states.CallbackUpdateCollectionSelectDescription, "", true},
}

func Collections(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(len(session.CollectionsState.Collections) > 0, func(k *Keyboard) {
			k.AddSearch(states.CallbackCollectionsFind)
		}).
		AddCollectionsSelect(session).
		AddNavigation(currentPage, lastPage, states.PrefixCollectionsPage, true).
		AddCollectionFiltersAndSorting(session).
		AddCollectionsNew().
		AddBack("").
		Build(session.Lang)
}

func CollectionManage(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddUpdate(states.CallbackManageCollectionSelectUpdate).
		AddDelete(states.CallbackManageCollectionSelectDelete).
		AddBack(states.CallbackManageCollectionSelectBack).
		Build(session.Lang)
}

func CollectionUpdate(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddButtons(updateCollectionButtons...).
		AddBack(states.CallbackUpdateCollectionSelectBack).
		Build(session.Lang)
}

func FindCollections(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddCollectionsSelect(session).
		AddNavigation(currentPage, lastPage, states.PrefixFindCollectionsPage, true).
		AddBack(states.CallbackFindCollectionsBack).
		Build(session.Lang)
}

func CollectionsSorting(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	sorting := session.CollectionsState.Sorting
	return New().
		AddButtons(getSortingCollectionsButtons(sorting, session.Lang)...).
		AddIf(sorting.IsSortingEnabled(), func(k *Keyboard) {
			k.AddResetAllSorting(states.CallbackSortingCollectionsAllReset)
		}).
		AddBack(states.CallbackSortingCollectionsBack).
		Build(session.Lang)
}
