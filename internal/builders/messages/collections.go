package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildCollectionsMessage(session *models.Session, metadata *filters.Metadata, isFind bool) string {
	collections := session.CollectionsState.Collections

	msg := ""

	if metadata.TotalRecords == 0 {
		msg += translator.Translate(session.Lang, "collectionsNotFound", nil, nil)
		return msg
	}

	totalCollectionsMsgKey := "totalCollections"
	if isFind {
		totalCollectionsMsgKey = "totalCollectionsFilms"
	}

	totalCollectionsMsg := translator.Translate(session.Lang, totalCollectionsMsgKey, nil, nil)
	msg += fmt.Sprintf("📚 <b>%s:</b> %d\n\n", totalCollectionsMsg, metadata.TotalRecords)

	for i, collection := range collections {
		itemID := utils.GetItemID(i, metadata.CurrentPage, metadata.PageSize)

		numberEmoji := utils.NumberToEmoji(itemID)

		msg += fmt.Sprintf("%s\n", numberEmoji)
		msg += BuildCollectionDetailMessage(session, &collection)
	}

	pageMsg := translator.Translate(session.Lang, "pageCounter", map[string]interface{}{
		"CurrentPage": metadata.CurrentPage,
		"LastPage":    metadata.LastPage,
	}, nil)

	msg += fmt.Sprintf("<b>📄 %s</b>\n\n", pageMsg)

	msg += translator.Translate(session.Lang, "choiceCollectionForDetails", nil, nil)

	return msg
}
