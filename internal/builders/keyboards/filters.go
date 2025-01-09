package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func BuildSortingDirectionKeyboard(session *models.Session, sorting *models.Sorting) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddSortingDirection()

	keyboard.AddResetSorting(sorting)

	keyboard.AddCancel()

	return keyboard.Build(session.Lang)
}

func (k *Keyboard) AddSearch(callback string) *Keyboard {
	return k.AddButton("🔎", "search", callback, "")
}

func (k *Keyboard) AddResetAllSorting(callback string) *Keyboard {
	return k.AddButton("", "resetSorting", callback, "")
}

func (k *Keyboard) AddResetSorting(sorting *models.Sorting) *Keyboard {
	if sorting.IsSortingFieldEnabled(sorting.Field) {
		return k.AddButton("", "reset", states.CallbackProcessReset, "")
	}

	return k
}

func (k *Keyboard) AddSortingDirection() *Keyboard {
	return k.AddButtonsWithRowSize(2,
		Button{"⬆️", "increaseOrder", states.CallbackIncrease, ""},
		Button{"⬇️", "decreaseOrder", states.CallbacktDecrease, ""},
	)
}
