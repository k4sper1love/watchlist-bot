package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/models"
)

func BuildProfileMessage(user *models.User) string {
	msg := fmt.Sprintf(
		"👤 <b>Профиль:</b>\n\n"+
			"🔹 <b>Имя:</b> %s\n"+
			"🔹 <b>ID:</b> %d\n"+
			"🔹 <b>Email:</b> %s\n"+
			"🔹 <b>Дата регистрации:</b> %s\n\n",
		user.Username,
		user.ID,
		user.Email,
		user.CreatedAt.Format("02.01.2006 15:04"),
	)

	return msg
}
