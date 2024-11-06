package builders

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
)

func BuildCollectionDetailMessage(collection *apiModels.Collection) string {
	msg := fmt.Sprintf("🆔 <b>ID</b>: %d\n", collection.ID)

	if collection.Name != "" {
		msg += fmt.Sprintf("📛 <b>Название</b>: %s\n", collection.Name)
	}
	if collection.Description != "" {
		msg += fmt.Sprintf("📝 <b>Описание</b>: %s\n", collection.Description)
	}

	msg += fmt.Sprintf("📈 <b>Всего фильмов</b>: %d\n", collection.TotalFilms)

	return msg
}
