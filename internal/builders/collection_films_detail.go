package builders

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/models"
)

func BuildCollectionFilmsDetailMessage(collectionFilms models.CollectionFilms, index int) string {
	film := collectionFilms.Films[index]

	msg := "Вот ваш фильм:\n"
	msg += fmt.Sprintf("№: %d\nID: %d\nTitle: %s\nDescription: %s\nRating %0.2f\n",
		index+1, film.ID, film.Title, film.Description, film.Rating)

	return msg
}
