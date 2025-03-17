package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

// SortingDirection creates an inline keyboard for managing sorting direction.
func SortingDirection(session *models.Session, sorting *models.Sorting) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddSortingDirection().
		AddResetSorting(sorting).
		AddCancel().
		Build(session.Lang)
}
