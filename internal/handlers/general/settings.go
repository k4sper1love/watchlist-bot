package general

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strconv"
)

func HandleSettingsCommand(app models.App, session *models.Session) {
	msg := "⚙️ <b>Настройки:</b>\n\nВыберите параметр для изменения"

	keyboard := keyboards.BuildSettingsKeyboard()

	app.SendMessage(msg, keyboard)
}

func HandleSettingsButton(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackSettingsCollectionsPageSize:
		handleCollectionsPageSize(app, session)
	case states.CallbackSettingsFilmsPageSize:
		handleFilmsPageSize(app, session)
	case states.CallbackSettingsObjectsPageSize:
		handleObjectsPageSize(app, session)
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
	case states.ProcessSettingsFilmsAwaitingPageSize:
		parseFilmsPageSize(app, session)
	case states.ProcessSettingsObjectsAwaitingPageSize:
		parseObjectsPageSize(app, session)
	}
}

func handleCollectionsPageSize(app models.App, session *models.Session) {
	msg := fmt.Sprintf("Текущий размер страницы для списка коллекций: %d\nВведите желаемый размер (от 1 и выше)", session.CollectionsState.PageSize)

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

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

func handleFilmsPageSize(app models.App, session *models.Session) {
	msg := fmt.Sprintf("Текущий размер страницы для списка фильмов: %d\nВведите желаемый размер (от 1 и выше)", session.FilmsState.PageSize)

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessSettingsFilmsAwaitingPageSize)
}

func parseFilmsPageSize(app models.App, session *models.Session) {
	pageSize, err := strconv.Atoi(utils.ParseMessageString(app.Upd))
	if err != nil || pageSize < 1 {
		app.SendMessage("Введите целое число от 1 и выше", nil)
		return
	}

	session.FilmsState.PageSize = pageSize

	msg := fmt.Sprintf("Новый размер страницы для списка фильмов: %d", pageSize)

	app.SendMessage(msg, nil)

	session.ClearState()

	HandleSettingsCommand(app, session)
}

func handleObjectsPageSize(app models.App, session *models.Session) {
	msg := fmt.Sprintf("Текущий размер страницы для списка объектов: %d\nВведите желаемый размер (от 1 и выше)", session.CollectionFilmsState.PageSize)

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessSettingsObjectsAwaitingPageSize)
}

func parseObjectsPageSize(app models.App, session *models.Session) {
	pageSize, err := strconv.Atoi(utils.ParseMessageString(app.Upd))
	if err != nil || pageSize < 1 {
		app.SendMessage("Введите целое число от 1 и выше", nil)
		return
	}

	session.CollectionFilmsState.PageSize = pageSize

	msg := fmt.Sprintf("Новый размер страницы для списка коллекций: %d", pageSize)

	app.SendMessage(msg, nil)

	session.ClearState()

	HandleSettingsCommand(app, session)
}
