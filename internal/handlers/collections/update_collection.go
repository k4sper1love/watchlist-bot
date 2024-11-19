package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

var updateCollectionButtons = []keyboards.Button{
	{"Название", states.CallbackUpdateCollectionSelectName},
	{"Описание", states.CallbackUpdateCollectionSelectDescription},
}

func HandleUpdateCollectionCommand(app models.App, session *models.Session) {
	msg := messages.BuildCollectionDetailMessage(&session.CollectionDetailState.Collection)
	msg += "\nВыберите, какое поле вы хотите изменить?"

	keyboard := keyboards.NewKeyboard().
		AddButtons(updateCollectionButtons...).
		AddBack(states.CallbackUpdateCollectionSelectBack).
		Build()

	app.SendMessage(msg, keyboard)
}

func HandleUpdateCollectionButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackUpdateCollectionSelectBack:
		HandleManageCollectionCommand(app, session)

	case states.CallbackUpdateCollectionSelectName:
		handleUpdateCollectionName(app, session)

	case states.CallbackUpdateCollectionSelectDescription:
		handleUpdateCollectionDescription(app, session)
	}
}

func HandleUpdateCollectionProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearState()
		session.CollectionDetailState.Clear()
		HandleUpdateCollectionCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessUpdateCollectionAwaitingName:
		parseUpdateCollectionName(app, session)
	case states.ProcessUpdateCollectionAwaitingDescription:
		parseUpdateCollectionDescription(app, session)
	}
}

func handleUpdateCollectionName(app models.App, session *models.Session) {
	msg := "Введите новое имя для коллекции"

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionAwaitingName)
}

func parseUpdateCollectionName(app models.App, session *models.Session) {
	session.CollectionDetailState.Name = utils.ParseMessageString(app.Upd)

	finishUpdateCollectionProcess(app, session)
}

func handleUpdateCollectionDescription(app models.App, session *models.Session) {
	msg := "Введите новое описание для коллекции"

	keyboard := keyboards.NewKeyboard().AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionAwaitingDescription)
}

func parseUpdateCollectionDescription(app models.App, session *models.Session) {
	session.CollectionDetailState.Description = utils.ParseMessageString(app.Upd)

	finishUpdateCollectionProcess(app, session)
}

func updateCollection(app models.App, session *models.Session) {
	collection, err := watchlist.UpdateCollection(app, session)
	if err != nil {
		app.SendMessage("Не удалось обновить коллекцию", nil)
		return
	}

	session.CollectionDetailState.Collection = *collection
	app.SendMessage("Коллекция успешно обновлена", nil)
}

func finishUpdateCollectionProcess(app models.App, session *models.Session) {
	updateCollection(app, session)
	session.ClearAllStates()
	HandleUpdateCollectionCommand(app, session)
}
