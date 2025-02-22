package films

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleManageFilmCommand(app models.App, session *models.Session) {
	msg := messages.BuildFilmDetailMessage(session)
	choiceMsg := translator.Translate(session.Lang, "choiceAction", nil, nil)
	msg += fmt.Sprintf("<b>%s</b>", choiceMsg)

	keyboard := keyboards.BuildFilmManageKeyboard(session)

	app.SendImage(session.FilmDetailState.Film.ImageURL, msg, keyboard)
}

func HandleManageFilmButtons(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Update) {
	case states.CallbackManageFilmSelectBack:
		HandleFilmsDetailCommand(app, session)

	case states.CallbackManageFilmSelectUpdate:
		HandleUpdateFilmCommand(app, session)

	case states.CallbackManageFilmSelectDelete:
		HandleDeleteFilmCommand(app, session)

	case states.CallbackManageFilmSelectRemoveFromCollection:
		handleRemoveFilmFromCollection(app, session)
	}
}

func handleRemoveFilmFromCollection(app models.App, session *models.Session) {
	if err := watchlist.DeleteCollectionFilm(app, session); err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "removeFilmFailure", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackFilmsManage).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		return
	}

	msg := "🧹󠁝 " + translator.Translate(session.Lang, "removeFilmSuccess", nil, nil)
	app.SendMessage(msg, nil)

	HandleFilmsCommand(app, session)
}
