package builders

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

type Button struct {
	Text         string
	CallbackData string
}

func BuildButtonKeyboard(buttons []Button, buttonsPerRow int) tgbotapi.InlineKeyboardMarkup {
	var inlineButtons [][]tgbotapi.InlineKeyboardButton
	var row []tgbotapi.InlineKeyboardButton

	for i, btn := range buttons {
		row = append(row, tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.CallbackData))

		if (i+1)%buttonsPerRow == 0 {
			inlineButtons = append(inlineButtons, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	if len(row) > 0 {
		inlineButtons = append(inlineButtons, row)
	}

	return tgbotapi.NewInlineKeyboardMarkup(inlineButtons...)
}

func BuildNavigationButtons(currentPage, lastPage int, prevData, nextData string) []Button {
	var buttons []Button
	if currentPage > 1 {
		buttons = append(buttons, Button{"Предыдущая страница", prevData})
	}

	if currentPage < lastPage {
		buttons = append(buttons, Button{"Следующая страница", nextData})
	}

	return buttons
}
