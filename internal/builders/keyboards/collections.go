package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func BuildCollectionsKeyboard(session *models.Session, currentPage, lastPage int) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddCollectionsSelect(session)

	keyboard.AddNavigation(
		currentPage,
		lastPage,
		states.CallbackCollectionsPrevPage,
		states.CallbackCollectionsNextPage,
	)

	keyboard.AddCollectionsNew()

	keyboard.AddBack("")

	return keyboard.Build(session.Lang)
}

func (k *Keyboard) AddCollectionsSelect(session *models.Session) *Keyboard {
	var buttons []Button

	for i, collection := range session.CollectionsState.Collections {
		itemID := utils.GetItemID(i, session.CollectionsState.CurrentPage, session.CollectionsState.PageSize)

		buttons = append(buttons, Button{"", fmt.Sprintf("%s (%d)", collection.Name, itemID), fmt.Sprintf("select_collection_%d", collection.ID), ""})
	}

	k.AddButtonsWithRowSize(2, buttons...)

	return k
}

func (k *Keyboard) AddCollectionsNew() *Keyboard {
	return k.AddButton("➕", "createCollection", states.CallbackCollectionsNew, "")
}

func (k *Keyboard) AddCollectionsDelete() *Keyboard {
	return k.AddButton("🗑️", "deleteCollection", states.CallbackManageCollectionSelectDelete, "")
}

func (k *Keyboard) AddCollectionsUpdate() *Keyboard {
	return k.AddButton("✏️", "updateCollection", states.CallbackManageCollectionSelectUpdate, "")
}

func (k *Keyboard) AddCollectionsManage() *Keyboard {
	return k.AddButton("⚙️", "manageCollection", states.CallbackCollectionsManage, "")
}
