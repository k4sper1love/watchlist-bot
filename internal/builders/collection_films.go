package builders

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func BuildCollectionFilmsMessage(collectionFilmsResponse *models.CollectionFilmsResponse) string {
	msg := "Вот ваши фильмы:\n"

	for i, film := range collectionFilmsResponse.CollectionFilms.Films {
		itemID := i + 1 + ((collectionFilmsResponse.Metadata.CurrentPage - 1) * collectionFilmsResponse.Metadata.PageSize)

		msg += fmt.Sprintf("%d. ID: %d\nTitle: %s\nGenre: %s\nDescription: %s\nRating: %.2f\nLast updated: %s\nCreated: %s\n",
			itemID, film.ID, film.Title, film.Genre, film.Description, film.Rating, film.UpdatedAt, film.CreatedAt)
	}
	msg += fmt.Sprintf("%d из %d страниц\n", collectionFilmsResponse.Metadata.CurrentPage, collectionFilmsResponse.Metadata.LastPage)

	return msg
}

func BuildCollectionFilmsSelectButtons(collectionFilmsResponse *models.CollectionFilmsResponse) []Button {
	var buttons []Button

	for _, film := range collectionFilmsResponse.CollectionFilms.Films {
		buttons = append(buttons, Button{film.Title, fmt.Sprintf("select_collection_film_%d", film.ID)})
	}

	return buttons
}
