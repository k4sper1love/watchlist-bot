package collectionFilms

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"log/slog"
	"strconv"
	"strings"
)

func HandleAddCollectionToFilmCommand(app models.App, session *models.Session) {
	collections, err := GetCollectionsExcludeFilm(app, session)
	if err != nil {
		app.SendMessage(err.Error(), nil)
		return
	}

	if len(collections) == 0 {
		msg := "❗️" + translator.Translate(session.Lang, "collectionsNotFound", nil, nil)
		keyboard := keyboards.BuildAddCollectionToFilmNotFoundKeyboard(session)
		app.SendMessage(msg, keyboard)
		return
	}

	choiceMsg := translator.Translate(session.Lang, "choiceCollection", nil, nil)
	msg := fmt.Sprintf("<b>%s</b>", choiceMsg)
	keyboard := keyboards.BuildAddCollectionToFilmKeyboard(session)
	app.SendMessage(msg, keyboard)
}

func HandleAddCollectionToFilmButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	switch {
	case callback == states.CallbackAddCollectionToFilmBack:
		films.HandleFilmsDetailCommand(app, session)

	case strings.HasPrefix(callback, "select_cf_collection_"):
		HandleAddCollectionToFilmSelect(app, session)

	case callback == states.CallbackAddCollectionToFilmFind:
		handleAddCollectionToFilmFind(app, session)

	case callback == states.CallbackAddCollectionToFilmAgain:
		session.CollectionsState.Name = ""
		handleAddCollectionToFilmFind(app, session)

	case callback == states.CallbackAddCollectionToFilmReset:
		session.CollectionsState.Name = ""
		HandleAddCollectionToFilmCommand(app, session)

	case callback == states.CallbackAddCollectionToFilmNextPage:
		if session.CollectionFilmsState.CurrentPage < session.CollectionFilmsState.LastPage {
			session.CollectionFilmsState.CurrentPage++
			HandleAddCollectionToFilmCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAddCollectionToFilmPrevPage:
		if session.CollectionFilmsState.CurrentPage > 1 {
			session.CollectionFilmsState.CurrentPage--
			HandleAddCollectionToFilmCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAddCollectionToFilmLastPage:
		if session.CollectionFilmsState.CurrentPage != session.CollectionFilmsState.LastPage {
			session.CollectionFilmsState.CurrentPage = session.CollectionFilmsState.LastPage
			HandleAddCollectionToFilmCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackAddCollectionToFilmFirstPage:
		if session.CollectionFilmsState.CurrentPage != 1 {
			session.CollectionFilmsState.CurrentPage = 1
			HandleAddCollectionToFilmCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}
	}
}

func HandleAddCollectionToFilmProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearAllStates()
		HandleAddCollectionToFilmCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessAddCollectionToFilmAwaitingName:
		parseAddCollectionToFilmName(app, session)
	}
}

func HandleAddCollectionToFilmSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	idStr := strings.TrimPrefix(callback, "select_cf_collection_")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "getCollectionFailure", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackCollectionFilmsFromFilm).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		sl.Log.Error("failed to parse collection ID", slog.Any("error", err), slog.String("callback", callback))
		return
	}

	session.CollectionDetailState.Collection.ID = id

	addFilmToCollection(app, session)
}

func handleAddCollectionToFilmFind(app models.App, session *models.Session) {
	msg := "❓ " + translator.Translate(session.Lang, "collectionRequestName", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessAddCollectionToFilmAwaitingName)
}

func parseAddCollectionToFilmName(app models.App, session *models.Session) {
	name := utils.ParseMessageString(app.Upd)

	session.CollectionsState.Name = name
	session.CollectionFilmsState.CurrentPage = 1

	session.ClearState()

	HandleAddCollectionToFilmCommand(app, session)
}

func GetCollectionsExcludeFilm(app models.App, session *models.Session) ([]apiModels.Collection, error) {
	collectionsResponse, err := watchlist.GetCollectionsExcludeFilm(app, session)
	if err != nil {
		return nil, err
	}

	session.CollectionsState.Collections = collectionsResponse.Collections
	session.CollectionFilmsState.CurrentPage = collectionsResponse.Metadata.CurrentPage
	session.CollectionFilmsState.LastPage = collectionsResponse.Metadata.LastPage

	return collectionsResponse.Collections, nil
}
