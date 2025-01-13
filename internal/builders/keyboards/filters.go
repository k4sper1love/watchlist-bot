package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildSortingDirectionKeyboard(session *models.Session, sorting *models.Sorting) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddSortingDirection()

	keyboard.AddResetSorting(sorting)

	keyboard.AddCancel()

	return keyboard.Build(session.Lang)
}

func (k *Keyboard) AddSearch(callback string) *Keyboard {
	return k.AddButton("🔎", "search", callback, "", true)
}

func (k *Keyboard) AddResetAllSorting(callback string) *Keyboard {
	return k.AddButton("", "resetSorting", callback, "", true)
}

func (k *Keyboard) AddResetSorting(sorting *models.Sorting) *Keyboard {
	if sorting.IsSortingFieldEnabled(sorting.Field) {
		return k.AddButton("", "reset", states.CallbackProcessReset, "", true)
	}

	return k
}

func (k *Keyboard) AddSortingDirection() *Keyboard {
	return k.AddButtonsWithRowSize(2,
		Button{"⬇️", "decreaseOrder", states.CallbacktDecrease, "", true},
		Button{"⬆️", "increaseOrder", states.CallbackIncrease, "", true},
	)
}

func addSortingButton(buttons []Button, sorting *models.Sorting, lang, field, callback string) []Button {
	sortingEnabled := sorting.IsSortingFieldEnabled(field)

	text := translator.Translate(lang, field, nil, nil)

	if sortingEnabled {
		text += fmt.Sprintf(": %s", utils.SortDirectionToEmoji(sorting.Sort))
	}

	button := Button{
		utils.BoolToEmoji(sortingEnabled),
		text,
		callback,
		"",
		false,
	}

	return append(buttons, button)
}

func addFiltersFilmsButton(buttons []Button, filter *models.FiltersFilm, lang, filterType, callback string, needTranslate bool) []Button {
	filterEnabled := filter.IsFilterEnabled(filterType)

	text := translator.Translate(lang, filterType, nil, nil)

	if filterEnabled {
		value := filter.ValueToString(filterType)
		if needTranslate {
			value = translator.Translate(lang, value, nil, nil)
		}
		text += fmt.Sprintf(": %s", value)
	}

	button := Button{
		utils.BoolToEmoji(filterEnabled),
		text,
		callback,
		"",
		false,
	}

	return append(buttons, button)
}
