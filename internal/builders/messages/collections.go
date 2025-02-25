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
		msg += "❗️" + translator.Translate(session.Lang, "collectionsNotFound", nil, nil)
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

		msg += utils.NumberToEmoji(itemID)

		if collection.IsFavorite {
			msg += "⭐"
		}

		msg += fmt.Sprintf(" <i>ID: %d</i>", collection.ID)

		msg += "\n" + BuildCollectionDetailMessage(session, &collection)
	}

	pageMsg := translator.Translate(session.Lang, "pageCounter", map[string]interface{}{
		"CurrentPage": metadata.CurrentPage,
		"LastPage":    metadata.LastPage,
	}, nil)

	msg += fmt.Sprintf("<b>📄 %s</b>", pageMsg)

	return msg
}

func BuildCollectionsFailureMessage(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "getCollectionsFailure", nil, nil)
}

func BuildCollectionRequestNameMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "collectionRequestName", nil, nil)
}

func BuildCollectionRequestDescriptionMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "collectionRequestDescription", nil, nil)
}

func BuildDeleteCollectionMessage(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "deleteCollectionConfirm", map[string]interface{}{
		"Collection": session.CollectionDetailState.Collection.Name,
	}, nil)
}

func BuildDeleteCollectionFailureMessage(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "deleteCollectionFailure", map[string]interface{}{
		"Collection": session.CollectionDetailState.Collection.Name,
	}, nil)
}

func BuildDeleteCollectionSuccessMessage(session *models.Session) string {
	return "🗑️ " + translator.Translate(session.Lang, "deleteCollectionSuccess", map[string]interface{}{
		"Collection": session.CollectionDetailState.Collection.Name,
	}, nil)
}

func BuildManageCollectionMessage(session *models.Session) string {
	msg := BuildCollectionHeader(session)
	choiceMsg := translator.Translate(session.Lang, "choiceAction", nil, nil)
	msg += fmt.Sprintf("<b>%s</b>", choiceMsg)
	return msg
}

func BuildCreateCollectionFailureMessage(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "createCollectionFailure", nil, nil)
}

func BuildCreateCollectionSuccessMessage(session *models.Session) string {
	return "📚 " + translator.Translate(session.Lang, "createCollectionSuccess", nil, nil)
}

func BuildUpdateCollectionMessage(session *models.Session) string {
	msg := BuildCollectionHeader(session)
	choiceMsg := translator.Translate(session.Lang, "choiceAction", nil, nil)
	msg += fmt.Sprintf("<b>%s</b>", choiceMsg)
	return msg
}

func BuildUpdateCollectionFailureMessage(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "updateCollectionFailure", nil, nil)
}

func BuildUpdateCollectionSuccessMessage(session *models.Session) string {
	return "✏️ " + translator.Translate(session.Lang, "updateCollectionSuccess", nil, nil)
}

func BuildCollectionsNotFoundMessage(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "collectionsNotFound", nil, nil)
}

func BuildChoiceCollectionMessage(session *models.Session) string {
	choiceMsg := translator.Translate(session.Lang, "choiceCollection", nil, nil)
	return fmt.Sprintf("<b>%s</b>", choiceMsg)
}

func BuildOptionsFilmToCollectionMessage(session *models.Session) string {
	msg := BuildCollectionHeader(session)
	msg += "<b>" + translator.Translate(session.Lang, "choiceAction", nil, nil) + "</b>"
	return msg
}
