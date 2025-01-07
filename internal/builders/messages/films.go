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
	var header string
	switch session.Context {
	case states.ContextFilm:
		header = "🎥 "
	case states.ContextCollection:
		header = buildCollectionHeader(session)
	default:
		return translator.Translate(session.Lang, "unknownContext", nil, nil)
	}

	return header + buildFilmsList(session, metadata, false)
}

func BuildFindFilmsMessage(session *models.Session, metadata *filters.Metadata) string {
	return "🎥 " + buildFilmsList(session, metadata, true)
}

func buildFilmsList(session *models.Session, metadata *filters.Metadata, isFind bool) string {
	films := session.FilmsState.Films

	totalFilmsMsgKey := "totalFilms"
	if isFind {
		totalFilmsMsgKey = "totalFindFilms"
	}

	totalFilmsMsg := translator.Translate(session.Lang, totalFilmsMsgKey, nil, nil)
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

	// Информация о страницах
	pageMsg := translator.Translate(session.Lang, "pageCounter", map[string]interface{}{
		"CurrentPage": metadata.CurrentPage,
		"LastPage":    metadata.LastPage,
	}, nil)
	msg += fmt.Sprintf("<b>📄 %s</b>\n", pageMsg)

	msg += translator.Translate(session.Lang, "choiceFilmForDetails", nil, nil)
	return msg
}

func buildCollectionHeader(session *models.Session) string {
	collection := session.CollectionDetailState.Collection

	collectionMsg := translator.Translate(session.Lang, "collection", nil, nil)
	msg := fmt.Sprintf("<b>%s:</b> \"%s\"", collectionMsg, collection.Name)

	if collection.Description != "" {
		descriptionMsg := translator.Translate(session.Lang, "description", nil, nil)
		msg += fmt.Sprintf("\n<b>%s:</b> %s", descriptionMsg, collection.Description)
	}

	msg += "\n\n"

	if collection.TotalFilms == 0 {
		msg += translator.Translate(session.Lang, "notFoundFilmsInCollection", nil, nil)
	}
	return msg
}
