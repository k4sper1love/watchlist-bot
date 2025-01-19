package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildProfileMessage(session *models.Session) string {
	user := session.User

	part1 := translator.Translate(session.Lang, "profile", nil, nil)
	part2 := translator.Translate(session.Lang, "role", nil, nil)
	part3 := translator.Translate(session.Lang, "telegramID", nil, nil)
	part4 := translator.Translate(session.Lang, "language", nil, nil)
	part5 := translator.Translate(session.Lang, "id", nil, nil)
	part6 := translator.Translate(session.Lang, "name", nil, nil)
	part7 := translator.Translate(session.Lang, "email", nil, nil)
	part8 := translator.Translate(session.Lang, "created", nil, nil)

	role := translator.Translate(session.Lang, session.Role.String(), nil, nil)

	lang := session.Lang
	if lang == "" {
		lang = translator.Translate(session.Lang, "empty", nil, nil)
	}

	email := user.Email
	if email == "" {
		email = translator.Translate(session.Lang, "empty", nil, nil)
	}

	createdAt := fmt.Sprintf("<code>%s</code>", user.CreatedAt.Format("02.01.2006 15:04"))

	msg := fmt.Sprintf(
		"👤 <b>%s</b>\n\n"+
			"🔹 <b>%s:</b> <code>%s</code>\n"+
			"🔹 <b>%s:</b> <code>%d</code>\n"+
			"🔹 <b>%s:</b> <code>%s</code>\n\n"+
			"🔹 <b>%s:</b> <code>%d</code>\n"+
			"🔹 <b>%s:</b> <code>%s</code>\n"+
			"🔹 <b>%s:</b> <code>%s</code>\n"+
			"🔹 <b>%s:</b> %s\n\n",
		part1,
		part2, role,
		part3, session.TelegramID,
		part4, lang,
		part5, user.ID,
		part6, user.Username,
		part7, email,
		part8, createdAt)

	return msg
}
