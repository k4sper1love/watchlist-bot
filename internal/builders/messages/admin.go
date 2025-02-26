package messages

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func BuildAdminUserListMessage(session *models.Session, users []models.Session) string {
	currentPage := session.AdminState.CurrentPage
	pageSize := session.AdminState.PageSize
	lastPage := session.AdminState.LastPage
	totalRecords := session.AdminState.TotalRecords

	foundMsg := translator.Translate(session.Lang, "foundUsers", nil, nil)
	msg := fmt.Sprintf("<b>%s:</b> %d\n\n", foundMsg, totalRecords)

	for i, user := range users {
		itemID := utils.GetItemID(i, currentPage, pageSize)

		numberEmoji := utils.NumberToEmoji(itemID)

		roleMsg := translator.Translate(session.Lang, user.Role.String(), nil, nil)
		msg += fmt.Sprintf("<b>%s %s</b>\n", numberEmoji, roleMsg)

		msg += fmt.Sprintf("Telegram ID: <code>%d</code>\n", user.TelegramID)

		if user.TelegramUsername != "" {
			msg += fmt.Sprintf("Username: @%s\n", user.TelegramUsername)
		}

		if user.User.ID != 1 {
			msg += fmt.Sprintf("API ID: <code>%d</code>\n", user.User.ID)
		}

		bannedMsg := translator.Translate(session.Lang, "access", nil, nil)
		msg += fmt.Sprintf("%s: %s\n", bannedMsg, utils.BoolToEmojiColored(!user.IsBanned))

		msg += "\n"
	}

	pageMsg := translator.Translate(session.Lang, "pageCounter", map[string]interface{}{
		"CurrentPage": currentPage,
		"LastPage":    lastPage,
	}, nil)

	msg += fmt.Sprintf("<b>📄 %s</b>", pageMsg)

	return msg
}

func BuildAdminUserDetailMessage(session *models.Session, user *models.Session) string {
	detailsMsg := translator.Translate(session.Lang, "sessionDetails", nil, nil)
	msg := fmt.Sprintf("💻 <b>%s:</b>\n", detailsMsg)

	roleMsg := translator.Translate(session.Lang, "role", nil, nil)
	roleValueMsg := translator.Translate(session.Lang, user.Role.String(), nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> %s\n", roleMsg, roleValueMsg)

	msg += fmt.Sprintf("<b>Telegram ID:</b> <code>%d</code>\n", user.TelegramID)

	if user.TelegramUsername != "" {
		msg += fmt.Sprintf("<b>Username:</b> @%s\n", user.TelegramUsername)
	}
	bannedMsg := translator.Translate(session.Lang, "access", nil, nil)
	msg += fmt.Sprintf("<b>%s</b>: %s\n", bannedMsg, utils.BoolToEmojiColored(!user.IsBanned))

	languageMsg := translator.Translate(session.Lang, "language", nil, nil)
	language := ""
	if user.Lang == "" {
		language = translator.Translate(session.Lang, "empty", nil, nil)
	} else {
		language = user.Lang
	}
	msg += fmt.Sprintf("<b>%s:</b> %s\n", languageMsg, language)

	stateMsg := translator.Translate(session.Lang, "state", nil, nil)
	state := ""
	if user.State == "" {
		state = translator.Translate(session.Lang, "empty", nil, nil)
	} else {
		state = translator.Translate(session.Lang, "empty", nil, nil)
	}
	msg += fmt.Sprintf("<b>%s:</b> %s\n", stateMsg, state)

	createdMsg := translator.Translate(session.Lang, "created", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> %s\n", createdMsg, user.CreatedAt.Format("02.01.2006 15:04"))

	updatedMsg := translator.Translate(session.Lang, "updated", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> %s\n\n", updatedMsg, user.UpdatedAt.Format("02.01.2006 15:04"))

	detailsMsg = translator.Translate(session.Lang, "apiDetails", nil, nil)
	msg += fmt.Sprintf("🌐 <b>%s:</b>\n", detailsMsg)

	msg += fmt.Sprintf("<b>ID:</b> <code>%d</code>\n", user.User.ID)

	nameMsg := translator.Translate(session.Lang, "name", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> <code>%s</code>\n", nameMsg, user.User.Username)

	if user.User.Email != "" {
		msg += fmt.Sprintf("<b>Email:</b> <code>%s</code>\n", user.User.Email)
	}

	createdMsg = translator.Translate(session.Lang, "created", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> %s\n", createdMsg, user.User.CreatedAt.Format("02.01.2006 15:04"))

	return msg
}

func BuildFeedbackListMessage(session *models.Session, feedbacks []models.Feedback) string {
	currentPage := session.AdminState.CurrentPage
	pageSize := session.AdminState.PageSize
	lastPage := session.AdminState.LastPage
	totalRecords := session.AdminState.TotalRecords

	if len(feedbacks) == 0 {
		msg := "❗️" + translator.Translate(session.Lang, "notFoundMessages", nil, nil)
		return msg
	}

	foundMsg := translator.Translate(session.Lang, "foundFeedbacks", nil, nil)
	msg := fmt.Sprintf("<b>%s:</b> %d\n\n", foundMsg, totalRecords)

	for i, feedback := range feedbacks {
		itemID := utils.GetItemID(i, currentPage, pageSize)

		numberEmoji := utils.NumberToEmoji(itemID)

		categoryMsg := translator.Translate(session.Lang, feedback.Category, nil, nil)
		msg += fmt.Sprintf("<b>%s %s</b>\n", numberEmoji, categoryMsg)

		msg += fmt.Sprintf("Telegram ID: <code>%d</code>\n", feedback.TelegramID)

		if feedback.TelegramUsername != "" {
			msg += fmt.Sprintf("Username: @%s\n", feedback.TelegramUsername)
		}

		createdMsg := translator.Translate(session.Lang, "created", nil, nil)
		msg += fmt.Sprintf("%s: %s", createdMsg, feedback.CreatedAt.Format("02.01.2006 15:04"))

		msg += fmt.Sprintf(" <i>(%d)</i>\n\n", feedback.ID)
	}

	pageMsg := translator.Translate(session.Lang, "pageCounter", map[string]interface{}{
		"CurrentPage": currentPage,
		"LastPage":    lastPage,
	}, nil)

	msg += fmt.Sprintf("<b>📄 %s</b>", pageMsg)

	return msg
}

func BuildFeedbackDetailMessage(session *models.Session, feedback *models.Feedback) string {
	detailsMsg := "💬 " + translator.Translate(session.Lang, "feedbackDetails", nil, nil)
	msg := fmt.Sprintf("<b>%s:</b>\n", detailsMsg)

	msg += fmt.Sprintf("<b>Telegram ID</b>: <code>%d</code>\n", feedback.TelegramID)

	if feedback.TelegramUsername != "" {
		msg += fmt.Sprintf("<b>Username</b>: @%s\n", feedback.TelegramUsername)
	}

	msg += fmt.Sprintf("\n")

	categoryMsg := translator.Translate(session.Lang, "category", nil, nil)
	translatedCategory := translator.Translate(session.Lang, feedback.Category, nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> %s\n\n", categoryMsg, translatedCategory)

	feedbackMsg := translator.Translate(session.Lang, "message", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b>\n<pre>%s</pre>\n\n", feedbackMsg, feedback.Message)

	createdMsg := translator.Translate(session.Lang, "created", nil, nil)
	msg += fmt.Sprintf("<b>%s:</b> %s", createdMsg, feedback.CreatedAt.Format("02.01.2006 15:04"))

	msg += fmt.Sprintf(" <i>(%d)</i>\n", feedback.ID)

	return msg
}

func BuildUserBanNotificationMessage(session *models.Session, reason string) string {
	roleMsg := translator.Translate(session.AdminState.UserLang, session.Role.String(), nil, nil)
	bannedByMsg := translator.Translate(session.AdminState.UserLang, "bannedBy", map[string]interface{}{
		"Role": roleMsg,
	}, nil)

	msg := fmt.Sprintf("❌ %s", bannedByMsg)

	if reason != "" {
		reasonMsg := translator.Translate(session.AdminState.UserLang, "reason", nil, nil)
		msg += fmt.Sprintf("\n\n<b>%s</b>: %s", reasonMsg, reason)
	}

	accessDeniedMsg := translator.Translate(session.AdminState.UserLang, "botAccessDenied", nil, nil)
	msg += fmt.Sprintf("\n\n⛔ %s", accessDeniedMsg)

	return msg
}

func BuildBanMessage(session *models.Session, reason string) string {
	banUserMsg := translator.Translate(session.Lang, "banUser", map[string]interface{}{
		"ID": session.AdminState.UserID,
	}, nil)

	msg := fmt.Sprintf("❌ %s", banUserMsg)

	if reason != "" {
		reasonMsg := translator.Translate(session.Lang, "reason", nil, nil)
		msg += fmt.Sprintf("\n\n<b>%s</b>: %s", reasonMsg, reason)
	}

	return msg
}

func BuildUserUnbanNotificationMessage(session *models.Session) string {
	roleMsg := translator.Translate(session.AdminState.UserLang, session.Role.String(), nil, nil)
	bannedByMsg := translator.Translate(session.AdminState.UserLang, "unbannedBy", map[string]interface{}{
		"Role": roleMsg,
	}, nil)

	msg := fmt.Sprintf("✅ %s", bannedByMsg)

	return msg
}

func BuildUnbanMessage(session *models.Session) string {
	unbanUserMsg := translator.Translate(session.Lang, "unbanUser", map[string]interface{}{
		"ID": session.AdminState.UserID,
	}, nil)

	msg := fmt.Sprintf("✅ %s", unbanUserMsg)

	return msg
}

func BuildRaiseRoleNotificationMessage(session *models.Session) string {
	roleMsg := translator.Translate(session.AdminState.UserLang, session.AdminState.UserRole.NextRole().String(), nil, nil)
	raiseMsg := translator.Translate(session.AdminState.UserLang, "raisedTo", map[string]interface{}{
		"Role": toBold(roleMsg),
	}, nil)

	msg := fmt.Sprintf("⬆️ %s", raiseMsg)

	return msg
}

func BuildRoleChangeMessage(session *models.Session, raise bool) string {
	var newRole, messageID string

	if raise {
		newRole = session.AdminState.UserRole.NextRole().String()
		messageID = "raiseUser"
	} else {
		newRole = session.AdminState.UserRole.PrevRole().String()
		messageID = "lowerUser"
	}

	roleMsg := translator.Translate(session.Lang, newRole, nil, nil)
	statusMsg := translator.Translate(session.Lang, messageID, map[string]interface{}{
		"ID":   toCode(fmt.Sprintf("%d", session.AdminState.UserID)),
		"Role": toBold(roleMsg),
	}, nil)

	return fmt.Sprintf("⬆️ %s", statusMsg)
}

func BuildLowerRoleNotificationMessage(session *models.Session) string {
	roleMsg := translator.Translate(session.AdminState.UserLang, session.AdminState.UserRole.PrevRole().String(), nil, nil)
	raiseMsg := translator.Translate(session.AdminState.UserLang, "lowerTo", map[string]interface{}{
		"Role": toBold(roleMsg),
	}, nil)

	msg := fmt.Sprintf("⬇️ %s", raiseMsg)

	return msg
}

func BuildRemoveRoleNotificationMessage(session *models.Session) string {
	removedMsg := translator.Translate(session.AdminState.UserLang, "removedRole", nil, nil)

	msg := fmt.Sprintf("⚠️ %s", removedMsg)

	return msg
}

func BuildRemoveRoleMessage(session *models.Session) string {
	removeMsg := translator.Translate(session.Lang, "removeUserRole", map[string]interface{}{
		"ID": toCode(fmt.Sprintf("%d", session.AdminState.UserID)),
	}, nil)

	msg := fmt.Sprintf("⚠️ %s", removeMsg)

	return msg
}

func BuildChangeRoleNotificationMessage(session *models.Session, newRole roles.Role) string {
	oldRoleMsg := translator.Translate(session.AdminState.UserLang, session.AdminState.UserRole.String(), nil, nil)
	newRoleMsg := translator.Translate(session.AdminState.UserLang, newRole.String(), nil, nil)

	changedMsg := translator.Translate(session.AdminState.UserLang, "changedRole", map[string]interface{}{
		"OldRole": toBold(oldRoleMsg),
		"NewRole": toBold(newRoleMsg),
	}, nil)

	msg := fmt.Sprintf("🔄 %s", changedMsg)

	return msg
}

func BuildChangeRoleMessage(session *models.Session, newRole roles.Role) string {
	oldRoleMsg := translator.Translate(session.Lang, session.AdminState.UserRole.String(), nil, nil)
	newRoleMsg := translator.Translate(session.AdminState.UserLang, newRole.String(), nil, nil)

	changeMsg := translator.Translate(session.Lang, "changeUserRole", map[string]interface{}{
		"ID":      toCode(fmt.Sprintf("%d", session.AdminState.UserID)),
		"OldRole": toBold(oldRoleMsg),
		"NewRole": toBold(newRoleMsg),
	}, nil)

	msg := fmt.Sprintf("🔄 %s", changeMsg)

	return msg
}

func BuildAdminMenuMessage(session *models.Session) string {
	part1 := translator.Translate(session.Lang, "adminPanel", nil, nil)
	part2 := translator.Translate(session.Lang, "choiceAction", nil, nil)
	return fmt.Sprintf("🛠️ <b>%s</b>\n\n%s", part1, part2)
}

func BuildAdminRequestIDOrUsernameMessage(session *models.Session) string {
	return translator.Translate(session.Lang, "requestIDOrUsername", nil, nil)
}

func BuildNoAccessMessage(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "noAccess", nil, nil)
}

func BuildRequestBroadcastMessage(session *models.Session) string {
	return "💬 " + translator.Translate(session.Lang, "requestBroadcastMessage", nil, nil)
}

func BuildRequestBroadcastImageMessage(session *models.Session) string {
	return "🏞️ " + translator.Translate(session.Lang, "requestBroadcastImage", nil, nil)
}

func BuildRequestBroadcastPinMessage(session *models.Session) string {
	return "📌 " + translator.Translate(session.Lang, "requestBroadcastPin", nil, nil)
}

func BuildBroadcastPreviewMessage(session *models.Session) string {
	previewMsg := translator.Translate(session.Lang, "preview", nil, nil)
	return fmt.Sprintf("👁️ <i>%s:</i>\n\n%s", previewMsg, session.AdminState.Message)
}

func BuildBroadcastEmptyMessage(session *models.Session) string {
	return "❗️" + translator.Translate(session.Lang, "broadcastEmpty", nil, nil)
}

func BuildBroadcastConfirmMessage(session *models.Session, count int64) string {
	countMsg := "👥 " + translator.Translate(session.Lang, "recipientCount", nil, nil)
	msg := fmt.Sprintf("<b>%s</b>: %d", countMsg, count)

	if session.AdminState.NeedFeedbackPin {
		msg += "\n\n📌 " + translator.Translate(session.Lang, "messageWillBePin", nil, nil)
	}

	return msg
}

func BuildFeedbackDeleteSuccessMessage(session *models.Session) string {
	return "🗑️ " + translator.Translate(session.Lang, "deleteFeedbackSuccess", map[string]interface{}{
		"ID": session.AdminState.FeedbackID,
	}, nil)
}
