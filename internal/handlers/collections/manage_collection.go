package collections

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleManageCollectionCommand(app models.App, session *models.Session) {
	msg := messages.BuildCollectionHeader(session)
	choiceMsg := translator.Translate(session.Lang, "choiceAction", nil, nil)
	msg += fmt.Sprintf("<b>%s</b>", choiceMsg)

	keyboard := keyboards.BuildCollectionManageKeyboard(session)

	app.SendMessage(msg, keyboard)
}

func HandleManageCollectionButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackManageCollectionSelectBack:
		session.SetContext(states.ContextCollection)
		films.HandleFilmsCommand(app, session)

	case states.CallbackManageCollectionSelectUpdate:
		HandleUpdateCollectionCommand(app, session)

	case states.CallbackManageCollectionSelectDelete:
		HandleDeleteCollectionCommand(app, session)
	}
}
