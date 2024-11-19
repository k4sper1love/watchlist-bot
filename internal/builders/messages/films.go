package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func BuildFilmsMessage(session *models.Session, metadata *filters.Metadata) string {
	if session.Context == states.ContextFilm {
		return "🎥 " + filmsToString(session, metadata)
	} else if session.Context == states.ContextCollection {
		return collectionFilmsToString(session, metadata)
	}

	return "Неизвестный контекст"
}

func filmsToString(session *models.Session, metadata *filters.Metadata) string {
	films := session.FilmsState.Films

	msg := fmt.Sprintf("<b>Всего фильмов:</b> %d\n\n", metadata.TotalRecords)

	if metadata.TotalRecords == 0 {
		msg += "Не найдено фильмов."
		return msg
	}

	for i, film := range films {
		itemID := utils.GetItemID(i, metadata.CurrentPage, metadata.PageSize)

		numberEmoji := numberToEmoji(itemID)

		msg += fmt.Sprintf("%s\n", numberEmoji)
		msg += BuildFilmGeneralMessage(&film)
	}

	msg += fmt.Sprintf("<b>📄 Страница %d из %d</b>\n", metadata.CurrentPage, metadata.LastPage)

	msg += "Выберите фильм из списка, чтобы узнать больше."

	return msg
}

func collectionFilmsToString(session *models.Session, metadata *filters.Metadata) string {
	collection := session.CollectionDetailState.Collection

	msg := fmt.Sprintf("<b>Коллекция:</b> \"%s\"\n", collection.Name)

	if collection.Description != "" {
		msg += fmt.Sprintf("<b>Описание:</b> %s\n", collection.Description)
	} else {
		msg += "\n"
	}

	if collection.TotalFilms == 0 {
		msg += "Не найдено фильмов в этой коллекции."
		return msg
	}

	return msg + filmsToString(session, metadata)
}
