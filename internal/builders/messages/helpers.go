package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

// wrapText wraps the given text with a specified format if the text is not empty.
func wrapText(text, format string) string {
	if text == "" {
		return ""
	}
	return fmt.Sprintf(format, text)
}

// toBold formats the given text as bold using HTML tags.
func toBold(text string) string {
	return wrapText(text, "<b>%s</b>")
}

// toCode formats the given text as inline code using HTML tags.
func toCode(text string) string {
	return wrapText(text, "<code>%s</code>")
}

// toItalic formats the given text as italic using HTML tags.
func toItalic(text string) string {
	return wrapText(text, "<i>%s</i>")
}

// toPre formats the given text as preformatted text using HTML tags.
func toPre(text string) string {
	return wrapText(text, "<pre>%s</pre>")
}

// formatPageCounter generates a pagination message displaying the current and last page numbers.
func formatPageCounter(session *models.Session, currentPage, lastPage int) string {
	return fmt.Sprintf("📄 %s",
		toBold(translator.Translate(session.Lang, "pageCounter", map[string]interface{}{
			"CurrentPage": currentPage,
			"LastPage":    lastPage,
		}, nil)))
}

// formatOptionalString formats a string value with a label if the value is not empty.
func formatOptionalString(label, value, format string) string {
	if value != "" {
		return fmt.Sprintf(format, label, value)
	}
	return ""
}

// formatOptionalNumber formats a numeric value with a label if the value is not equal to the zero value.
func formatOptionalNumber[T int | float64](label string, value, zeroValue T, format string) string {
	if value != zeroValue {
		return fmt.Sprintf(format, label, value)
	}
	return ""
}

// formatOptionalBool formats a boolean value with a label if the value is true.
func formatOptionalBool(label string, value bool, format string) string {
	if value {
		return fmt.Sprintf(format, label)
	}
	return ""
}

// nonEmpty returns the fallback value if the given value is empty; otherwise, returns the value itself.
func nonEmpty(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
