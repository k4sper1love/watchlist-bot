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
	"strings"
)

// Films generates a message listing films based on the current session context (e.g., general list or collection).
func Films(session *models.Session, metadata *filters.Metadata) string {
	switch session.Context {
	case states.CtxFilm:
		return FilmList(session, metadata, false, true)
	case states.CtxCollection:
		return CollectionHeader(session) + FilmList(session, metadata, false, true)
	default:
		return translator.Translate(session.Lang, "unknownContext", nil, nil)
	}
}

// FindFilms generates a message listing films for the "find" operation with pagination details.
func FindFilms(session *models.Session, metadata *filters.Metadata) string {
	return FilmList(session, metadata, true, true)
}

// FindNewFilm generates a message listing new films for the "find new film" operation with pagination details.
func FindNewFilm(session *models.Session, metadata *filters.Metadata) string {
	return FilmList(session, metadata, true, false)
}

// FilmList generates a paginated list of films with details like title, ID, and viewed status.
func FilmList(session *models.Session, metadata *filters.Metadata, isFind bool, needViewed bool) string {
	if metadata.TotalRecords == 0 {
		return "❗️" + translator.Translate(session.Lang, "filmsNotFound", nil, nil)
	}

	var msg strings.Builder
	totalFilmsKey := "totalFilms"
	if isFind {
		totalFilmsKey = "totalFindFilms"
	}

	msg.WriteString(fmt.Sprintf("🎥 %s: %d\n\n",
		toBold(translator.Translate(session.Lang, totalFilmsKey, nil, nil)),
		metadata.TotalRecords))

	for i, film := range session.FilmsState.Films {
		msg.WriteString(formatFilm(session, needViewed, metadata, &film, i))
	}

	msg.WriteString(formatPageCounter(session, metadata.CurrentPage, metadata.LastPage))
	return msg.String()
}

// CollectionHeader generates a header for a collection, including its name, favorite status, and description.
func CollectionHeader(session *models.Session) string {
	return fmt.Sprintf("%s%s%s\n\n",
		toBold(session.CollectionDetailState.Collection.Name),
		formatOptionalBool("⭐", session.CollectionDetailState.Collection.IsFavorite, " %s"),
		formatOptionalString("", toItalic(session.CollectionDetailState.Collection.Description), "\n%s%s"))
}

// FilterRange generates a message for configuring a range-based filter (e.g., year, rating).
func FilterRange(session *models.Session, filterType string) string {
	filterEnabled := session.GetFilmFiltersByCtx().IsFieldEnabled(filterType)
	return fmt.Sprintf("↕️ %s\n\n%s%s%s",
		translator.Translate(session.Lang, "filterInstructionRange", nil, nil),
		translator.Translate(session.Lang, "filterInstructionPartialRange", nil, nil),
		formatOptionalBool(toBold(translator.Translate(session.Lang, "currentValue", nil, nil)),
			filterEnabled, "\n\n%s:"),
		formatOptionalBool(session.GetFilmFiltersByCtx().String(filterType),
			filterEnabled, " %s"))
}

// FilterSwitch generates a message for configuring a switch-based filter (e.g., genre, viewed status).
func FilterSwitch(session *models.Session, filterType string) string {
	filterEnabled := session.GetFilmFiltersByCtx().IsFieldEnabled(filterType)
	return fmt.Sprintf("🔀 %s%s%s",
		translator.Translate(session.Lang, "filterInstructionSwitch", map[string]interface{}{
			"Filter": translator.Translate(session.Lang, filterType, nil, nil),
		}, nil),
		formatOptionalBool(toBold(translator.Translate(session.Lang, "currentValue", nil, nil)),
			filterEnabled, "\n\n%s:"),
		formatOptionalBool(translator.Translate(session.Lang, session.GetFilmFiltersByCtx().String(filterType), nil, nil),
			filterEnabled, " %s"))
}

// FilmsFailure generates an error message when retrieving films fails.
func FilmsFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "getFilmsFailure", nil, nil)
}

// RequestFilmTitle generates a message prompting the user to enter a film title.
func RequestFilmTitle(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestTitle", nil, nil)
}

// DeleteFilm generates a confirmation message for deleting a film.
func DeleteFilm(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "deleteFilmConfirm", map[string]interface{}{
		"Film": session.FilmDetailState.Film.Title,
	}, nil)
}

// DeleteFilmFailure generates an error message when deleting a film fails.
func DeleteFilmFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "deleteFilmFailure", map[string]interface{}{
		"Film": session.FilmDetailState.Film.Title,
	}, nil)
}

// DeleteFilmSuccess generates a success message after deleting a film.
func DeleteFilmSuccess(session *models.Session) string {
	return "🗑 " + translator.Translate(session.Lang, "deleteFilmSuccess", map[string]interface{}{
		"Film": session.FilmDetailState.Film.Title,
	}, nil)
}

// FilmsNotFound generates a message indicating that no films were found.
func FilmsNotFound(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "filmsNotFound", nil, nil)
}

// ManageFilm generates a message prompting the user to choose an action for managing a specific film.
func ManageFilm(session *models.Session) string {
	return fmt.Sprintf("%s%s",
		FilmDetail(session),
		toBold(translator.Translate(session.Lang, "choiceAction", nil, nil)))
}

// RemoveFilmFailure generates an error message when removing a film from a collection fails.
func RemoveFilmFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "removeFilmFailure", nil, nil)
}

// RemoveFilmSuccess generates a success message after removing a film from a collection.
func RemoveFilmSuccess(session *models.Session) string {
	return "🧹󠁝 " + translator.Translate(session.Lang, "removeFilmSuccess", nil, nil)
}

// NewFilmFromURL generates a message prompting the user to enter a URL for creating a new film.
func NewFilmFromURL(session *models.Session) string {
	return fmt.Sprintf("❓%s\n\n%s:\n%s",
		toBold(translator.Translate(session.Lang, "filmRequestLink", nil, nil)),
		translator.Translate(session.Lang, "supportedServices", nil, nil),
		toItalic(parsing.GetSupportedServicesInline()))
}

// RequestFilmYear generates a message prompting the user to enter a film's release year.
func RequestFilmYear(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestYear", nil, nil)
}

// RequestFilmGenre generates a message prompting the user to enter a film's genre.
func RequestFilmGenre(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestGenre", nil, nil)
}

// RequestFilmDescription generates a message prompting the user to enter a film's description.
func RequestFilmDescription(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestDescription", nil, nil)
}

// RequestFilmRating generates a message prompting the user to enter a film's rating.
func RequestFilmRating(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestRating", nil, nil)
}

// RequestFilmImage generates a message prompting the user to upload a film's image.
func RequestFilmImage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestImage", nil, nil)
}

// RequestFilmURL generates a message prompting the user to enter a film's URL.
func RequestFilmURL(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestLink", nil, nil)
}

// RequestFilmComment generates a message prompting the user to enter a comment for a film.
func RequestFilmComment(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestComment", nil, nil)
}

// RequestFilmViewed generates a message prompting the user to confirm whether they have viewed a film.
func RequestFilmViewed(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestViewed", nil, nil)
}

// RequestFilmUserRating generates a message prompting the user to enter their personal rating for a film.
func RequestFilmUserRating(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestUserRating", nil, nil)
}

// RequestFilmReview generates a message prompting the user to enter a review for a film.
func RequestFilmReview(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestReview", nil, nil)
}

// CreateFilmFailure generates an error message when creating a film fails.
func CreateFilmFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "createFilmFailure", nil, nil)
}

// CreateFilmSuccess generates a success message after creating a film.
func CreateFilmSuccess(session *models.Session) string {
	return "🎬 " + translator.Translate(session.Lang, "createFilmSuccess", nil, nil)
}

// CreateCollectionFilmSuccess generates a success message after adding a film to a collection.
func CreateCollectionFilmSuccess(session *models.Session, collectionName string) string {
	return "🎬 " + translator.Translate(session.Lang, "createCollectionFilmSuccess", map[string]interface{}{
		"Collection": collectionName,
	}, nil)
}

// UpdateFilm generates a message prompting the user to choose a field to update for a specific film.
func UpdateFilm(session *models.Session) string {
	return fmt.Sprintf("%s%s",
		FilmDetail(session),
		toBold(translator.Translate(session.Lang, "updateChoiceField", nil, nil)))
}

// UpdateFilmFailure generates an error message when updating a film fails.
func UpdateFilmFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "updateFilmFailure", nil, nil)
}

// UpdateFilmSuccess generates a success message after updating a film.
func UpdateFilmSuccess(session *models.Session) string {
	return "✏️ " + translator.Translate(session.Lang, "updateFilmSuccess", nil, nil)
}

// RequestViewedFilmUserRating generates a message prompting the user to enter their rating for a viewed film.
func RequestViewedFilmUserRating(session *models.Session) string {
	return fmt.Sprintf("✔️ %s\n\n%s",
		toBold(translator.Translate(session.Lang, "viewedFilmRequestRating", nil, nil)),
		toItalic(translator.Translate(session.Lang, "viewedFilmCanCancel", nil, nil)))
}

// RequestViewedFilmReview generates a message prompting the user to enter a review for a viewed film.
func RequestViewedFilmReview(session *models.Session) string {
	return fmt.Sprintf("✔️ %s\n\n%s",
		toBold(translator.Translate(session.Lang, "viewedFilmRequestReview", nil, nil)),
		toItalic(translator.Translate(session.Lang, "viewedFilmCanCancel", nil, nil)))
}

// ChoiceFilter generates a message prompting the user to choose a filter.
func ChoiceFilter(session *models.Session) string {
	return toBold(translator.Translate(session.Lang, "choiceFilter", nil, nil))
}

// ChoiceFilm generates a message prompting the user to choose a film.
func ChoiceFilm(session *models.Session) string {
	return toBold(translator.Translate(session.Lang, "choiceFilm", nil, nil))
}

// AddFilmToCollectionSuccess generates a success message after adding a film to a collection.
func AddFilmToCollectionSuccess(session *models.Session, collectionFilm *apiModels.CollectionFilm) string {
	return "➕ " + translator.Translate(session.Lang, "filmToCollectionSuccess", map[string]interface{}{
		"Film":       collectionFilm.Film.Title,
		"Collection": collectionFilm.Collection.Name,
	}, nil)
}

// formatFilm formats a single film entry with details like favorite status, ID, and general information.
func formatFilm(session *models.Session, needViewed bool, metadata *filters.Metadata, film *apiModels.Film, index int) string {
	return fmt.Sprintf("%s%s %s\n%s",
		utils.NumberToEmoji(utils.GetItemID(index, metadata.CurrentPage, metadata.PageSize)),
		formatOptionalBool("⭐", film.IsFavorite, "%s"),
		toItalic(fmt.Sprintf("ID: %d", film.ID)),
		FilmGeneral(session, film, needViewed))
}
