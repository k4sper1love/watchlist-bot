package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"strings"
)

func Collections(session *models.Session, metadata *filters.Metadata, isFind bool) string {
	if metadata.TotalRecords == 0 {
		return "❗️" + translator.Translate(session.Lang, "collectionsNotFound", nil, nil)
	}

	var msg strings.Builder
	totalCollectionsKey := "totalCollections"
	if isFind {
		totalCollectionsKey = "totalCollectionsFilms"
	}

	msg.WriteString(fmt.Sprintf("📚 %s: %d\n\n",
		toBold(translator.Translate(session.Lang, totalCollectionsKey, nil, nil)),
		metadata.TotalRecords))

	for i, collection := range session.CollectionsState.Collections {
		msg.WriteString(formatCollection(metadata, &collection, i))
	}

	msg.WriteString(formatPageCounter(session, metadata.CurrentPage, metadata.LastPage))
	return msg.String()
}

func CollectionsFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "getCollectionsFailure", nil, nil)
}

func RequestCollectionName(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "collectionRequestName", nil, nil)
}

func RequestCollectionDescription(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "collectionRequestDescription", nil, nil)
}

func DeleteCollection(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "deleteCollectionConfirm", map[string]interface{}{
		"Collection": session.CollectionDetailState.Collection.Name,
	}, nil)
}

func DeleteCollectionFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "deleteCollectionFailure", map[string]interface{}{
		"Collection": session.CollectionDetailState.Collection.Name,
	}, nil)
}

func DeleteCollectionSuccess(session *models.Session) string {
	return "🗑️ " + translator.Translate(session.Lang, "deleteCollectionSuccess", map[string]interface{}{
		"Collection": session.CollectionDetailState.Collection.Name,
	}, nil)
}

func CollectionChoiceAction(session *models.Session) string {
	return fmt.Sprintf("%s%s",
		CollectionHeader(session),
		toBold(translator.Translate(session.Lang, "choiceAction", nil, nil)))
}

func CreateCollectionFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "createCollectionFailure", nil, nil)
}

func CreateCollectionSuccess(session *models.Session) string {
	return "📚 " + translator.Translate(session.Lang, "createCollectionSuccess", nil, nil)
}

func UpdateCollection(session *models.Session) string {
	return fmt.Sprintf("%s%s",
		CollectionHeader(session),
		toBold(translator.Translate(session.Lang, "choiceAction", nil, nil)))
}

func UpdateCollectionFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "updateCollectionFailure", nil, nil)
}

func UpdateCollectionSuccess(session *models.Session) string {
	return "✏️ " + translator.Translate(session.Lang, "updateCollectionSuccess", nil, nil)
}

func CollectionsNotFound(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "collectionsNotFound", nil, nil)
}

func ChoiceCollection(session *models.Session) string {
	return toBold(translator.Translate(session.Lang, "choiceCollection", nil, nil))
}

func formatCollection(metadata *filters.Metadata, collection *apiModels.Collection, index int) string {
	return fmt.Sprintf("%s%s %s\n%s",
		utils.NumberToEmoji(utils.GetItemID(index, metadata.CurrentPage, metadata.PageSize)),
		formatOptionalBool("⭐", collection.IsFavorite, "%s"),
		toItalic(fmt.Sprintf("ID: %d", collection.ID)),
		CollectionDetail(collection))
}
