package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/validator"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleNewCollectionCommand(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "collectionRequestName", nil, nil)

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
	name := utils.ParseMessageString(app.Upd)
	if ok := utils.ValidStringLength(name, 3, 100); !ok {
		validator.HandleInvalidInputLength(app, session, 3, 100)
		HandleNewCollectionCommand(app, session)
		return
	}
	session.CollectionDetailState.Name = name

	requestNewCollectionDescription(app, session)
}

func requestNewCollectionDescription(app models.App, session *models.Session) {
	keyboard := keyboards.NewKeyboard().AddSkip().AddCancel().Build(session.Lang)

	msg := "❓" + translator.Translate(session.Lang, "collectionRequestDescription", nil, nil)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessNewCollectionAwaitingDescription)
}

func parseNewCollectionDescription(app models.App, session *models.Session) {
	if utils.IsSkip(app.Upd) {
		session.CollectionDetailState.Description = ""
		createCollection(app, session)
		return
	}
	description := utils.ParseMessageString(app.Upd)
	if ok := utils.ValidStringLength(description, 0, 500); !ok {
		validator.HandleInvalidInputLength(app, session, 0, 500)
		requestNewCollectionDescription(app, session)
		return
	}
	session.CollectionDetailState.Description = description

	createCollection(app, session)
}

func createCollection(app models.App, session *models.Session) {
	collection, err := watchlist.CreateCollection(app, session)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "createCollectionFailure", nil, nil)
		app.SendMessage(msg, nil)
		return
	}

	msg := "📚 " + translator.Translate(session.Lang, "createCollectionSuccess", nil, nil)
	app.SendMessage(msg, nil)

	session.CollectionDetailState.ObjectID = collection.ID

	session.SetContext(states.ContextCollection)
	films.HandleFilmsCommand(app, session)

	session.ClearAllStates()
}
