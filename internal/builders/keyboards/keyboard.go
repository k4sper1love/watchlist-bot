package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

type Button struct {
	Emoji         string
	Text          string
	Callback      string
	URL           string
	NeedTranslate bool
}

type Keyboard struct {
	Rows [][]Button
}

func New() *Keyboard {
	return &Keyboard{}
}

func (k *Keyboard) Build(lang string) *tgbotapi.InlineKeyboardMarkup {
	if lang != "" {
		k.translate(lang)
	}

	inlineKeyboard := k.createKeyboard()
	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}
}

func (k *Keyboard) translate(languageCode string) *Keyboard {
	for i, row := range k.Rows {
		for j, btn := range row {
			if btn.NeedTranslate {
				k.Rows[i][j].Text = translator.Translate(languageCode, btn.Text, nil, nil)
			}
		}
	}
	return k
}

func (k *Keyboard) createKeyboard() [][]tgbotapi.InlineKeyboardButton {
	var inlineKeyboard [][]tgbotapi.InlineKeyboardButton

	for _, row := range k.Rows {
		inlineRow := make([]tgbotapi.InlineKeyboardButton, 0, len(row))
		for _, btn := range row {
			inlineRow = append(inlineRow, k.createButton(btn))
		}
		inlineKeyboard = append(inlineKeyboard, inlineRow)
	}
	return inlineKeyboard
}

func (k *Keyboard) createButton(btn Button) tgbotapi.InlineKeyboardButton {
	if btn.Emoji != "" {
		btn.Text = fmt.Sprintf("%s %s", btn.Emoji, btn.Text)
	}

	if btn.URL != "" {
		return tgbotapi.NewInlineKeyboardButtonURL(btn.Text, btn.URL)
	}
	return tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.Callback)
}

func (k *Keyboard) AddRow(buttons ...Button) *Keyboard {
	k.Rows = append(k.Rows, buttons)
	return k
}

func (k *Keyboard) AddButton(emoji, text, callback, url string, needTranslate bool) *Keyboard {
	return k.AddRow(Button{Emoji: emoji, Text: text, Callback: callback, URL: url, NeedTranslate: needTranslate})
}

func (k *Keyboard) AddButtons(buttons ...Button) *Keyboard {
	for _, button := range buttons {
		k.AddButton(button.Emoji, button.Text, button.Callback, button.URL, button.NeedTranslate)
	}
	return k
}

func (k *Keyboard) AddButtonsWithRowSize(rowSize int, buttons ...Button) *Keyboard {
	for i := 0; i < len(buttons); i += rowSize {
		end := i + rowSize
		if end > len(buttons) {
			end = len(buttons)
		}
		k.AddRow(buttons[i:end]...)
	}
	return k
}
