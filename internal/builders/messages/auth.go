package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func RegistrationSuccess(session *models.Session) string {
	return "✅ " + translator.Translate(session.Lang, "registrationSuccess", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

func Logout(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "logoutConfirm", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

func LogoutFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "logoutFailure", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

func LogoutSuccess(session *models.Session) string {
	return "🚪 " + translator.Translate(session.Lang, "logoutSuccess", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

func AuthFailure(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "authFailure", nil, nil)
}

func PermissionsNotEnough(session *models.Session) string {
	return "❌ " + translator.Translate(session.Lang, "permissionsNotEnough", nil, nil)
}

func Banned(session *models.Session) string {
	return fmt.Sprintf("❌ %s\n\n%s",
		translator.Translate(session.Lang, "bannedHeader", nil, nil),
		translator.Translate(session.Lang, "bannedBody", nil, nil))
}
