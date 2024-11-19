package admin

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log"
	"strconv"
	"strings"
)

var adminButtons = []keyboards.Button{
	{"Узнать количество пользателей", states.CallbackAdminSelectUserCount},
	{"Отправить рассылку", states.CallbackAdminSelectBroadcastMessage},
	{"Просмотр фидбеков", states.CallbackAdminSelectFeedback},
	{"Список пользователей", states.CallbackAdminSelectUsers},
}

func HandleAdminCommand(app models.App, session *models.Session) {
	msg := "Выберите действие"

	keyboard := keyboards.NewKeyboard().
		AddButtons(adminButtons...).
		AddBack(states.CallbackAdminSelectBack).
		Build()

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
		app.SendMessage("Не удалось получить инфу", nil)
		return
	}

	app.SendMessage(fmt.Sprintf("Уникальных юзеров бота: %d", count), nil)
	HandleAdminCommand(app, session)
}

func handleAdminBroadcastMessage(app models.App, session *models.Session) {
	count, err := postgres.GetUserCounts()
	if err != nil {
		app.SendMessage("Произошла ошибка при подсчете получателей", nil)
		return
	}

	msg := fmt.Sprintf("Количество получателей: %d\n", count)
	msg += "Введите сообщение, которое будет использвано для рассылки"

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessAdminAwaitingBroadcastMessageText)
}

func parseAdminBroadcastMessageText(app models.App, session *models.Session) {
	msg := utils.ParseMessageString(app.Upd)

	telegramIDs, err := postgres.GetAllTelegramID()
	if err != nil {
		app.SendMessage("Ошибка при получении IDs пользователей", nil)
		return
	}

	app.SendBroadcastMessage(telegramIDs, msg, nil)

	session.ClearState()

	HandleAdminCommand(app, session)
}

func handleAdminFeedback(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddBack(states.CallbackAdminSelectBackPanel).Build()

	feedbacks, err := postgres.FetchAllFeedbacks()
	if err != nil || len(feedbacks) == 0 {
		app.SendMessage("📭 Нет новых фидбеков для просмотра.", keyboard)
		return
	}

	const maxMessageLength = 4000

	msg := "📄 <b>Список фидбеков:</b>\n\n"

	for _, feedback := range feedbacks {
		entry := fmt.Sprintf(
			"🆔 ID: %d\n"+
				"👤 Пользователь: tg://user?id=%d\n"+
				"📂 Категория: %s\n"+
				"💬 %s\n"+
				"📅 %s\n\n"+
				"🗑️ /delete_feedback_%d\n"+
				"━━━━━━━━━━━━━\n",
			feedback.ID, feedback.TelegramID, feedback.Category, feedback.Message,
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
		app.SendMessage("Ошибка при получении ID фидбека.", nil)
		return
	}

	err = postgres.DeleteFeedbackByID(feedbackID)
	if err != nil {
		app.SendMessage("❌ Не удалось удалить фидбек.", nil)
		return
	}

	app.SendMessage("✅ Фидбек успешно удален.", nil)

	handleAdminFeedback(app, session)
}

func handleAdminUsers(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddBack(states.CallbackAdminSelectBackPanel).Build()

	users, err := postgres.FetchAllUsers()
	if err != nil || len(users) == 0 {
		app.SendMessage("📭 Нет пользователей для отображения.", keyboard)
		return
	}

	const maxMessageLength = 4000

	msg := "📄 <b>Список пользователей:</b>\n\n"

	for _, user := range users {
		entry := fmt.Sprintf(
			"🆔 Telegram ID: tg://user?id=%d\n"+
				"👤 <b>Админ:</b> %s\n"+
				"🔐 <b>Заблокирован:</b> %s\n"+
				"📅 <b>Создан:</b> %s\n",
			user.TelegramID,
			boolToString(user.IsAdmin),
			boolToString(user.IsBanned),
			user.CreatedAt.Format("02.01.2006 15:04"),
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
		app.SendMessage("❌ Неверный формат ID.", nil)
		return
	}

	err = postgres.BanUser(telegramID)
	if err != nil {
		app.SendMessage("❌ Ошибка при блокировке пользователя.", nil)
		return
	}

	app.SendMessage(fmt.Sprintf("✅ Пользователь %d успешно заблокирован.", telegramID), nil)

	handleAdminUsers(app, session)
}

func handleUnbanUser(app models.App, session *models.Session) {
	telegramID, err := parseUnbanTelegramID(app)
	if err != nil {
		app.SendMessage("❌ Неверный формат ID.", nil)
		return
	}

	err = postgres.UnbanUser(telegramID)
	if err != nil {
		app.SendMessage("❌ Ошибка при разблокировке пользователя.", nil)
		return
	}

	app.SendMessage(fmt.Sprintf("✅ Пользователь %d успешно разблокирован.", telegramID), nil)

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
