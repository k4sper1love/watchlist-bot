package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleNewCollectionCommand(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	msg := translator.Translate(session.Lang, "collectionRequestName", nil, nil)

	app.SendMessage(msg, keyboard)
	session.SetState(states.ProcessNewCollectionAwaitingName)
}

func HandleNewCollectionProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearState()
		session.CollectionDetailState.Clear()
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

	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := translator.Translate(session.Lang, "collectionRequestDescription", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewCollectionAwaitingDescription)
}

func parseNewCollectionDescription(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.CollectionDetailState.Description = ""
	} else {
		session.CollectionDetailState.Description = utils.ParseMessageString(app.Upd)
	}

	createCollection(app, session)
	session.ClearAllStates()
}

func createCollection(app models.App, session *models.Session) {
	collection, err := watchlist.CreateCollection(app, session)
	if err != nil {
		msg := translator.Translate(session.Lang, "createCollectionFailure", nil, nil)
		app.SendMessage(msg, nil)
		return
	}

	msg := translator.Translate(session.Lang, "createCollectionSuccess", nil, nil)
	app.SendMessage(msg, nil)

	session.CollectionDetailState.ObjectID = collection.ID

	session.SetContext(states.ContextCollection)
	films.HandleFilmsCommand(app, session)
}
