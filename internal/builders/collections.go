package builders

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func BuildCollectionsMessage(collectionsResponse *models.CollectionsResponse) string {
	metadata := collectionsResponse.Metadata

	msg := "Вот ваши коллекции:\n"

	for i, collection := range collectionsResponse.Collections {
		itemID := utils.GetItemID(i, metadata.CurrentPage, metadata.PageSize)

		msg += fmt.Sprintf("%d. ID: %d\nName: %s\nDescription: %s\nLast updated: %s\nCreated: %s\n",
			itemID, collection.ID, collection.Name, collection.Description, collection.UpdatedAt, collection.CreatedAt)
	}

	msg += fmt.Sprintf("%d из %d страниц\n", collectionsResponse.Metadata.CurrentPage, collectionsResponse.Metadata.LastPage)

	return msg + "Выберите колллекцию из списка"
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
	k.Buttons = append(k.Buttons, Button{"Удалить коллекцию", states.CallbackCollectionsDelete})
	return k
}

func (k *Keyboard) AddCollectionsUpdate() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Обновить коллекцию", states.CallbackCollectionsUpdate})
	return k
}
