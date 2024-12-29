package admin

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"log"
	"strconv"
	"strings"
)

var adminButtons = []keyboards.Button{
	{"", "adminOptionUserCount", states.CallbackAdminSelectUserCount},
	{"", "adminOptionBroadcast", states.CallbackAdminSelectBroadcastMessage},
	{"", "adminOptionFeedback", states.CallbackAdminSelectFeedback},
	{"", "adminOptionUserList", states.CallbackAdminSelectUsers},
}

func HandleAdminCommand(app models.App, session *models.Session) {
	msg := translator.Translate(session.Lang, "choiceAction", nil, nil)

	keyboard := keyboards.NewKeyboard().
		AddButtons(adminButtons...).
		AddBack(states.CallbackAdminSelectBack).
		Build(session.Lang)

	app.SendMessage(msg, keyboard)
}

func HandleAdminButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	command := utils.ParseMessageCommand(app.Upd)
	switch {
	case callback == states.CallbackAdminSelectBack:
		general.HandleMenuCommand(app, session)

	case callback == states.CallbackAdminSelectBackPanel:
		HandleAdminCommand(app, session)

	case callback == states.CallbackAdminSelectUserCount:
		handleAdminUserCount(app, session)

	case callback == states.CallbackAdminSelectBroadcastMessage:
		handleAdminBroadcastMessage(app, session)

	case callback == states.CallbackAdminSelectFeedback:
		handleAdminFeedback(app, session)

	case callback == states.CallbackAdminSelectUsers:
		handleAdminUsers(app, session)

	case strings.HasPrefix(command, "delete_feedback_"):
		handleDeleteFeedback(app, session)

	case strings.HasPrefix(command, "ban_"):
		handleBanUser(app, session)

	case strings.HasPrefix(command, "unban_"):
		handleUnbanUser(app, session)
	}

}

func HandleAdminProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		HandleAdminCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessAdminAwaitingBroadcastMessageText:
		parseAdminBroadcastMessageText(app, session)
	}
}

func handleAdminUserCount(app models.App, session *models.Session) {
	count, err := postgres.GetUserCounts()
	if err != nil {
		msg := translator.Translate(session.Lang, "requestFailure", nil, nil)
		app.SendMessage(msg, nil)
		return
	}

	uniqueUsersMsg := translator.Translate(session.Lang, "uniqueUsersCount", nil, nil)
	msg := fmt.Sprintf("%s: %d", uniqueUsersMsg, count)

	app.SendMessage(msg, nil)
	HandleAdminCommand(app, session)
}

func handleAdminBroadcastMessage(app models.App, session *models.Session) {
	count, err := postgres.GetUserCounts()
	if err != nil {
		msg := translator.Translate(session.Lang, "requestFailure", nil, nil)
		app.SendMessage(msg, nil)
		return
	}

	part1 := translator.Translate(session.Lang, "recipientCount", nil, nil)
	part2 := translator.Translate(session.Lang, "requestBroadcastMessage", nil, nil)
	msg := fmt.Sprintf("%s: %d\n\n%s", part1, count, part2)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessAdminAwaitingBroadcastMessageText)
}

func parseAdminBroadcastMessageText(app models.App, session *models.Session) {
	msg := utils.ParseMessageString(app.Upd)

	telegramIDs, err := postgres.GetAllTelegramID()
	if err != nil {
		msg = translator.Translate(session.Lang, "requestFailure", nil, nil)
		app.SendMessage(msg, nil)
		return
	}

	app.SendBroadcastMessage(telegramIDs, msg, nil)

	session.ClearState()

	HandleAdminCommand(app, session)
}

func handleAdminFeedback(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddBack(states.CallbackAdminSelectBackPanel).Build(session.Lang)

	feedbacks, err := postgres.FetchAllFeedbacks()
	if err != nil || len(feedbacks) == 0 {
		emptyListMsg := translator.Translate(session.Lang, "emptyFeedbackList", nil, nil)
		msg := fmt.Sprintf("📭 %s", emptyListMsg)
		app.SendMessage(msg, keyboard)
		return
	}

	const maxMessageLength = 4000

	feedbackListMsg := translator.Translate(session.Lang, "feedbackList", nil, nil)
	msg := fmt.Sprintf("📄 <b>%s:</b>\n\n", feedbackListMsg)

	for _, feedback := range feedbacks {
		idMsg := translator.Translate(session.Lang, "id", nil, nil)
		userMsg := translator.Translate(session.Lang, "user", nil, nil)
		categoryMsg := translator.Translate(session.Lang, "category", nil, nil)

		entry := fmt.Sprintf(
			"🆔 %s: %d\n"+
				"👤 %s: tg://user?id=%d\n"+
				"📂 %s: %s\n"+
				"💬 %s\n"+
				"📅 %s\n\n"+
				"🗑️ /delete_feedback_%d\n"+
				"━━━━━━━━━━━━━\n",
			idMsg, feedback.ID, userMsg, feedback.TelegramID, categoryMsg, feedback.Category, feedback.Message,
			feedback.CreatedAt.Format("02.01.2006 15:04"), feedback.ID,
		)

		if len(msg)+len(entry) > maxMessageLength {
			app.SendMessage(msg, keyboard)
			msg = ""
		}

		msg += entry
	}

	if len(msg) > 0 {
		app.SendMessage(msg, keyboard)
	}
}

func handleDeleteFeedback(app models.App, session *models.Session) {
	feedbackID, err := parseFeedbackID(app)
	if err != nil {
		msg := translator.Translate(session.Lang, "requestFailure", nil, nil)
		app.SendMessage(msg, nil)
		return
	}

	err = postgres.DeleteFeedbackByID(feedbackID)
	if err != nil {
		failureMsg := translator.Translate(session.Lang, "deleteFeedbackFailure", nil, nil)
		msg := fmt.Sprintf("❌ %s", failureMsg)
		app.SendMessage(msg, nil)
		return
	}

	successMsg := translator.Translate(session.Lang, "deleteFeedbackSuccess", nil, nil)
	msg := fmt.Sprintf("✅ %s", successMsg)
	app.SendMessage(msg, nil)

	handleAdminFeedback(app, session)
}

func handleAdminUsers(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddBack(states.CallbackAdminSelectBackPanel).Build(session.Lang)

	users, err := postgres.FetchAllUsers()
	if err != nil || len(users) == 0 {
		msg := translator.Translate(session.Lang, "emptyUserList", nil, nil)
		app.SendMessage(msg, keyboard)
		return
	}

	const maxMessageLength = 4000

	userListMsg := translator.Translate(session.Lang, "userList", nil, nil)
	msg := fmt.Sprintf("📄 <b>%s:</b>\n\n", userListMsg)

	for _, user := range users {
		telegramIDMsg := translator.Translate(session.Lang, "telegramID", nil, nil)
		adminMsg := translator.Translate(session.Lang, "admin", nil, nil)
		bannedMsg := translator.Translate(session.Lang, "banned", nil, nil)
		createdMsg := translator.Translate(session.Lang, "created", nil, nil)

		entry := fmt.Sprintf(
			"🆔 %s tg://user?id=%d\n"+
				"👤 <b>%s:</b> %s\n"+
				"🔐 <b>%s:</b> %s\n"+
				"📅 <b>%s:</b> %s\n",
			telegramIDMsg, user.TelegramID,
			adminMsg, boolToString(user.IsAdmin),
			bannedMsg, boolToString(user.IsBanned),
			createdMsg, user.CreatedAt.Format("02.01.2006 15:04"),
		)

		if !user.IsAdmin && !user.IsBanned {
			entry += fmt.Sprintf("🛑 /ban_%d\n", user.TelegramID)
		}

		if user.IsBanned {
			entry += fmt.Sprintf("🟢 /unban_%d\n", user.TelegramID)
		}

		entry += "━━━━━━━━━━━━━\n"

		if len(msg)+len(entry) > maxMessageLength {
			app.SendMessage(msg, keyboard)
			msg = ""
		}

		msg += entry
	}

	if len(msg) > 0 {
		app.SendMessage(msg, keyboard)
	}
}

func handleBanUser(app models.App, session *models.Session) {
	telegramID, err := parseBanTelegramID(app)
	if err != nil {
		badMsg := translator.Translate(session.Lang, "badFormatID", nil, nil)
		msg := fmt.Sprintf("❌ %s", badMsg)
		app.SendMessage(msg, nil)
		return
	}

	err = postgres.BanUser(telegramID)
	if err != nil {
		failureMsg := translator.Translate(session.Lang, "banFailure", nil, nil)
		msg := fmt.Sprintf("❌ %s", failureMsg)
		app.SendMessage(msg, nil)
		return
	}

	successMsg := translator.Translate(session.Lang, "banSuccess", map[string]interface{}{
		"User": telegramID,
	}, nil)

	msg := fmt.Sprintf("✅ %s", successMsg)

	app.SendMessage(msg, nil)

	handleAdminUsers(app, session)
}

func handleUnbanUser(app models.App, session *models.Session) {
	telegramID, err := parseUnbanTelegramID(app)
	if err != nil {
		badMsg := translator.Translate(session.Lang, "badFormatID", nil, nil)
		msg := fmt.Sprintf("❌ %s", badMsg)
		app.SendMessage(msg, nil)
		return
	}

	err = postgres.UnbanUser(telegramID)
	if err != nil {
		failureMsg := translator.Translate(session.Lang, "unbanFailure", nil, nil)
		msg := fmt.Sprintf("❌ %s", failureMsg)
		app.SendMessage(msg, nil)
		return
	}

	successMsg := translator.Translate(session.Lang, "unbanSuccess", map[string]interface{}{
		"User": telegramID,
	}, nil)

	msg := fmt.Sprintf("✅ %s", successMsg)

	app.SendMessage(msg, nil)

	handleAdminUsers(app, session)
}

func parseFeedbackID(app models.App) (int, error) {
	command := utils.ParseMessageCommand(app.Upd)
	idStr := strings.TrimPrefix(command, "delete_feedback_")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("error parsing feedback ID: %v", err)
		return -1, err
	}

	return id, nil
}

func parseBanTelegramID(app models.App) (int, error) {
	command := utils.ParseMessageCommand(app.Upd)
	idStr := strings.TrimPrefix(command, "ban_")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("error parsing ban telegram ID: %v", err)
		return -1, err
	}

	return id, nil
}

func parseUnbanTelegramID(app models.App) (int, error) {
	command := utils.ParseMessageCommand(app.Upd)
	idStr := strings.TrimPrefix(command, "unban_")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		log.Printf("error parsing unban telegram ID: %v", err)
		return -1, err
	}

	return id, nil
}

func boolToString(value bool) string {
	if value {
		return "✅"
	}
	return "❌"
}
