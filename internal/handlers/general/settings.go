package general

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strconv"
)

var settingsButtons = []builders.Button{
	{"Изменить количество коллекций на странице", states.CallbackSettingsCollectionsPageSize},
}

func HandleSettingsCommand(app models.App, session *models.Session) {
	msg := "Выберите, что хотите настроить\n"
	keyboard := builders.NewKeyboard(1).AddSeveral(settingsButtons).AddBack(states.CallbackSettingsBack).Build()
	app.SendMessage(msg, keyboard)
}

func HandleSettingsButton(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackSettingsBack:
		HandleMenuCommand(app, session)

	case states.CallbackSettingsCollectionsPageSize:
		handleCollectionsPageSize(app, session)
	}
}

func HandleSettingsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearState()
		HandleSettingsCommand(app, session)
	}

	switch session.State {
	case states.ProcessSettingsCollectionsAwaitingPageSize:
		parseCollectionsPageSize(app, session)
	}
}

func handleCollectionsPageSize(app models.App, session *models.Session) {
	msg := fmt.Sprintf("Текущий размер страницы для списка коллекций: %d\nВведите желаемый размер (от 1 и выше)", session.CollectionsState.PageSize)

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessSettingsCollectionsAwaitingPageSize)
}

func parseCollectionsPageSize(app models.App, session *models.Session) {
	pageSize, err := strconv.Atoi(utils.ParseMessageString(app.Upd))
	if err != nil || pageSize < 1 {
		app.SendMessage("Введите целое число от 1 и выше", nil)
		return
	}

	session.CollectionsState.PageSize = pageSize

	msg := fmt.Sprintf("Новый размер страницы для списка коллекций: %d", pageSize)

	app.SendMessage(msg, nil)

	session.ClearState()

	HandleSettingsCommand(app, session)
}
