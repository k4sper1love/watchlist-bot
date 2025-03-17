package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

// Predefined buttons for updating collection details.
var updateCollectionButtons = []Button{
	{"", "title", states.CallUpdateCollectionName, "", true},
	{"", "description", states.CallUpdateCollectionDescription, "", true},
}

// Collections creates an inline keyboard for managing collections.
func Collections(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(len(session.CollectionsState.Collections) > 0, func(k *Keyboard) {
			k.AddSearch(states.CallCollectionsFind)
		}).
		AddCollectionsSelect(session).
		AddNavigation(currentPage, lastPage, states.CollectionsPage, true).
		AddCollectionFiltersAndSorting(session).
		AddCollectionsNew().
		AddBack("").
		Build(session.Lang)
}

// CollectionManage creates an inline keyboard for managing a specific collection.
func CollectionManage(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddUpdate(states.CallManageCollectionUpdate).
		AddDelete(states.CallManageCollectionDelete).
		AddBack(states.CallManageCollectionBack).
		Build(session.Lang)
}

// CollectionUpdate creates an inline keyboard for updating collection details (e.g., title, description).
func CollectionUpdate(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddButtons(updateCollectionButtons...).
		AddBack(states.CallUpdateCollectionBack).
		Build(session.Lang)
}

// FindCollections creates an inline keyboard for selecting collections with navigation and back options.
func FindCollections(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddCollectionsSelect(session).
		AddNavigation(currentPage, lastPage, states.FindCollectionsPage, true).
		AddBack(states.CallFindCollectionsBack).
		Build(session.Lang)
}

// CollectionsSorting creates an inline keyboard for managing collection sorting options.
func CollectionsSorting(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	sorting := session.CollectionsState.Sorting
	return New().
		AddButtons(getSortingCollectionsButtons(sorting, session.Lang)...).
		AddIf(sorting.IsEnabled(), func(k *Keyboard) {
			k.AddResetAllSorting(states.CallCollectionSortingAllReset)
		}).
		AddBack(states.CallCollectionSortingBack).
		Build(session.Lang)
}
