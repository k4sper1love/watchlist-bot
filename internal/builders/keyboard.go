package builders

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
)

type Button struct {
	Text         string
	CallbackData string
}

type Keyboard struct {
	Buttons []Button
	Size    int
}

func NewKeyboard(size int) *Keyboard {
	return &Keyboard{Size: size}
}

func (k *Keyboard) Build() *tgbotapi.InlineKeyboardMarkup {
	var inlineButtons [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for i, btn := range k.Buttons {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.CallbackData))

		if (i+1)%k.Size == 0 {
			inlineButtons = append(inlineButtons, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	if len(row) > 0 {
		inlineButtons = append(inlineButtons, row)
	}

	keyboard := tgbotapi.NewInlineKeyboardMarkup(inlineButtons...)

	return &keyboard
}

func (k *Keyboard) Add(button Button) *Keyboard {
	k.Buttons = append(k.Buttons, button)
	return k
}

func (k *Keyboard) AddSeveral(buttons []Button) *Keyboard {
	k.Buttons = append(k.Buttons, buttons...)
	return k
}

func (k *Keyboard) AddNavigation(currentPage, lastPage int, prevData, nextData string) *Keyboard {
	if currentPage > 1 {
		k.Buttons = append(k.Buttons, Button{"Предыдущая страница", prevData})
	}

	if currentPage < lastPage {
		k.Buttons = append(k.Buttons, Button{"Следующая страница", nextData})
	}

	return k
}

func (k *Keyboard) AddCancel() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Отмена", states.CallbackProcessCancel})
	return k
}

func (k *Keyboard) AddSkip() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Пропустить", states.CallbackProcessSkip})
	return k
}

func (k *Keyboard) AddSurvey() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Да", states.CallbackYes})

	k.Buttons = append(k.Buttons, Button{"Нет", states.CallbackNo})

	return k
}

func (k *Keyboard) AddBack(callbackData string) *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Назад", callbackData})
	return k
}
