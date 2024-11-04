package collections

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleNewCollectionCommand(app models.App, session *models.Session) {
	keyboard := builders.NewKeyboard(1).AddCancel().Build()
	app.SendMessage("Введите название коллекции", keyboard)
	session.SetState(states.ProcessNewCollectionAwaitingName)
}

func HandleNewCollectionProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearState()
		HandleCollectionsCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessNewCollectionAwaitingName:
		parseNewCollectionName(app, session)

	case states.ProcessNewCollectionAwaitingDescription:
		parseNewCollectionDescription(app, session)
	}
}

func parseNewCollectionName(app models.App, session *models.Session) {
	session.CollectionDetailState.Name = utils.ParseMessageString(app.Upd)

	keyboard := builders.NewKeyboard(1).AddSkip().AddCancel().Build()

	app.SendMessage("Введите описание коллекции", keyboard)

	session.SetState(states.ProcessNewCollectionAwaitingDescription)
}

func parseNewCollectionDescription(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.CollectionDetailState.Description = ""
	} else {
		session.CollectionDetailState.Description = utils.ParseMessageString(app.Upd)
	}

	createCollection(app, session)
	session.CollectionDetailState.Clear()
	session.ClearState()
}

func createCollection(app models.App, session *models.Session) {
	collection, err := watchlist.CreateCollection(app, session)
	if err != nil {
		app.SendMessage("Не удалось создать коллекцию", nil)
		return
	}

	app.SendMessage("Новая коллекция успешно создана!", nil)
	msg := fmt.Sprintf("ID: %d\n", collection.ID) +
		fmt.Sprintf("Name: %s\n", collection.Name) +
		fmt.Sprintf("Description: %s\n", collection.Description) +
		fmt.Sprintf("Last updated: %s", collection.UpdatedAt.String()) +
		fmt.Sprintf("Created: %s\n", collection.CreatedAt.String())
	app.SendMessage(msg, nil)
}
