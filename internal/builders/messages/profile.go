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
			"🔹 <b>%s:</b> %s",
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

func BuildUpdateProfileMessage(session *models.Session) string {
	msg := translator.Translate(session.Lang, "choiceField", nil, nil)
	return fmt.Sprintf("%s\n\n<b>%s</b>", BuildProfileMessage(session), msg)
}

func BuildUpdateProfileUsernameMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "updateProfileUsername", nil, nil)
}

func BuildUpdateProfileEmailMessage(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "updateProfileEmail", nil, nil)
}

func BuildUpdateProfileFailureMessage(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "updateProfileFailure", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

func BuildUpdateProfileSuccessMessage(session *models.Session) string {
	return "✏️ " + translator.Translate(session.Lang, "updateProfileSuccess", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

func BuildDeleteProfileMessage(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "deleteProfileConfirm", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

func BuildDeleteProfileFailureMessage(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "deleteProfileFailure", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

func BuildDeleteProfileSuccessMessage(session *models.Session) string {
	return "🗑️ " + translator.Translate(session.Lang, "deleteProfileSuccess", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}
