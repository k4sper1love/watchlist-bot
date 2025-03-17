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

// Collections generates a message listing collections with pagination details.
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

// CollectionsFailure generates an error message when retrieving collections fails.
func CollectionsFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "getCollectionsFailure", nil, nil)
}

// RequestCollectionName generates a message prompting the user to enter a collection name.
func RequestCollectionName(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "collectionRequestName", nil, nil)
}

// RequestCollectionDescription generates a message prompting the user to enter a collection description.
func RequestCollectionDescription(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "collectionRequestDescription", nil, nil)
}

// DeleteCollection generates a confirmation message for deleting a collection.
func DeleteCollection(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "deleteCollectionConfirm", map[string]interface{}{
		"Collection": session.CollectionDetailState.Collection.Name,
	}, nil)
}

// DeleteCollectionFailure generates an error message when deleting a collection fails.
func DeleteCollectionFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "deleteCollectionFailure", map[string]interface{}{
		"Collection": session.CollectionDetailState.Collection.Name,
	}, nil)
}

// DeleteCollectionSuccess generates a success message after deleting a collection.
func DeleteCollectionSuccess(session *models.Session) string {
	return "🗑️ " + translator.Translate(session.Lang, "deleteCollectionSuccess", map[string]interface{}{
		"Collection": session.CollectionDetailState.Collection.Name,
	}, nil)
}

// CollectionChoiceAction generates a message prompting the user to choose an action for a specific collection.
func CollectionChoiceAction(session *models.Session) string {
	return fmt.Sprintf("%s%s",
		CollectionHeader(session),
		toBold(translator.Translate(session.Lang, "choiceAction", nil, nil)))
}

// CreateCollectionFailure generates an error message when creating a collection fails.
func CreateCollectionFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "createCollectionFailure", nil, nil)
}

// CreateCollectionSuccess generates a success message after creating a collection.
func CreateCollectionSuccess(session *models.Session) string {
	return "📚 " + translator.Translate(session.Lang, "createCollectionSuccess", nil, nil)
}

// UpdateCollection generates a message prompting the user to choose an action for updating a collection.
func UpdateCollection(session *models.Session) string {
	return fmt.Sprintf("%s%s",
		CollectionHeader(session),
		toBold(translator.Translate(session.Lang, "choiceAction", nil, nil)))
}

// UpdateCollectionFailure generates an error message when updating a collection fails.
func UpdateCollectionFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "updateCollectionFailure", nil, nil)
}

// UpdateCollectionSuccess generates a success message after updating a collection.
func UpdateCollectionSuccess(session *models.Session) string {
	return "✏️ " + translator.Translate(session.Lang, "updateCollectionSuccess", nil, nil)
}

// CollectionsNotFound generates a message indicating that no collections were found.
func CollectionsNotFound(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "collectionsNotFound", nil, nil)
}

// ChoiceCollection generates a message prompting the user to choose a collection.
func ChoiceCollection(session *models.Session) string {
	return toBold(translator.Translate(session.Lang, "choiceCollection", nil, nil))
}

// formatCollection formats a single collection entry with details like favorite status, ID, and name.
func formatCollection(metadata *filters.Metadata, collection *apiModels.Collection, index int) string {
	return fmt.Sprintf("%s%s %s\n%s",
		utils.NumberToEmoji(utils.GetItemID(index, metadata.CurrentPage, metadata.PageSize)),
		formatOptionalBool("⭐", collection.IsFavorite, "%s"),
		toItalic(fmt.Sprintf("ID: %d", collection.ID)),
		CollectionDetail(collection))
}
