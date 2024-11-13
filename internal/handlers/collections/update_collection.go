package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

var updateCollectionButtons = []builders.Button{
	{"Название", states.CallbackUpdateCollectionSelectName},
	{"Описание", states.CallbackUpdateCollectionSelectDescription},
}

func HandleUpdateCollectionCommand(app models.App, session *models.Session) {
	msg := builders.BuildCollectionDetailMessage(-1, &session.CollectionDetailState.Object.Collection)
	msg += "\nВыберите, какое поле вы хотите изменить?"

	keyboard := builders.NewKeyboard(1).
		AddSeveral(updateCollectionButtons).
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

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessUpdateCollectionAwaitingName)
}

func parseUpdateCollectionName(app models.App, session *models.Session) {
	session.CollectionDetailState.Name = utils.ParseMessageString(app.Upd)

	finishUpdateCollectionProcess(app, session)
}

func handleUpdateCollectionDescription(app models.App, session *models.Session) {
	msg := "Введите новое описание для коллекции"

	keyboard := builders.NewKeyboard(1).AddCancel().Build()

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

	session.CollectionDetailState.Object.Collection = *collection
	app.SendMessage("Коллекция успешно обновлена", nil)
}

func finishUpdateCollectionProcess(app models.App, session *models.Session) {
	state := session.CollectionDetailState
	collection := state.Object.Collection

	if state.Name == "" {
		state.Name = collection.Name
	}

	if state.Description == "" {
		state.Description = collection.Description
	}

	updateCollection(app, session)
	session.CollectionDetailState.Clear()
	session.ClearState()
	HandleUpdateCollectionCommand(app, session)
}
