package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

type Button struct {
	Emoji        string
	Text         string
	CallbackData string
	URL          string
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

func (k *Keyboard) AddButton(emoji, text, callbackData, url string) *Keyboard {
	return k.AddRow(Button{Emoji: emoji, Text: text, CallbackData: callbackData, URL: url})
}

func (k *Keyboard) AddButtons(buttons ...Button) *Keyboard {
	for _, button := range buttons {
		k.AddButton(button.Emoji, button.Text, button.CallbackData, "")
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

func (k *Keyboard) Build(languageCode string) *tgbotapi.InlineKeyboardMarkup {
	if languageCode != "" {
		k.translate(languageCode)
	}

	var inlineButtons [][]tgbotapi.InlineKeyboardButton
	for _, row := range k.Rows {
		var inlineRow []tgbotapi.InlineKeyboardButton
		for _, btn := range row {
			fullText := btn.Text
			if btn.Emoji != "" {
				fullText = fmt.Sprintf("%s %s", btn.Emoji, btn.Text)
			}

			if btn.URL != "" {
				inlineRow = append(inlineRow, tgbotapi.NewInlineKeyboardButtonURL(fullText, btn.URL))
			} else if btn.CallbackData != "" {
				inlineRow = append(inlineRow, tgbotapi.NewInlineKeyboardButtonData(fullText, btn.CallbackData))
			}
		}

		if len(inlineRow) > 0 {
			inlineButtons = append(inlineButtons, inlineRow)
		}
	}

	keyboard := tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineButtons}

	return &keyboard
}

func (k *Keyboard) AddNavigation(currentPage, lastPage int, prevData, nextData string) *Keyboard {
	var buttons []Button

	if currentPage > 1 {
		buttons = append(buttons, Button{Emoji: "⬅", Text: "backward", CallbackData: prevData})
	}
	if currentPage < lastPage {
		buttons = append(buttons, Button{Emoji: "➡", Text: "forward", CallbackData: nextData})
	}

	if len(buttons) > 0 {
		k.AddRow(buttons...)
	}

	return k
}

func (k *Keyboard) AddURLButton(emoji, text, url string) *Keyboard {
	return k.AddButton(emoji, text, "", url)
}

func (k *Keyboard) AddCancel() *Keyboard {
	return k.AddButton("", "cancel", states.CallbackProcessCancel, "")
}

func (k *Keyboard) AddSkip() *Keyboard {
	return k.AddButton("", "skip", states.CallbackProcessSkip, "")
}

func (k *Keyboard) AddSurvey() *Keyboard {
	return k.AddButtonsWithRowSize(2,
		Button{"", "yes", states.CallbackYes, ""},
		Button{"", "no", states.CallbackNo, ""},
	)
}

func (k *Keyboard) AddBack(callbackData string) *Keyboard {
	var buttons []Button

	if callbackData != "" {
		buttons = append(buttons, Button{"←", "back", callbackData, ""})
	}

	buttons = append(buttons, Button{"🏠 ", "mainMenu", states.CallbackMainMenu, ""})

	return k.AddButtonsWithRowSize(len(buttons), buttons...)
}

func (k *Keyboard) translate(languageCode string) *Keyboard {
	for i, row := range k.Rows {
		for j, btn := range row {
			translatedText := translator.Translate(languageCode, btn.Text, nil, nil)
			k.Rows[i][j].Text = translatedText
		}
	}

	return k
}
