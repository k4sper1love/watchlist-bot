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

func Films(session *models.Session, metadata *filters.Metadata) string {
	switch session.Context {
	case states.ContextFilm:
		return FilmList(session, metadata, false, true)
	case states.ContextCollection:
		return CollectionHeader(session) + FilmList(session, metadata, false, true)
	default:
		return translator.Translate(session.Lang, "unknownContext", nil, nil)
	}
}

func FindFilms(session *models.Session, metadata *filters.Metadata) string {
	return FilmList(session, metadata, true, true)
}

func FindNewFilm(session *models.Session, metadata *filters.Metadata) string {
	return FilmList(session, metadata, true, false)
}

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

func CollectionHeader(session *models.Session) string {
	return fmt.Sprintf("%s%s%s\n\n",
		toBold(session.CollectionDetailState.Collection.Name),
		formatOptionalBool("⭐", session.CollectionDetailState.Collection.IsFavorite, " %s"),
		formatOptionalString("", toItalic(session.CollectionDetailState.Collection.Description), "\n%s%s"))
}

func FilterRange(session *models.Session, filterType string) string {
	filterEnabled := session.GetFilmsFiltersByContext().IsFilterEnabled(filterType)
	return fmt.Sprintf("↕️ %s\n\n%s%s%s",
		translator.Translate(session.Lang, "filterInstructionRange", nil, nil),
		translator.Translate(session.Lang, "filterInstructionPartialRange", nil, nil),
		formatOptionalBool(toBold(translator.Translate(session.Lang, "currentValue", nil, nil)),
			filterEnabled, "\n\n%s:"),
		formatOptionalBool(session.GetFilmsFiltersByContext().ValueToString(filterType),
			filterEnabled, " %s"))
}

func FilterSwitch(session *models.Session, filterType string) string {
	filterEnabled := session.GetFilmsFiltersByContext().IsFilterEnabled(filterType)
	return fmt.Sprintf("🔀 %s%s%s",
		translator.Translate(session.Lang, "filterInstructionSwitch", map[string]interface{}{
			"Filter": translator.Translate(session.Lang, filterType, nil, nil),
		}, nil),
		formatOptionalBool(toBold(translator.Translate(session.Lang, "currentValue", nil, nil)),
			filterEnabled, "\n\n%s:"),
		formatOptionalBool(translator.Translate(session.Lang, session.GetFilmsFiltersByContext().ValueToString(filterType), nil, nil),
			filterEnabled, " %s"))
}

func FilmsFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "getFilmsFailure", nil, nil)
}

func RequestFilmTitle(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestTitle", nil, nil)
}

func DeleteFilm(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "deleteFilmConfirm", map[string]interface{}{
		"Film": session.FilmDetailState.Film.Title,
	}, nil)
}

func DeleteFilmFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "deleteFilmFailure", map[string]interface{}{
		"Film": session.FilmDetailState.Film.Title,
	}, nil)
}

func DeleteFilmSuccess(session *models.Session) string {
	return "🗑 " + translator.Translate(session.Lang, "deleteFilmSuccess", map[string]interface{}{
		"Film": session.FilmDetailState.Film.Title,
	}, nil)
}

func FilmsNotFound(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "filmsNotFound", nil, nil)
}

func ManageFilm(session *models.Session) string {
	return fmt.Sprintf("%s%s",
		FilmDetail(session),
		toBold(translator.Translate(session.Lang, "choiceAction", nil, nil)))
}

func RemoveFilmFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "removeFilmFailure", nil, nil)
}

func RemoveFilmSuccess(session *models.Session) string {
	return "🧹󠁝 " + translator.Translate(session.Lang, "removeFilmSuccess", nil, nil)
}

func NewFilmFromURL(session *models.Session) string {
	return fmt.Sprintf("❓%s\n\n%s:\n%s",
		toBold(translator.Translate(session.Lang, "filmRequestLink", nil, nil)),
		translator.Translate(session.Lang, "supportedServices", nil, nil),
		toItalic(parsing.GetSupportedServicesInline()))
}

func RequestFilmYear(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestYear", nil, nil)
}

func RequestFilmGenre(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestGenre", nil, nil)
}

func RequestFilmDescription(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestDescription", nil, nil)
}

func RequestFilmRating(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestRating", nil, nil)
}

func RequestFilmImage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestImage", nil, nil)
}

func RequestFilmURL(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestLink", nil, nil)
}

func RequestFilmComment(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestComment", nil, nil)
}

func RequestFilmViewed(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestViewed", nil, nil)
}

func RequestFilmUserRating(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestUserRating", nil, nil)
}

func RequestFilmReview(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "filmRequestReview", nil, nil)
}

func CreateFilmFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "createFilmFailure", nil, nil)
}

func CreateFilmSuccess(session *models.Session) string {
	return "🎬 " + translator.Translate(session.Lang, "createFilmSuccess", nil, nil)
}

func CreateCollectionFilmSuccess(session *models.Session, collectionName string) string {
	return "🎬 " + translator.Translate(session.Lang, "createCollectionFilmSuccess", map[string]interface{}{
		"Collection": collectionName,
	}, nil)
}

func UpdateFilm(session *models.Session) string {
	return fmt.Sprintf("%s%s",
		FilmDetail(session),
		toBold(translator.Translate(session.Lang, "updateChoiceField", nil, nil)))
}

func UpdateFilmFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "updateFilmFailure", nil, nil)
}

func UpdateFilmSuccess(session *models.Session) string {
	return "✏️ " + translator.Translate(session.Lang, "updateFilmSuccess", nil, nil)
}

func RequestViewedFilmUserRating(session *models.Session) string {
	return fmt.Sprintf("✔️ %s\n\n%s",
		toBold(translator.Translate(session.Lang, "viewedFilmRequestRating", nil, nil)),
		toItalic(translator.Translate(session.Lang, "viewedFilmCanCancel", nil, nil)))
}

func RequestViewedFilmReview(session *models.Session) string {
	return fmt.Sprintf("✔️ %s\n\n%s",
		toBold(translator.Translate(session.Lang, "viewedFilmRequestReview", nil, nil)),
		toItalic(translator.Translate(session.Lang, "viewedFilmCanCancel", nil, nil)))
}

func ChoiceFilter(session *models.Session) string {
	return toBold(translator.Translate(session.Lang, "choiceFilter", nil, nil))
}

func ChoiceFilm(session *models.Session) string {
	return toBold(translator.Translate(session.Lang, "choiceFilm", nil, nil))
}

func AddFilmToCollectionSuccess(session *models.Session, collectionFilm *apiModels.CollectionFilm) string {
	return "➕ " + translator.Translate(session.Lang, "filmToCollectionSuccess", map[string]interface{}{
		"Film":       collectionFilm.Film.Title,
		"Collection": collectionFilm.Collection.Name,
	}, nil)
}

func formatFilm(session *models.Session, needViewed bool, metadata *filters.Metadata, film *apiModels.Film, index int) string {
	return fmt.Sprintf("%s%s %s\n%s",
		utils.NumberToEmoji(utils.GetItemID(index, metadata.CurrentPage, metadata.PageSize)),
		formatOptionalBool("⭐", film.IsFavorite, "%s"),
		toItalic(fmt.Sprintf("ID: %d", film.ID)),
		FilmGeneral(session, film, needViewed))
}
