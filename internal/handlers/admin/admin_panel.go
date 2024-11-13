package admin

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/database/postgres"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

var adminButtons = []builders.Button{
	{"Узнать количество пользателей", states.CallbackAdminSelectUserCount},
	{"Отправить рассылку", states.CallbackAdminSelectBroadcastMessage},
}

func HandleAdminCommand(app models.App, session *models.Session) {
	msg := "Выберите действие"

	keyboard := builders.NewKeyboard(1).
		AddSeveral(adminButtons).
		AddBack(states.CallbackAdminSelectBack).
		Build()

	app.SendMessage(msg, keyboard)
}

func HandleAdminButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackAdminSelectBack:
		general.HandleMenuCommand(app, session)

	case states.CallbackAdminSelectUserCount:
		handleAdminUserCount(app, session)

	case states.CallbackAdminSelectBroadcastMessage:
		handleAdminBroadcastMessage(app, session)
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

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

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
