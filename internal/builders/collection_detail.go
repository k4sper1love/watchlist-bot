package builders

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
)

func BuildCollectionDetailMessage(itemID int, collection *apiModels.Collection) string {
	msg := "<b>🔍︎ Коллекция </b>"

	if itemID != -1 {
		msg += fmt.Sprintf("<b>№%d</b>. ", itemID)
	}

	msg += fmt.Sprintf("<b>ID:</b> %d\n", collection.ID)
	msg += fmt.Sprintf("<b>📖 Название:</b> %s\n", collection.Name)

	if collection.Description != "" {
		msg += fmt.Sprintf("<b>📝 Описание:</b> %s\n", collection.Description)
	}

	msg += fmt.Sprintf("<b>🎬 Всего фильмов:</b> %d\n", collection.TotalFilms)

	msg += fmt.Sprintf("<b>🕒 Последнее обновление:</b> %s\n", collection.UpdatedAt.Format("02.01.2006 15:04"))
	msg += fmt.Sprintf("<b>📅 Создана:</b> %s\n\n", collection.CreatedAt.Format("02.01.2006 15:04"))

	return msg
}
