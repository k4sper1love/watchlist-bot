package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func RequestSortDirection(session *models.Session, sorting *models.Sorting) string {
	return fmt.Sprintf("🗂️ %s\n\n%s",
		translator.Translate(session.Lang, "selectedSortField", map[string]interface{}{
			"Field": translator.Translate(session.Lang, sorting.Field, nil, nil),
		}, nil),
		translator.Translate(session.Lang, "requestDirection", nil, nil))
}

func ChoiceSorting(session *models.Session) string {
	return toBold(translator.Translate(session.Lang, "choiceSorting", nil, nil))
}

func ResetSortingSuccess(session *models.Session) string {
	return "🔄 " + translator.Translate(session.Lang, "sortingResetSuccess", nil, nil)
}

func SortingApplied(session *models.Session, sorting *models.Sorting) string {
	return fmt.Sprintf("%s %s",
		utils.SortDirectionToEmoji(sorting.Direction),
		translator.Translate(session.Lang, "sortingApplied", map[string]interface{}{
			"Field": translator.Translate(session.Lang, sorting.Field, nil, nil),
		}, nil))
}

func ResetFiltersSuccess(session *models.Session) string {
	return "🔄 " + translator.Translate(session.Lang, "filterResetSuccess", nil, nil)
}

func ResetFilterSuccess(session *models.Session, filterType string) string {
	return "🔄 " + translator.Translate(session.Lang, "filterResetSuccess", map[string]interface{}{
		"Filter": translator.Translate(session.Lang, filterType, nil, nil),
	}, 1)
}

func FilterApplied(session *models.Session, filterType, emoji string) string {
	return fmt.Sprintf("%s %s",
		emoji,
		translator.Translate(session.Lang, "filterApplied", map[string]interface{}{
			"Filter": translator.Translate(session.Lang, filterType, nil, nil),
		}, nil))
}

func InvalidFilterRange(session *models.Session, config utils.FilterRangeConfig) string {
	return fmt.Sprintf("❌ %s\n\n%s\n- %s: %s\n- %s: %s\n- %s: %s, %s\n\n⚠️ %s",
		translator.Translate(session.Lang, "invalidInput", nil, nil),
		toBold(translator.Translate(session.Lang, "requestRangeInFormat", nil, nil)),
		translator.Translate(session.Lang, "exampleValue", nil, nil),
		toCode("5.5"),
		translator.Translate(session.Lang, "exampleRange", nil, nil),
		toCode("1990-2023"),
		translator.Translate(session.Lang, "examplePartialRange", nil, nil),
		toCode("5-"), toCode("-10"),
		toItalic(translator.Translate(session.Lang, "rangeLimits", map[string]interface{}{
			"Min": fmt.Sprintf("%.f", config.MinValue),
			"Max": fmt.Sprintf("%.f", config.MaxValue),
		}, nil)))
}
