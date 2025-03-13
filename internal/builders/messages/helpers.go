package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func wrapText(text, format string) string {
	if text == "" {
		return ""
	}
	return fmt.Sprintf(format, text)
}

func toBold(text string) string {
	return wrapText(text, "<b>%s</b>")
}

func toCode(text string) string {
	return wrapText(text, "<code>%s</code>")
}

func toItalic(text string) string {
	return wrapText(text, "<i>%s</i>")
}

func toPre(text string) string {
	return wrapText(text, "<pre>%s</pre>")
}

func formatPageCounter(session *models.Session, currentPage, lastPage int) string {
	return fmt.Sprintf("📄 %s",
		toBold(translator.Translate(session.Lang, "pageCounter", map[string]interface{}{
			"CurrentPage": currentPage,
			"LastPage":    lastPage,
		}, nil)))
}

func formatOptionalString(label, value, format string) string {
	if value != "" {
		return fmt.Sprintf(format, label, value)
	}
	return ""
}

func formatOptionalNumber[T int | float64](label string, value, zeroValue T, format string) string {
	if value != zeroValue {
		return fmt.Sprintf(format, label, value)
	}
	return ""
}

func formatOptionalBool(label string, value bool, format string) string {
	if value {
		return fmt.Sprintf(format, label)
	}
	return ""
}

func nonEmpty(value, fallback string) string {
	if value == "" {
		return fallback
	}
	return value
}
