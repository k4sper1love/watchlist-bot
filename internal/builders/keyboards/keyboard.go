package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

type Button struct {
	Emoji         string
	Text          string
	CallbackData  string
	URL           string
	NeedTranslate bool
}

type Keyboard struct {
	Rows [][]Button
}

func NewKeyboard() *Keyboard {
	return &Keyboard{}
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

func (k *Keyboard) AddRow(buttons ...Button) *Keyboard {
	k.Rows = append(k.Rows, buttons)
	return k
}

func (k *Keyboard) AddButton(emoji, text, callbackData, url string, needTranslate bool) *Keyboard {
	return k.AddRow(Button{Emoji: emoji, Text: text, CallbackData: callbackData, URL: url, NeedTranslate: needTranslate})
}

func (k *Keyboard) AddButtons(buttons ...Button) *Keyboard {
	for _, button := range buttons {
		k.AddButton(button.Emoji, button.Text, button.CallbackData, button.URL, button.NeedTranslate)
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

func (k *Keyboard) AddNavigation(currentPage, lastPage int, prevData, nextData, firstData, lastData string) *Keyboard {
	var buttons []Button

	if currentPage > 1 {
		if firstData != "" {
			buttons = append(buttons, Button{Emoji: "⏮", Text: "", CallbackData: firstData, NeedTranslate: false})
		}
		buttons = append(buttons, Button{Emoji: "⬅", Text: "", CallbackData: prevData, NeedTranslate: false})
	}

	if currentPage < lastPage {
		if nextData != "" {
			buttons = append(buttons, Button{Emoji: "➡", Text: "", CallbackData: nextData, NeedTranslate: false})
		}
		buttons = append(buttons, Button{Emoji: "⏭", Text: "", CallbackData: lastData, NeedTranslate: false})
	}

	if len(buttons) > 0 {
		k.AddRow(buttons...)
	}

	return k
}

func (k *Keyboard) AddURLButton(emoji, text, url string, needTranslate bool) *Keyboard {
	return k.AddButton(emoji, text, "", url, needTranslate)
}

func (k *Keyboard) AddCancel() *Keyboard {
	return k.AddButton("", "cancel", states.CallbackProcessCancel, "", true)
}

func (k *Keyboard) AddSkip() *Keyboard {
	return k.AddButton("", "skip", states.CallbackProcessSkip, "", true)
}

func (k *Keyboard) AddSurvey() *Keyboard {
	return k.AddButtonsWithRowSize(2,
		Button{"", "yes", states.CallbackYes, "", true},
		Button{"", "no", states.CallbackNo, "", true},
	)
}

func (k *Keyboard) AddBack(callbackData string) *Keyboard {
	var buttons []Button

	if callbackData != "" {
		buttons = append(buttons, Button{"←", "back", callbackData, "", true})
	}

	buttons = append(buttons, Button{"🏠 ", "mainMenu", states.CallbackMainMenu, "", true})

	return k.AddButtonsWithRowSize(len(buttons), buttons...)
}

func (k *Keyboard) AddUpdate(callbackData string) *Keyboard {
	return k.AddButton("✏️", "update", callbackData, "", true)
}

func (k *Keyboard) AddDelete(callbackData string) *Keyboard {
	return k.AddButton("🗑️", "delete", callbackData, "", true)
}

func (k *Keyboard) AddManage(callbackData string) *Keyboard {
	return k.AddButton("⚙️", "manage", callbackData, "", true)
}

func (k *Keyboard) translate(languageCode string) *Keyboard {
	for i, row := range k.Rows {
		for j, btn := range row {
			if btn.NeedTranslate {
				translatedText := translator.Translate(languageCode, btn.Text, nil, nil)
				k.Rows[i][j].Text = translatedText
			}
		}
	}

	return k
}
