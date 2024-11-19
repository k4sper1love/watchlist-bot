package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
)

type Button struct {
	Text         string
	CallbackData string
}

type Keyboard struct {
	Rows [][]Button
}

func NewKeyboard() *Keyboard {
	return &Keyboard{}
}

func (k *Keyboard) AddRow(buttons ...Button) *Keyboard {
	k.Rows = append(k.Rows, buttons)
	return k
}

func (k *Keyboard) AddButton(text, callbackData string) *Keyboard {
	return k.AddRow(Button{Text: text, CallbackData: callbackData})
}

func (k *Keyboard) AddButtons(buttons ...Button) *Keyboard {
	for _, button := range buttons {
		k.AddButton(button.Text, button.CallbackData)
	}

	return k
}

func (k *Keyboard) AddButtonsWithRowSize(rowSize int, buttons ...Button) *Keyboard {
	var row []Button
	for i, button := range buttons {
		row = append(row, button)
		if (i+1)%rowSize == 0 {
			k.AddRow(row...)
			row = []Button{}
		}
	}

	if len(row) > 0 {
		k.AddRow(row...)
	}

	return k
}

func (k *Keyboard) Build() *tgbotapi.InlineKeyboardMarkup {
	var inlineButtons [][]tgbotapi.InlineKeyboardButton

	for _, row := range k.Rows {
		var inlineRow []tgbotapi.InlineKeyboardButton
		for _, btn := range row {
			inlineRow = append(inlineRow, tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.CallbackData))
		}
		inlineButtons = append(inlineButtons, inlineRow)
	}

	keyboard := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineButtons}

	return &keyboard
}

func (k *Keyboard) AddNavigation(currentPage, lastPage int, prevData, nextData string) *Keyboard {
	var buttons []Button

	if currentPage > 1 {
		buttons = append(buttons, Button{Text: "⬅ Назад", CallbackData: prevData})
	}
	if currentPage < lastPage {
		buttons = append(buttons, Button{Text: "➡ Вперед", CallbackData: nextData})
	}

	if len(buttons) > 0 {
		k.AddRow(buttons...)
	}

	return k
}

func (k *Keyboard) AddCancel() *Keyboard {
	return k.AddButton("Отмена", states.CallbackProcessCancel)
}

func (k *Keyboard) AddSkip() *Keyboard {
	return k.AddButton("Пропустить", states.CallbackProcessSkip)
}

func (k *Keyboard) AddSurvey() *Keyboard {
	return k.AddButtonsWithRowSize(2,
		Button{"Да", states.CallbackYes},
		Button{"Нет", states.CallbackNo},
	)
}

func (k *Keyboard) AddBack(callbackData string) *Keyboard {
	var buttons []Button

	if callbackData != "" {
		buttons = append(buttons, Button{"← Обратно", callbackData})
	}

	buttons = append(buttons, Button{"🏠 Главное меню", states.CallbackMainMenu})

	return k.AddButtonsWithRowSize(len(buttons), buttons...)
}
