package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildProfileMessage(session *models.Session) string {
	user := session.User

	part1 := translator.Translate(session.Lang, "profile", nil, nil)
	part2 := translator.Translate(session.Lang, "name", nil, nil)
	part3 := translator.Translate(session.Lang, "id", nil, nil)
	part4 := translator.Translate(session.Lang, "email", nil, nil)
	part5 := translator.Translate(session.Lang, "joinDate", nil, nil)

	msg := fmt.Sprintf(
		"👤 <b>%s:</b>\n\n"+
			"🔹 <b>%s:</b> <code>%s</code>\n"+
			"🔹 <b>%s:</b> <code>%d</code>\n"+
			"🔹 <b>%s:</b> <code>%s</code>\n"+
			"🔹 <b>%s:</b> %s\n\n",
		part1,
		part2, user.Username,
		part3, user.ID,
		part4, user.Email,
		part5, user.CreatedAt.Format("02.01.2006 15:04"))

	return msg
}
