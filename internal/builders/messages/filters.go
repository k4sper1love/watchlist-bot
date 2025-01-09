package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildSelectedSortMessage(session *models.Session, sorting *models.Sorting) string {
	field := translator.Translate(session.Lang, sorting.Field, nil, nil)
	part1 := translator.Translate(session.Lang, "selectedSortField", map[string]interface{}{
		"Field": field,
	}, nil)
	part2 := translator.Translate(session.Lang, "requestDirection", nil, nil)

	msg := fmt.Sprintf("%s\n\n%s", part1, part2)

	return msg
}
