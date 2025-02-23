package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildRegistrationSuccessMessage(session *models.Session) string {
	return "✅ " + translator.Translate(session.Lang, "registrationSuccess", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

func BuildLogoutMessage(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "logoutConfirm", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

func BuildLogoutFailureMessage(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "logoutFailure", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

func BuildLogoutSuccessMessage(session *models.Session) string {
	return "🚪 " + translator.Translate(session.Lang, "logoutSuccess", map[string]interface{}{
		"Username": session.User.Username,
	}, nil)
}

func BuildAuthFailureMessage(session *models.Session) string {
	return "🚨 " + translator.Translate(session.Lang, "authFailure", nil, nil)
}

func BuildPermissionsNotEnoughMessage(session *models.Session) string {
	return "❌ " + translator.Translate(session.Lang, "permissionsNotEnough", nil, nil)
}

func BuildBannedMessage(session *models.Session) string {
	part1 := translator.Translate(session.Lang, "bannedHeader", nil, nil)
	part2 := translator.Translate(session.Lang, "bannedBody", nil, nil)

	return fmt.Sprintf("❌ %s\n\n%s", part1, part2)
}
