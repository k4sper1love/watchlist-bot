package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/parsing"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildFilmsMessage(session *models.Session, metadata *filters.Metadata) string {
	var header string
	switch session.Context {
	case states.ContextFilm:
		header = ""
	case states.ContextCollection:
		header = BuildCollectionHeader(session)
	default:
		return translator.Translate(session.Lang, "unknownContext", nil, nil)
	}

	return header + buildFilmsList(session, metadata, false, true)
}

func BuildFindFilmsMessage(session *models.Session, metadata *filters.Metadata) string {
	return buildFilmsList(session, metadata, true, true)
}

func BuildFindNewFilmMessage(session *models.Session, metadata *filters.Metadata) string {
	return buildFilmsList(session, metadata, true, false)
}

func buildFilmsList(session *models.Session, metadata *filters.Metadata, isFind bool, needViewed bool) string {
	films := session.FilmsState.Films

	if metadata.TotalRecords == 0 {
		msg := "❗️" + translator.Translate(session.Lang, "filmsNotFound", nil, nil)
		return msg
	}

	totalFilmsMsgKey := "totalFilms"
	if isFind {
		totalFilmsMsgKey = "totalFindFilms"
	}

	totalFilmsMsg := translator.Translate(session.Lang, totalFilmsMsgKey, nil, nil)
	msg := fmt.Sprintf("🎥 <b>%s:</b> %d\n\n", totalFilmsMsg, metadata.TotalRecords)

	for i, film := range films {
		itemID := utils.GetItemID(i, metadata.CurrentPage, metadata.PageSize)
		numberEmoji := utils.NumberToEmoji(itemID)
		msg += fmt.Sprintf("%s ", numberEmoji)

		if film.IsFavorite {
			msg += "⭐ "
		}

		msg += fmt.Sprintf(" <i>ID: %d</i>", film.ID)

		msg += "\n" + BuildFilmGeneralMessage(session, &film, needViewed)
	}

	pageMsg := translator.Translate(session.Lang, "pageCounter", map[string]interface{}{
		"CurrentPage": metadata.CurrentPage,
		"LastPage":    metadata.LastPage,
	}, nil)
	msg += fmt.Sprintf("<b>📄 %s</b>", pageMsg)

	return msg
}

func BuildCollectionHeader(session *models.Session) string {
	collection := session.CollectionDetailState.Collection

	msg := fmt.Sprintf("<b>%s</b>", collection.Name)

	if collection.IsFavorite {
		msg += " ⭐"
	}

	if collection.Description != "" {
		msg += fmt.Sprintf("\n<i>%s</i>", collection.Description)
	}

	msg += "\n\n"

	return msg
}

func BuildFilterRangeMessage(session *models.Session, filterType string) string {
	filter := session.GetFilmsFiltersByContext()

	part1 := translator.Translate(session.Lang, "filterInstructionRange", nil, nil)
	part2 := translator.Translate(session.Lang, "filterInstructionPartialRange", nil, nil)

	msg := fmt.Sprintf("↕️ %s\n\n<i>%s</i>", part1, part2)

	if filter.IsFilterEnabled(filterType) {
		currentValueMsg := translator.Translate(session.Lang, "currentValue", nil, nil)
		value := filter.ValueToString(filterType)
		msg += fmt.Sprintf("\n\n<b>%s</b>: %s", currentValueMsg, value)
	}

	return msg
}

func BuildFilterSwitchMessage(session *models.Session, filterType string) string {
	filter := session.GetFilmsFiltersByContext()

	filterMsg := translator.Translate(session.Lang, filterType, nil, nil)
	msg := "🔀 " + translator.Translate(session.Lang, "filterInstructionSwitch", map[string]interface{}{
		"Filter": filterMsg,
	}, nil)

	if filter.IsFilterEnabled(filterType) {
		currentValueMsg := translator.Translate(session.Lang, "currentValue", nil, nil)
		value := translator.Translate(session.Lang, filter.ValueToString(filterType), nil, nil)
		msg += fmt.Sprintf("\n\n<b>%s</b>: %s", currentValueMsg, value)
	}

	return msg
}

func BuildFilmsFailureMessage(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "getFilmsFailure", nil, nil)
}

func BuildFilmRequestTitleMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestTitle", nil, nil)
}

func BuildDeleteFilmMessage(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "deleteFilmConfirm", map[string]interface{}{
		"Film": session.FilmDetailState.Film.Title,
	}, nil)
}

func BuildDeleteFilmFailureMessage(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "deleteFilmFailure", map[string]interface{}{
		"Film": session.FilmDetailState.Film.Title,
	}, nil)
}

func BuildDeleteFilmSuccessMessage(session *models.Session) string {
	return "🗑 " + translator.Translate(session.Lang, "deleteFilmSuccess", map[string]interface{}{
		"Film": session.FilmDetailState.Film.Title,
	}, nil)
}

func BuildFilmsNotFoundMessage(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "filmsNotFound", nil, nil)
}

func BuildManageFilmMessage(session *models.Session) string {
	return fmt.Sprintf("%s<b>%s</b>", BuildFilmDetailMessage(session), translator.Translate(session.Lang, "choiceAction", nil, nil))
}

func BuildRemoveFilmFailureMessage(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "removeFilmFailure", nil, nil)
}

func BuildRemoveFilmSuccessMessage(session *models.Session) string {
	return "🧹󠁝 " + translator.Translate(session.Lang, "removeFilmSuccess", nil, nil)
}

func BuildNewFilmFromURLMessage(session *models.Session) string {
	part1 := translator.Translate(session.Lang, "filmRequestLink", nil, nil)
	part2 := translator.Translate(session.Lang, "supportedServices", nil, nil)
	supportedServices := parsing.GetSupportedServicesInline()

	return fmt.Sprintf("❓<b>%s</b>\n\n%s:\n<i>%s</i>", part1, part2, supportedServices)
}

func BuildFilmRequestYearMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestYear", nil, nil)
}

func BuildFilmRequestGenreMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestGenre", nil, nil)
}

func BuildFilmRequestDescriptionMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestDescription", nil, nil)
}

func BuildFilmRequestRatingMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestRating", nil, nil)
}

func BuildFilmRequestImageMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestImage", nil, nil)
}

func BuildFilmRequestURLMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestLink", nil, nil)
}

func BuildFilmRequestCommentMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestComment", nil, nil)
}

func BuildFilmRequestViewedMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestViewed", nil, nil)
}

func BuildFilmRequestUserRatingMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestUserRating", nil, nil)
}

func BuildFilmRequestReviewMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestReview", nil, nil)
}

func BuildCreateFilmFailureMessage(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "createFilmFailure", nil, nil)
}

func BuildCreateFilmSuccessMessage(session *models.Session) string {
	return "🎬 " + translator.Translate(session.Lang, "createFilmSuccess", nil, nil)
}

func BuildCreateCollectionFilmSuccessMessage(session *models.Session, collectionName string) string {
	return "🎬 " + translator.Translate(session.Lang, "createCollectionFilmSuccess", map[string]interface{}{
		"Collection": collectionName,
	}, nil)
}

func BuildUpdateFilmMessage(session *models.Session) string {
	msg := BuildFilmDetailMessage(session)
	choiceMsg := translator.Translate(session.Lang, "updateChoiceField", nil, nil)
	msg += fmt.Sprintf("<b>%s</b>", choiceMsg)

	return msg
}

func BuildUpdateFilmFailureMessage(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "updateFilmFailure", nil, nil)
}

func BuildUpdateFilmSuccessMessage(session *models.Session) string {
	return "✏️ " + translator.Translate(session.Lang, "updateFilmSuccess", nil, nil)
}

func BuildViewedFilmRequestUserRatingMessage(session *models.Session) string {
	part1 := translator.Translate(session.Lang, "viewedFilmRequestRating", nil, nil)
	part2 := translator.Translate(session.Lang, "viewedFilmCanCancel", nil, nil)

	return fmt.Sprintf("✔️ <b>%s</b>\n\n<i>%s</i>", part1, part2)
}

func BuildViewedFilmRequestReviewMessage(session *models.Session) string {
	part1 := translator.Translate(session.Lang, "viewedFilmRequestReview", nil, nil)
	part2 := translator.Translate(session.Lang, "viewedFilmCanCancel", nil, nil)

	return fmt.Sprintf("✔️ <b>%s</b>\n\n<i>%s</i>", part1, part2)
}

func BuildFiltersFilmsMessage(session *models.Session) string {
	choiceMsg := translator.Translate(session.Lang, "choiceFilter", nil, nil)
	return fmt.Sprintf("<b>%s</b>", choiceMsg)
}

func BuildChoiceFilmMessage(session *models.Session) string {
	choiceMsg := translator.Translate(session.Lang, "choiceFilm", nil, nil)
	return fmt.Sprintf("<b>%s</b>", choiceMsg)
}

func BuildFilmToCollectionSuccessMessage(session *models.Session, collectionFilm *apiModels.CollectionFilm) string {
	return "➕ " + translator.Translate(session.Lang, "filmToCollectionSuccess", map[string]interface{}{
		"Film":       collectionFilm.Film.Title,
		"Collection": collectionFilm.Collection.Name,
	}, nil)
}
