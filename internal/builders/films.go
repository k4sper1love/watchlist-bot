package builders

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func BuildFilmsMessage(filmsResponse *models.FilmsResponse) string {
	films := filmsResponse.Films
	metadata := filmsResponse.Metadata

	if metadata.TotalRecords == 0 {
		msg := "Не найдено фильмов в этой коллекции."
		return msg
	}

	return filmsToString(films, metadata)
}

func filmsToString(films []apiModels.Film, metadata filters.Metadata) string {
	msg := fmt.Sprintf("<b>📊 Всего фильмов:</b> %d\n\n", metadata.TotalRecords)

	for i, film := range films {
		itemID := utils.GetItemID(i, metadata.CurrentPage, metadata.PageSize)

		msg += "<b>──────────────────────────</b>\n\n"
		msg += BuildCollectionFilmGeneralMessage(itemID, &film)
	}

	msg += "<b>──────────────────────────</b>\n\n"
	msg += fmt.Sprintf("<b>📄 Страница %d из %d</b>\n", metadata.CurrentPage, metadata.LastPage)
	msg += "Выберите фильм из списка, чтобы узнать больше."

	return msg
}

func (k *Keyboard) AddFilmsSelect(filmsResponse *models.FilmsResponse) *Keyboard {
	for i, film := range filmsResponse.Films {
		k.Buttons = append(k.Buttons, Button{film.Title, fmt.Sprintf("select_film_%d", i)})
	}
	return k
}

func (k *Keyboard) AddFilmsNew() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Добавить новый фильм", states.CallbackFilmsNew})

	return k
}

func (k *Keyboard) AddFilmsDelete() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Удалить фильм", states.CallbackManageFilmSelectDelete})

	return k
}

func (k *Keyboard) AddFilmsUpdate() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Обновить фильм", states.CallbackManageFilmSelectUpdate})

	return k
}

func (k *Keyboard) AddFilmsManage() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Управление фильмом", states.CallbackFilmsManage})

	return k
}
