package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

// RegistrationSuccess generates a success message after a user successfully registers.
func RegistrationSuccess(session *models.Session) string {
	return "✅ " + translator.Translate(session.Lang, "registrationSuccess", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

// Logout generates a confirmation message for logging out.
func Logout(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "logoutConfirm", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

// LogoutFailure generates an error message if the logout process fails.
func LogoutFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "logoutFailure", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

// LogoutSuccess generates a success message after the user logs out.
func LogoutSuccess(session *models.Session) string {
	return "🚪 " + translator.Translate(session.Lang, "logoutSuccess", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

// AuthFailure generates an error message if authentication fails.
func AuthFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "authFailure", nil, nil)
}

// PermissionsNotEnough generates a message indicating that the user does not have sufficient permissions.
func PermissionsNotEnough(session *models.Session) string {
	return "❌ " + translator.Translate(session.Lang, "permissionsNotEnough", nil, nil)
}

// Banned generates a message notifying the user that they have been banned.
func Banned(session *models.Session) string {
	return fmt.Sprintf("❌ %s\n\n%s",
		translator.Translate(session.Lang, "bannedHeader", nil, nil),
		translator.Translate(session.Lang, "bannedBody", nil, nil))
}
