package messages

import (
	"fmt"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
)

func BuildCollectionDetailMessage(collection *apiModels.Collection) string {
	msg := fmt.Sprintf("<b>Коллекция:</b> %s\n", collection.Name)

	if collection.Description != "" {
		msg += fmt.Sprintf("<b>Описание:</b> %s\n", collection.Description)
	}

	msg += fmt.Sprintf("<b>Всего фильмов:</b> %d\n\n", collection.TotalFilms)
	return msg
}
