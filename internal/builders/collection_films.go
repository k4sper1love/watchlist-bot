package builders

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func BuildCollectionFilmsMessage(collectionFilmsResponse *models.CollectionFilmsResponse) string {
	collection := collectionFilmsResponse.CollectionFilms.Collection
	films := collectionFilmsResponse.CollectionFilms.Films
	metadata := collectionFilmsResponse.Metadata

	msg := fmt.Sprintf("<b>🎬 Коллекция фильмов:</b> \"%s\"\n\n", collection.Name)

	if collection.TotalFilms == 0 {
		msg += "Не найдено фильмов в этой коллекции."
		return msg
	}

	msg += filmsToString(films, metadata)

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
	k.Buttons = append(k.Buttons, Button{"Удалить фильм", states.CallbackManageCollectionFilmSelectDelete})

	return k
}

func (k *Keyboard) AddCollectionFilmsUpdate() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Обновить фильм", states.CallbackManageCollectionFilmSelectUpdate})

	return k
}

func (k *Keyboard) AddCollectionFilmsManage() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Управление фильмом", states.CallbackCollectionFilmsManage})

	return k
}
