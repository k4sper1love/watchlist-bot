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

func UserDetail(session *models.Session, user *models.Session) string {
	return fmt.Sprintf("💻 %s:\n%s\n🌐 %s:\n%s",
		toBold(translator.Translate(session.Lang, "sessionDetails", nil, nil)),
		formatUserDetails(session, user),
		toBold(translator.Translate(session.Lang, "apiDetails", nil, nil)),
		formatAPIDetails(session, user))
}

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

func BanNotification(session *models.Session, reason string) string {
	return fmt.Sprintf("❌ %s\n\n%s⛔ %s",
		translator.Translate(session.AdminState.UserLang, "bannedBy", map[string]interface{}{
			"Role": translator.Translate(session.AdminState.UserLang, session.Role.String(), nil, nil),
		}, nil),
		formatOptionalString(toBold(translator.Translate(session.AdminState.UserLang, "reason", nil, nil)), reason, "%s: %s\n\n"),
		translator.Translate(session.AdminState.UserLang, "botAccessDenied", nil, nil))
}

func Ban(session *models.Session, reason string) string {
	return fmt.Sprintf("❌ %s%s",
		translator.Translate(session.Lang, "banUser", map[string]interface{}{
			"ID": session.AdminState.UserID,
		}, nil),
		formatOptionalString(toBold(translator.Translate(session.Lang, "reason", nil, nil)),
			reason, "\n\n%s: %s"))
}

func UnbanNotification(session *models.Session) string {
	return "✅ " + translator.Translate(session.AdminState.UserLang, "unbannedBy", map[string]interface{}{
		"Role": translator.Translate(session.AdminState.UserLang, session.Role.String(), nil, nil),
	}, nil)
}

func Unban(session *models.Session) string {
	return "✅ " + translator.Translate(session.Lang, "unbanUser", map[string]interface{}{
		"ID": session.AdminState.UserID,
	}, nil)
}

func ShiftRoleNotification(session *models.Session, raise bool) string {
	emoji, newRole, action := getShiftRoleData(session, raise)
	return fmt.Sprintf("%s %s",
		emoji,
		translator.Translate(session.AdminState.UserLang, action+"To", map[string]interface{}{
			"Role": toBold(translator.Translate(session.AdminState.UserLang, newRole, nil, nil)),
		}, nil))
}

func ShiftRole(session *models.Session, raise bool) string {
	emoji, newRole, action := getShiftRoleData(session, raise)
	return fmt.Sprintf("%s %s",
		emoji,
		translator.Translate(session.Lang, action+"User", map[string]interface{}{
			"ID":   toCode(fmt.Sprintf("%d", session.AdminState.UserID)),
			"Role": toBold(translator.Translate(session.Lang, newRole, nil, nil)),
		}, nil))
}

func RemoveRoleNotification(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.AdminState.UserLang, "removedRole", nil, nil)
}

func RemoveRole(session *models.Session) string {
	return "⚠️ " + translator.Translate(session.Lang, "removeUserRole", map[string]interface{}{
		"ID": toCode(strconv.Itoa(session.AdminState.UserID)),
	}, nil)
}

func ChangeRoleNotification(session *models.Session, newRole roles.Role) string {
	return "🔄 " + translator.Translate(session.AdminState.UserLang, "changedRole", map[string]interface{}{
		"OldRole": toBold(translator.Translate(session.AdminState.UserLang, session.AdminState.UserRole.String(), nil, nil)),
		"NewRole": toBold(translator.Translate(session.AdminState.UserLang, newRole.String(), nil, nil)),
	}, nil)
}

func ChangeRole(session *models.Session, newRole roles.Role) string {
	return "🔄 " + translator.Translate(session.Lang, "changeUserRole", map[string]interface{}{
		"ID":      toCode(fmt.Sprintf("%d", session.AdminState.UserID)),
		"OldRole": toBold(translator.Translate(session.Lang, session.AdminState.UserRole.String(), nil, nil)),
		"NewRole": toBold(translator.Translate(session.AdminState.UserLang, newRole.String(), nil, nil)),
	}, nil)
}

func AdminMenu(session *models.Session) string {
	return fmt.Sprintf("🛠️ %s\n\n%s",
		toBold(translator.Translate(session.Lang, "adminPanel", nil, nil)),
		translator.Translate(session.Lang, "choiceAction", nil, nil))
}

func RequestEntityField(session *models.Session) string {
	return fmt.Sprintf("%s\n\n%s",
		translator.Translate(session.Lang, "requestIDOrUsername", nil, nil),
		translator.Translate(session.Lang, "hintAPIUserID", nil, nil))
}

func NoAccess(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "noAccess", nil, nil)
}

func RequestBroadcastMessage(session *models.Session) string {
	return "💬 " + translator.Translate(session.Lang, "requestBroadcastMessage", nil, nil)
}

func RequestBroadcastImage(session *models.Session) string {
	return "🏞️ " + translator.Translate(session.Lang, "requestBroadcastImage", nil, nil)
}

func RequestBroadcastPin(session *models.Session) string {
	return "📌 " + translator.Translate(session.Lang, "requestBroadcastPin", nil, nil)
}

func BroadcastPreview(session *models.Session) string {
	return fmt.Sprintf("👁️ %s:\n\n%s",
		toItalic(translator.Translate(session.Lang, "preview", nil, nil)),
		session.AdminState.Message)
}

func BroadcastEmpty(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "broadcastEmpty", nil, nil)
}

func BroadcastConfirm(session *models.Session, count int64) string {
	return fmt.Sprintf("👥 %s: %d%s",
		toBold(translator.Translate(session.Lang, "recipientCount", nil, nil)),
		count,
		formatOptionalBool(translator.Translate(session.Lang, "messageWillBePin", nil, nil),
			session.AdminState.NeedFeedbackPin, "\n\n📌 %s"))
}

func FeedbackDeleteSuccess(session *models.Session) string {
	return "🗑️ " + translator.Translate(session.Lang, "deleteFeedbackSuccess", map[string]interface{}{
		"ID": session.AdminState.FeedbackID,
	}, nil)
}

func LogsNotFound(session *models.Session) string {
	return "❗" + translator.Translate(session.Lang, "logsNotFound", nil, nil)
}

func LogsFound(session *models.Session) string {
	return "💾 " + translator.Translate(session.Lang, "logsFound", map[string]interface{}{
		"ID": session.AdminState.UserID,
	}, nil)
}

func NeedRemoveRole(session *models.Session) string {
	return "❗" + translator.Translate(session.Lang, "needRemoveRole", nil, nil)
}

func RequestBanReason(session *models.Session) string {
	return "❓" + translator.Translate(session.Lang, "requestBanReason", nil, nil)
}

func ChoiceRole(session *models.Session) string {
	return fmt.Sprintf("%s: %s\n\n%s",
		toBold(translator.Translate(session.Lang, "currentRole", nil, nil)),
		translator.Translate(session.Lang, session.AdminState.UserRole.String(), nil, nil),
		translator.Translate(session.Lang, "choiceRole", nil, nil))
}

func formatUserEntry(session *models.Session, user models.Session, index int) string {
	return fmt.Sprintf("%s %s\n%s%s%s: %s\n\n",
		utils.NumberToEmoji(utils.GetItemID(index, session.AdminState.CurrentPage, session.AdminState.PageSize)),
		toBold(translator.Translate(session.Lang, user.Role.String(), nil, nil)),
		formatOptionalString("Username", user.TelegramUsername, "%s: @%s\n"),
		formatOptionalNumber("API ID", user.User.ID, -1, "%s: %d\n"),
		translator.Translate(session.Lang, "access", nil, nil),
		utils.BoolToEmojiColored(!user.IsBanned))
}

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

func getShiftRoleData(session *models.Session, raise bool) (string, string, string) {
	if raise {
		return "⬆️", session.AdminState.UserRole.NextRole().String(), "raise"
	}
	return "⬇️", session.AdminState.UserRole.PrevRole().String(), "lower"
}
