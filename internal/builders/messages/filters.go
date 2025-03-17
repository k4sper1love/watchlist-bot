package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

// RequestSortDirection generates a message prompting the user to choose a sorting direction for a specific field.
func RequestSortDirection(session *models.Session, sorting *models.Sorting) string {
	return fmt.Sprintf("🗂️ %s\n\n%s",
		translator.Translate(session.Lang, "selectedSortField", map[string]interface{}{
			"Field": translator.Translate(session.Lang, sorting.Field, nil, nil),
		}, nil),
		translator.Translate(session.Lang, "requestDirection", nil, nil))
}

// ChoiceSorting generates a message prompting the user to choose a sorting option.
func ChoiceSorting(session *models.Session) string {
	return toBold(translator.Translate(session.Lang, "choiceSorting", nil, nil))
}

// ResetSortingSuccess generates a success message after resetting sorting settings.
func ResetSortingSuccess(session *models.Session) string {
	return "🔄 " + translator.Translate(session.Lang, "sortingResetSuccess", nil, nil)
}

// SortingApplied generates a message indicating that sorting has been applied successfully.
func SortingApplied(session *models.Session, sorting *models.Sorting) string {
	return fmt.Sprintf("%s %s",
		utils.SortDirectionToEmoji(sorting.Direction),
		translator.Translate(session.Lang, "sortingApplied", map[string]interface{}{
			"Field": translator.Translate(session.Lang, sorting.Field, nil, nil),
		}, nil))
}

// ResetFiltersSuccess generates a success message after resetting all filters.
func ResetFiltersSuccess(session *models.Session) string {
	return "🔄 " + translator.Translate(session.Lang, "filterResetSuccess", nil, nil)
}

// ResetFilterSuccess generates a success message after resetting a specific filter.
func ResetFilterSuccess(session *models.Session, filterType string) string {
	return "🔄 " + translator.Translate(session.Lang, "filterResetSuccess", map[string]interface{}{
		"Filter": translator.Translate(session.Lang, filterType, nil, nil),
	}, 1)
}

// FilterApplied generates a message indicating that a specific filter has been applied successfully.
func FilterApplied(session *models.Session, filterType, emoji string) string {
	return fmt.Sprintf("%s %s",
		emoji,
		translator.Translate(session.Lang, "filterApplied", map[string]interface{}{
			"Filter": translator.Translate(session.Lang, filterType, nil, nil),
		}, nil))
}

// InvalidFilterRange generates an error message when the user provides an invalid range for a filter.
// Includes examples and range limits for clarification.
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
