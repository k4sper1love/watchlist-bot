package builders

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func BuildCollectionsMessage(collectionsResponse *models.CollectionsResponse) string {
	collections := collectionsResponse.Collections
	metadata := collectionsResponse.Metadata

	msg := ""

	if metadata.TotalRecords == 0 {
		msg += "Не найдено коллекций."
		return msg
	}

	msg += fmt.Sprintf("<b>📊 Всего коллекций:</b> %d\n\n", metadata.TotalRecords)

	for i, collection := range collections {
		itemID := utils.GetItemID(i, metadata.CurrentPage, metadata.PageSize)

		msg += "<b>──────────────────────────</b>\n\n"
		msg += BuildCollectionDetailMessage(itemID, &collection)
	}

	msg += "<b>──────────────────────────</b>\n\n"
	msg += fmt.Sprintf("<b>📄 Страница %d из %d</b>\n\n", metadata.CurrentPage, metadata.LastPage)
	msg += "Выберите коллекцию из списка, чтобы узнать больше."

	return msg
}

func (k *Keyboard) AddCollectionsSelect(collectionsResponse *models.CollectionsResponse) *Keyboard {
	for _, collection := range collectionsResponse.Collections {
		k.Buttons = append(k.Buttons, Button{collection.Name, fmt.Sprintf("select_collection_%d", collection.ID)})
	}
	return k
}

func (k *Keyboard) AddCollectionsNew() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Добавить новую коллекцию", states.CallbackCollectionsNew})
	return k
}

func (k *Keyboard) AddCollectionsDelete() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Удалить коллекцию", states.CallbackManageCollectionSelectDelete})
	return k
}

func (k *Keyboard) AddCollectionsUpdate() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Обновить коллекцию", states.CallbackManageCollectionSelectUpdate})
	return k
}

func (k *Keyboard) AddCollectionsManage() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Управление коллекцией", states.CallbackCollectionsManage})

	return k
}
