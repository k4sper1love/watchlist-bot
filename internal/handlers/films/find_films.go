package films

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleFindFilmsCommand(app models.App, session *models.Session) {
	metadata, err := GetFilms(app, session)
	if err != nil {
		app.SendMessage(err.Error(), nil)
		return
	}

	if metadata.TotalRecords == 0 {
		msg := translator.Translate(session.Lang, "filmsNotFound", nil, nil)
		keyboard := keyboards.NewKeyboard().AddAgain(states.CallbackFindFilmsAgain).AddBack(states.CallbackFindFilmsBack).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		return
	}

	msg := messages.BuildFindFilmsMessage(session, metadata)

	keyboard := keyboards.BuildFindFilmsKeyboard(session, metadata.CurrentPage, metadata.LastPage)

	app.SendMessage(msg, keyboard)
}

func HandleFindFilmsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)

	switch {
	case callback == states.CallbackFindFilmsBack:
		session.ClearAllStates()
		HandleFilmsCommand(app, session)
		return

	case callback == states.CallbackFindFilmsAgain:
		session.ClearAllStates()
		handleFilmsFindByTitle(app, session)
		return

	case callback == states.CallbackFindFilmsNextPage:
		if session.FilmsState.CurrentPage < session.FilmsState.LastPage {
			session.FilmsState.CurrentPage++
			HandleFindFilmsCommand(app, session)
		} else {
			msg := translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackFindFilmsPrevPage:
		if session.FilmsState.CurrentPage > 1 {
			session.FilmsState.CurrentPage--
			HandleFindFilmsCommand(app, session)
		} else {
			msg := translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}
	}
}
