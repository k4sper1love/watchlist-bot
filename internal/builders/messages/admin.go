package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"strconv"
	"strings"
)

// UserList generates a message listing users with pagination details.
func UserList(session *models.Session, users []models.Session) string {
	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("%s: %d\n\n",
		toBold(translator.Translate(session.Lang, "foundUsers", nil, nil)),
		session.AdminState.TotalRecords))

	for i, user := range users {
		msg.WriteString(formatUserEntry(session, user, i))
	}

	msg.WriteString(formatPageCounter(session, session.AdminState.CurrentPage, session.AdminState.LastPage))
	return msg.String()
}

// UserDetail generates a detailed message about a specific user's session and API details.
func UserDetail(session *models.Session, user *models.Session) string {
	return fmt.Sprintf("💻 %s:\n%s\n🌐 %s:\n%s",
		toBold(translator.Translate(session.Lang, "sessionDetails", nil, nil)),
		formatUserDetails(session, user),
		toBold(translator.Translate(session.Lang, "apiDetails", nil, nil)),
		formatAPIDetails(session, user))
}

// FeedbackList generates a message listing feedback entries with pagination details.
func FeedbackList(session *models.Session, feedbacks []models.Feedback) string {
	if session.AdminState.TotalRecords == 0 {
		return "❗️" + translator.Translate(session.Lang, "notFoundMessages", nil, nil)
	}

	var msg strings.Builder

	msg.WriteString(fmt.Sprintf("%s: %d\n\n",
		toBold(translator.Translate(session.Lang, "foundFeedbacks", nil, nil)),
		session.AdminState.TotalRecords))

	for i, feedback := range feedbacks {
		msg.WriteString(formatFeedback(session, feedback, i))
	}

	msg.WriteString(formatPageCounter(session, session.AdminState.CurrentPage, session.AdminState.LastPage))
	return msg.String()
}

// FeedbackDetail generates a detailed message about a specific feedback entry.
func FeedbackDetail(session *models.Session, feedback *models.Feedback) string {
	return fmt.Sprintf("💬 %s:\n%s: %s\n%s\n%s: %s\n%s:\n%s\n\n%s: %s %s\n",
		toBold(translator.Translate(session.Lang, "feedbackDetails", nil, nil)),
		toBold("Telegram ID"),
		toCode(strconv.Itoa(feedback.TelegramID)),
		formatOptionalString(toBold("Username"), feedback.TelegramUsername, "%s: @%s\n"),
		toBold(translator.Translate(session.Lang, "category", nil, nil)),
		translator.Translate(session.Lang, feedback.Category, nil, nil),
		toBold(translator.Translate(session.Lang, "message", nil, nil)),
		toPre(feedback.Message),
		toBold(translator.Translate(session.Lang, "created", nil, nil)),
		feedback.CreatedAt.Format("02.01.2006 15:04"),
		toItalic(fmt.Sprintf("(%d)", feedback.ID)))
}

// BanNotification generates a ban notification message for the user being banned.
func BanNotification(session *models.Session, reason string) string {
	return fmt.Sprintf("❌ %s\n\n%s⛔ %s",
		translator.Translate(session.AdminState.UserLang, "bannedBy", map[string]interface{}{
			"Role": translator.Translate(session.AdminState.UserLang, session.Role.String(), nil, nil),
		}, nil),
		formatOptionalString(toBold(translator.Translate(session.AdminState.UserLang, "reason", nil, nil)), reason, "%s: %s\n\n"),
		translator.Translate(session.AdminState.UserLang, "botAccessDenied", nil, nil))
}

// Ban generates a ban confirmation message for the admin banning a user.
func Ban(session *models.Session, reason string) string {
	return fmt.Sprintf("❌ %s%s",
		translator.Translate(session.Lang, "banUser", map[string]interface{}{
			"ID": session.AdminState.UserID,
		}, nil),
		formatOptionalString(toBold(translator.Translate(session.Lang, "reason", nil, nil)),
			reason, "\n\n%s: %s"))
}

// UnbanNotification generates an unban notification message for the user being unbanned.
func UnbanNotification(session *models.Session) string {
	return "✅ " + translator.Translate(session.AdminState.UserLang, "unbannedBy", map[string]interface{}{
		"Role": translator.Translate(session.AdminState.UserLang, session.Role.String(), nil, nil),
	}, nil)
}

// Unban generates an unban confirmation message for the admin unbanning a user.
func Unban(session *models.Session) string {
	return "✅ " + translator.Translate(session.Lang, "unbanUser", map[string]interface{}{
		"ID": session.AdminState.UserID,
	}, nil)
}

// ShiftRoleNotification generates a role shift notification message for the user whose role is being shifted.
func ShiftRoleNotification(session *models.Session, raise bool) string {
	emoji, newRole, action := getShiftRoleData(session, raise)
	return fmt.Sprintf("%s %s",
		emoji,
		translator.Translate(session.AdminState.UserLang, action+"To", map[string]interface{}{
			"Role": toBold(translator.Translate(session.AdminState.UserLang, newRole, nil, nil)),
		}, nil))
}

// ShiftRole generates a role shift confirmation message for the admin shifting a user's role.
func ShiftRole(session *models.Session, raise bool) string {
	emoji, newRole, action := getShiftRoleData(session, raise)
	return fmt.Sprintf("%s %s",
		emoji,
		translator.Translate(session.Lang, action+"User", map[string]interface{}{
			"ID":   toCode(fmt.Sprintf("%d", session.AdminState.UserID)),
			"Role": toBold(translator.Translate(session.Lang, newRole, nil, nil)),
		}, nil))
}

// RemoveRoleNotification generates a role removal notification message for the user whose role is being removed.
func RemoveRoleNotification(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.AdminState.UserLang, "removedRole", nil, nil)
}

// RemoveRole generates a role removal confirmation message for the admin removing a user's role.
func RemoveRole(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "removeUserRole", map[string]interface{}{
		"ID": toCode(strconv.Itoa(session.AdminState.UserID)),
	}, nil)
}

// ChangeRoleNotification generates a role change notification message for the user whose role is being updated.
func ChangeRoleNotification(session *models.Session, newRole roles.Role) string {
	return "🔄 " + translator.Translate(session.AdminState.UserLang, "changedRole", map[string]interface{}{
		"OldRole": toBold(translator.Translate(session.AdminState.UserLang, session.AdminState.UserRole.String(), nil, nil)),
		"NewRole": toBold(translator.Translate(session.AdminState.UserLang, newRole.String(), nil, nil)),
	}, nil)
}

// ChangeRole generates a role change confirmation message for the admin updating a user's role.
func ChangeRole(session *models.Session, newRole roles.Role) string {
	return "🔄 " + translator.Translate(session.Lang, "changeUserRole", map[string]interface{}{
		"ID":      toCode(fmt.Sprintf("%d", session.AdminState.UserID)),
		"OldRole": toBold(translator.Translate(session.Lang, session.AdminState.UserRole.String(), nil, nil)),
		"NewRole": toBold(translator.Translate(session.AdminState.UserLang, newRole.String(), nil, nil)),
	}, nil)
}

// AdminMenu generates a message for the admin panel menu.
func AdminMenu(session *models.Session) string {
	return fmt.Sprintf("🛠️ %s\n\n%s",
		toBold(translator.Translate(session.Lang, "adminPanel", nil, nil)),
		translator.Translate(session.Lang, "choiceAction", nil, nil))
}

// RequestEntityField generates a message prompting the admin to enter a user ID or username.
func RequestEntityField(session *models.Session) string {
	return fmt.Sprintf("%s\n\n%s",
		translator.Translate(session.Lang, "requestIDOrUsername", nil, nil),
		translator.Translate(session.Lang, "hintAPIUserID", nil, nil))
}

// NoAccess generates a message indicating that the user does not have access to a specific feature.
func NoAccess(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "noAccess", nil, nil)
}

// RequestBroadcastMessage generates a message prompting the admin to enter a broadcast message.
func RequestBroadcastMessage(session *models.Session) string {
	return "💬 " + translator.Translate(session.Lang, "requestBroadcastMessage", nil, nil)
}

// RequestBroadcastImage generates a message prompting the admin to upload a broadcast image.
func RequestBroadcastImage(session *models.Session) string {
	return "🏞️ " + translator.Translate(session.Lang, "requestBroadcastImage", nil, nil)
}

// RequestBroadcastPin generates a message prompting the admin to confirm pinning the broadcast message.
func RequestBroadcastPin(session *models.Session) string {
	return "📌 " + translator.Translate(session.Lang, "requestBroadcastPin", nil, nil)
}

// BroadcastPreview generates a preview of the broadcast message.
func BroadcastPreview(session *models.Session) string {
	return fmt.Sprintf("👁️ %s:\n\n%s",
		toItalic(translator.Translate(session.Lang, "preview", nil, nil)),
		session.AdminState.Message)
}

// BroadcastEmpty generates a message indicating that the broadcast content is empty.
func BroadcastEmpty(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "broadcastEmpty", nil, nil)
}

// BroadcastConfirm generates a confirmation message for sending the broadcast, including recipient count and pin status.
func BroadcastConfirm(session *models.Session, count int64) string {
	return fmt.Sprintf("👥 %s: %d%s",
		toBold(translator.Translate(session.Lang, "recipientCount", nil, nil)),
		count,
		formatOptionalBool(translator.Translate(session.Lang, "messageWillBePin", nil, nil),
			session.AdminState.NeedPin, "\n\n📌 %s"))
}

// FeedbackDeleteSuccess generates a success message after deleting a feedback entry.
func FeedbackDeleteSuccess(session *models.Session) string {
	return "🗑️ " + translator.Translate(session.Lang, "deleteFeedbackSuccess", map[string]interface{}{
		"ID": session.AdminState.FeedbackID,
	}, nil)
}

// LogsNotFound generates a message indicating that no logs were found.
func LogsNotFound(session *models.Session) string {
	return "❗" + translator.Translate(session.Lang, "logsNotFound", nil, nil)
}

// LogsFound generates a message indicating that logs were found for a specific user.
func LogsFound(session *models.Session) string {
	return "💾 " + translator.Translate(session.Lang, "logsFound", map[string]interface{}{
		"ID": session.AdminState.UserID,
	}, nil)
}

// NeedRemoveRole generates a message indicating that a role needs to be removed.
func NeedRemoveRole(session *models.Session) string {
	return "❗" + translator.Translate(session.Lang, "needRemoveRole", nil, nil)
}

// RequestBanReason generates a message prompting the admin to enter a reason for banning a user.
func RequestBanReason(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "requestBanReason", nil, nil)
}

// ChoiceRole generates a message prompting the admin to choose a role for a user.
func ChoiceRole(session *models.Session) string {
	return fmt.Sprintf("%s: %s\n\n%s",
		toBold(translator.Translate(session.Lang, "currentRole", nil, nil)),
		translator.Translate(session.Lang, session.AdminState.UserRole.String(), nil, nil),
		translator.Translate(session.Lang, "choiceRole", nil, nil))
}

// formatUserEntry formats a single user entry in a list of users with details like role, Telegram ID, and access status.
func formatUserEntry(session *models.Session, user models.Session, index int) string {
	return fmt.Sprintf("%s %s\n%s%s%s: %s\n\n",
		utils.NumberToEmoji(utils.GetItemID(index, session.AdminState.CurrentPage, session.AdminState.PageSize)),
		toBold(translator.Translate(session.Lang, user.Role.String(), nil, nil)),
		formatOptionalString("Username", user.TelegramUsername, "%s: @%s\n"),
		formatOptionalNumber("API ID", user.User.ID, -1, "%s: %d\n"),
		translator.Translate(session.Lang, "access", nil, nil),
		utils.BoolToEmojiColored(!user.IsBanned))
}

// formatUserDetails formats detailed information about a user's session, including role, Telegram ID, language, state, and timestamps.
func formatUserDetails(session *models.Session, user *models.Session) string {
	return fmt.Sprintf("%s: %s\n%s: %s\n%s%s: %s\n%s: %s\n%s: %s\n%s: %s\n%s: %s\n",
		toBold(translator.Translate(session.Lang, "role", nil, nil)),
		translator.Translate(session.Lang, user.Role.String(), nil, nil),
		toBold("Telegram ID"),
		toCode(strconv.Itoa(user.TelegramID)),
		formatOptionalString(toBold("Username"), user.TelegramUsername, "%s: @%s\n"),
		toBold(translator.Translate(session.Lang, "access", nil, nil)),
		utils.BoolToEmojiColored(!user.IsBanned),
		toBold(translator.Translate(session.Lang, "language", nil, nil)),
		nonEmpty(user.Lang, translator.Translate(session.Lang, "empty", nil, nil)),
		toBold(translator.Translate(session.Lang, "state", nil, nil)),
		nonEmpty(user.State, translator.Translate(session.Lang, "empty", nil, nil)),
		toBold(translator.Translate(session.Lang, "created", nil, nil)),
		user.CreatedAt.Format("02.01.2006 15:04"),
		toBold(translator.Translate(session.Lang, "updated", nil, nil)),
		user.UpdatedAt.Format("02.01.2006 15:04"))
}

// formatAPIDetails formats API-related details for a user, such as ID, username, email, and creation timestamp.
func formatAPIDetails(session *models.Session, user *models.Session) string {
	return fmt.Sprintf("%s: %s\n%s: %s\n%s%s: %s\n",
		toBold("ID"),
		toCode(strconv.Itoa(user.User.ID)),
		toBold(translator.Translate(session.Lang, "name", nil, nil)),
		toCode(user.User.Username),
		formatOptionalString(toBold("Email"), user.User.Email, "%s: %s\n"),
		toBold(translator.Translate(session.Lang, "created", nil, nil)),
		user.User.CreatedAt.Format("02.01.2006 15:04"))
}

// formatFeedback formats a single feedback entry in a list of feedbacks with details like category, Telegram ID, and creation timestamp.
func formatFeedback(session *models.Session, feedback models.Feedback, index int) string {
	return fmt.Sprintf("%s %s\n%s%s%s: %s %s\n\n",
		utils.NumberToEmoji(utils.GetItemID(index, session.AdminState.CurrentPage, session.AdminState.PageSize)),
		toBold(translator.Translate(session.Lang, feedback.Category, nil, nil)),
		formatOptionalNumber("Telegram ID", feedback.TelegramID, -1, "%s: <code>%d</code>\n"),
		formatOptionalString("Username", feedback.TelegramUsername, "%s: @%s\n"),
		translator.Translate(session.Lang, "created", nil, nil),
		feedback.CreatedAt.Format("02.01.2006 15:04"),
		toItalic(fmt.Sprintf("(%d)", feedback.ID)))
}

// getShiftRoleData determines the emoji, new role, and action type (raise/lower) based on whether the role is being increased or decreased.
func getShiftRoleData(session *models.Session, raise bool) (string, string, string) {
	if raise {
		return "⬆️", session.AdminState.UserRole.NextRole().String(), "raise"
	}
	return "⬇️", session.AdminState.UserRole.PrevRole().String(), "lower"
}
