package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildSelectedSortMessage(session *models.Session, sorting *models.Sorting) string {
	field := translator.Translate(session.Lang, sorting.Field, nil, nil)
	part1 := translator.Translate(session.Lang, "selectedSortField", map[string]interface{}{
		"Field": field,
	}, nil)
	part2 := translator.Translate(session.Lang, "requestDirection", nil, nil)

	msg := fmt.Sprintf("🗂️ %s\n\n%s", part1, part2)

	return msg
}

func BuildSortingMessage(session *models.Session) string {
	choiceMsg := translator.Translate(session.Lang, "choiceSorting", nil, nil)
	return fmt.Sprintf("<b>%s</b>", choiceMsg)
}

func BuildSortingResetSuccessMessage(session *models.Session) string {
	return "🔄 " + translator.Translate(session.Lang, "sortingResetSuccess", nil, nil)
}

func BuildSortingAppliedMessage(session *models.Session, sorting *models.Sorting) string {
	fieldMsg := translator.Translate(session.Lang, sorting.Field, nil, nil)
	directionEmoji := utils.SortDirectionToEmoji(sorting.Direction)
	return directionEmoji + " " + translator.Translate(session.Lang, "sortingApplied", map[string]interface{}{
		"Field": fieldMsg,
	}, nil)
}

func BuildFilterResetSuccessSimpleMessage(session *models.Session) string {
	return "🔄 " + translator.Translate(session.Lang, "filterResetSuccess", nil, nil)
}

func BuildFilterResetSuccessMessage(session *models.Session, filterType string) string {
	filterMsg := translator.Translate(session.Lang, filterType, nil, nil)
	return "🔄 " + translator.Translate(session.Lang, "filterResetSuccess", map[string]interface{}{
		"Filter": filterMsg,
	}, 1)
}

func BuildFilterAppliedMessage(session *models.Session, filterType, emoji string) string {
	filterMsg := translator.Translate(session.Lang, filterType, nil, nil)
	return emoji + " " + translator.Translate(session.Lang, "filterApplied", map[string]interface{}{
		"Filter": filterMsg,
	}, nil)
}

func BuildInvalidFilterRangeInputMessage(session *models.Session, config utils.FilterRangeConfig) string {
	exampleValue := translator.Translate(session.Lang, "exampleValue", nil, nil)
	exampleRange := translator.Translate(session.Lang, "exampleRange", nil, nil)
	examplePartialRange := translator.Translate(session.Lang, "examplePartialRange", nil, nil)
	rangeLimits := translator.Translate(session.Lang, "rangeLimits", map[string]interface{}{
		"Min": fmt.Sprintf("%.f", config.MinValue),
		"Max": fmt.Sprintf("%.f", config.MaxValue),
	}, nil)

	msg := "❌ " + translator.Translate(session.Lang, "invalidInput", nil, nil)
	msg += "\n\n<b>" + translator.Translate(session.Lang, "requestRangeInFormat", nil, nil) + "</b>"
	msg += fmt.Sprintf("\n- %s: <code>%s</code>", exampleValue, "5.5")
	msg += fmt.Sprintf("\n- %s: <code>%s</code>", exampleRange, "1990-2023")
	msg += fmt.Sprintf("\n- %s: <code>%s</code> или <code>%s</code>", examplePartialRange, "5-", "-10")
	msg += fmt.Sprintf("\n\n⚠️ <i>%s</i>", rangeLimits)

	return msg
}
