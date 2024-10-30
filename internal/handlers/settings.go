package handlers

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"strconv"
)

var settingsButtons = []builders.Button{
	{"Изменить количество коллекций на странице", "settings_collections_page_size"},
}

func handleSettingCommand(app config.App, session *models.Session) {
	msg := "Выберите, что хотите настроить\n"
	keyboard := builders.BuildButtonKeyboard(settingsButtons, 1)
	sendMessageWithKeyboard(app, msg, keyboard)
}

func handleSettingButton(app config.App, session *models.Session) {
	switch session.State {
	case CallbackSettingsCollectionsPageSize:
		msg := fmt.Sprintf("Текущий размер страницы для списка коллекций: %d\nВведите желаемый размер (от 1 и выше)", session.CollectionState.PageSize)
		sendMessage(app, msg)

		setState(session, ProcessSettingsCollectionsAwaitingPageSize)
	}
}

func handleSettingProcess(app config.App, session *models.Session) {
	switch session.State {
	case ProcessSettingsCollectionsAwaitingPageSize:
		pageSize, err := strconv.Atoi(parseMessageText(app.Upd))
		if err != nil || pageSize < 1 {
			sendMessage(app, "Введите целое число от 1 и выше")
			return
		}

		session.CollectionState.PageSize = pageSize
		sendMessage(app, fmt.Sprintf("Новый размер страницы для списка коллекций: %d", pageSize))

		resetState(session)
	}
}
