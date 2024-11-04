package builders

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
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

func (k *Keyboard) AddCollectionFilmsSelect(collectionFilmsResponse *models.CollectionFilmsResponse) *Keyboard {
	for i, film := range collectionFilmsResponse.CollectionFilms.Films {
		k.Buttons = append(k.Buttons, Button{film.Title, fmt.Sprintf("select_cf_%d", i)})
	}
	return k
}

func (k *Keyboard) AddCollectionFilmsNew() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Добавить новый фильм", states.CallbackCollectionFilmsNew})

	return k
}

func (k *Keyboard) AddCollectionFilmsDelete() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Удалить фильм", states.CallbackCollectionFilmsDelete})

	return k
}

func (k *Keyboard) AddCollectionFilmsUpdate() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Обновить фильм", states.CallbackCollectionFilmsUpdate})

	return k
}
