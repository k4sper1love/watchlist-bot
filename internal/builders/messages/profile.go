package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"strconv"
)

// Profile generates a detailed message about the user's profile.
func Profile(session *models.Session) string {
	return fmt.Sprintf("👤 %s\n\n🔹 %s: %s\n🔹 %s: %s\n🔹 %s: %s\n\n🔹 %s: %s\n🔹 %s: %s\n🔹 %s: %s\n🔹 %s: %s",
		toBold(translator.Translate(session.Lang, "profile", nil, nil)),
		toBold(translator.Translate(session.Lang, "role", nil, nil)),
		toCode(translator.Translate(session.Lang, session.Role.String(), nil, nil)),
		toBold(translator.Translate(session.Lang, "telegramID", nil, nil)),
		toCode(strconv.Itoa(session.TelegramID)),
		toBold(translator.Translate(session.Lang, "language", nil, nil)),
		toCode(nonEmpty(session.Lang, translator.Translate(session.Lang, "empty", nil, nil))),
		toBold(translator.Translate(session.Lang, "id", nil, nil)),
		toCode(strconv.Itoa(session.User.ID)),
		toBold(translator.Translate(session.Lang, "name", nil, nil)),
		toCode(session.User.Username),
		toBold(translator.Translate(session.Lang, "email", nil, nil)),
		toCode(nonEmpty(session.User.Email, translator.Translate(session.Lang, "empty", nil, nil))),
		toBold(translator.Translate(session.Lang, "created", nil, nil)),
		fmt.Sprintf("<code>%s</code>", session.User.CreatedAt.Format("02.01.2006 15:04")))
}

// UpdateProfile generates a message prompting the user to choose a field to update in their profile.
func UpdateProfile(session *models.Session) string {
	return fmt.Sprintf("%s\n\n%s",
		Profile(session),
		toBold(translator.Translate(session.Lang, "choiceField", nil, nil)))
}

// RequestProfileUsername generates a message prompting the user to enter a new username for their profile.
func RequestProfileUsername(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "updateProfileUsername", nil, nil)
}

// RequestProfileEmail generates a message prompting the user to enter a new email for their profile.
func RequestProfileEmail(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "updateProfileEmail", nil, nil)
}

// UpdateProfileFailure generates an error message when updating the user's profile fails.
func UpdateProfileFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "updateProfileFailure", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

// UpdateProfileSuccess generates a success message after updating the user's profile.
func UpdateProfileSuccess(session *models.Session) string {
	return "✏️ " + translator.Translate(session.Lang, "updateProfileSuccess", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

// DeleteProfile generates a confirmation message for deleting the user's profile.
func DeleteProfile(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "deleteProfileConfirm", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

// DeleteProfileFailure generates an error message when deleting the user's profile fails.
func DeleteProfileFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "deleteProfileFailure", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

// DeleteProfileSuccess generates a success message after deleting the user's profile.
func DeleteProfileSuccess(session *models.Session) string {
	return "🗑️ " + translator.Translate(session.Lang, "deleteProfileSuccess", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}
