package builders

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func BuildCollectionsMessage(collectionsResponse *models.CollectionsResponse) string {
	msg := "Вот ваши коллекции:\n"

	for i, collection := range collectionsResponse.Collections {
		itemID := i + 1 + ((collectionsResponse.Metadata.CurrentPage - 1) * collectionsResponse.Metadata.PageSize)

		msg += fmt.Sprintf("%d. ID: %d\nName: %s\nDescription: %s\nLast updated: %s\nCreated: %s\n",
			itemID, collection.ID, collection.Name, collection.Description, collection.UpdatedAt, collection.CreatedAt)
	}

	msg += fmt.Sprintf("%d из %d страниц\n", collectionsResponse.Metadata.CurrentPage, collectionsResponse.Metadata.LastPage)

	return msg + "Выберите колллекцию из списка"
}

func BuildCollectionsSelectButtons(collectionsResponse *models.CollectionsResponse) []Button {
	var buttons []Button

	for _, collection := range collectionsResponse.Collections {
		buttons = append(buttons, Button{collection.Name, fmt.Sprintf("select_collection_%d", collection.ID)})
	}

	return buttons
}
