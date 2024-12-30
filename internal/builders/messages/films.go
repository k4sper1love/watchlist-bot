package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildFilmsMessage(session *models.Session, metadata *filters.Metadata) string {
	if session.Context == states.ContextFilm {
		return "🎥 " + filmsToString(session, metadata)
	} else if session.Context == states.ContextCollection {
		return collectionFilmsToString(session, metadata)
	}

	msg := translator.Translate(session.Lang, "unknownContext", nil, nil)
	return msg
}

func filmsToString(session *models.Session, metadata *filters.Metadata) string {
	films := session.FilmsState.Films

	totalFilmsMsg := translator.Translate(session.Lang, "totalFilms", nil, nil)
	msg := fmt.Sprintf("<b>%s</b> %d\n\n", totalFilmsMsg, metadata.TotalRecords)

	if metadata.TotalRecords == 0 {
		msg += translator.Translate(session.Lang, "filmsNotFound", nil, nil)
		return msg
	}

	for i, film := range films {
		itemID := utils.GetItemID(i, metadata.CurrentPage, metadata.PageSize)

		numberEmoji := numberToEmoji(itemID)

		msg += fmt.Sprintf("%s\n", numberEmoji)
		msg += BuildFilmGeneralMessage(session, &film)
	}

	pageMsg := translator.Translate(session.Lang, "pageCounter", map[string]interface{}{
		"CurrentPage": metadata.CurrentPage,
		"LastPage":    metadata.LastPage,
	}, nil)

	msg += fmt.Sprintf("<b>📄 %s</b>\n", pageMsg)

	msg += translator.Translate(session.Lang, "choiceFilmForDetails", nil, nil)

	return msg
}

func collectionFilmsToString(session *models.Session, metadata *filters.Metadata) string {
	collection := session.CollectionDetailState.Collection

	collectionMsg := translator.Translate(session.Lang, "collection", nil, nil)
	msg := fmt.Sprintf("<b>%s:</b> \"%s\"\n", collectionMsg, collection.Name)

	if collection.Description != "" {
		descriptionMsg := translator.Translate(session.Lang, "description", nil, nil)
		msg += fmt.Sprintf("<b>%s:</b> %s\n", descriptionMsg, collection.Description)
	} else {
		msg += "\n"
	}

	if collection.TotalFilms == 0 {
		msg += translator.Translate(session.Lang, "notFoundFilmsInCollection", nil, nil)
		return msg
	}

	return msg + filmsToString(session, metadata)
}
