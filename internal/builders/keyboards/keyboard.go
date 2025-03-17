package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

// Button represents a single button in a keyboard.
// It includes text, emoji, callback data, URL, and translation flag.
type Button struct {
	Emoji         string // Emoji to display alongside the button text.
	Text          string // Text displayed on the button.
	Callback      string // Callback data sent when the button is pressed.
	URL           string // URL opened when the button is pressed (optional).
	NeedTranslate bool   // Indicates whether the button text needs translation.
}

// Keyboard represents a Telegram inline keyboard with rows of buttons.
type Keyboard struct {
	Rows [][]Button // Rows of buttons in the keyboard.
}

// New creates and returns a new empty Keyboard instance.
func New() *Keyboard {
	return &Keyboard{}
}

// Build constructs a Telegram inline keyboard markup from the current Keyboard.
// Translates button texts if a language code is provided.
func (k *Keyboard) Build(lang string) *tgbotapi.InlineKeyboardMarkup {
	if lang != "" {
		k.translate(lang)
	}

	inlineKeyboard := k.createKeyboard()
	return &tgbotapi.InlineKeyboardMarkup{InlineKeyboard: inlineKeyboard}
}

// translate translates button texts into the specified language.
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

// createKeyboard converts the Keyboard's rows of Buttons into a Telegram inline keyboard format.
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

// createButton converts a Button into a Telegram inline keyboard button.
func (k *Keyboard) createButton(btn Button) tgbotapi.InlineKeyboardButton {
	if btn.Emoji != "" {
		btn.Text = fmt.Sprintf("%s %s", btn.Emoji, btn.Text)
	}

	if btn.URL != "" {
		return tgbotapi.NewInlineKeyboardButtonURL(btn.Text, btn.URL)
	}
	return tgbotapi.NewInlineKeyboardButtonData(btn.Text, btn.Callback)
}

// AddRow adds a new row of buttons to the keyboard.
func (k *Keyboard) AddRow(buttons ...Button) *Keyboard {
	k.Rows = append(k.Rows, buttons)
	return k
}

// AddButton adds a single button to the keyboard as a new row.
func (k *Keyboard) AddButton(emoji, text, callback, url string, needTranslate bool) *Keyboard {
	return k.AddRow(Button{Emoji: emoji, Text: text, Callback: callback, URL: url, NeedTranslate: needTranslate})
}

// AddButtons adds multiple buttons to the keyboard, each as a separate row.
func (k *Keyboard) AddButtons(buttons ...Button) *Keyboard {
	for _, button := range buttons {
		k.AddButton(button.Emoji, button.Text, button.Callback, button.URL, button.NeedTranslate)
	}
	return k
}

// AddButtonsWithRowSize adds buttons to the keyboard, grouping them into rows of a specified size.
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
