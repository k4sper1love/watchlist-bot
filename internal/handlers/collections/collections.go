package collections

import (
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"log/slog"
	"strconv"
	"strings"
)

func HandleCollectionsCommand(app models.App, session *models.Session) {
	session.CollectionsState.Name = ""

	metadata, err := GetCollections(app, session)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "getCollectionsFailure", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack("").Build(session.Lang)
		app.SendMessage(msg, keyboard)
		return
	}

	msg := messages.BuildCollectionsMessage(session, metadata, false)

	keyboard := keyboards.BuildCollectionsKeyboard(session, metadata.CurrentPage, metadata.LastPage)

	app.SendMessage(msg, keyboard)
}

func HandleCollectionsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	switch {
	case callback == states.CallbackCollectionsBack:
		general.HandleMenuCommand(app, session)

	case callback == states.CallbackCollectionsNextPage:
		if session.CollectionsState.CurrentPage < session.CollectionsState.LastPage {
			session.CollectionsState.CurrentPage++
			HandleCollectionsCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackCollectionsPrevPage:
		if session.CollectionsState.CurrentPage > 1 {
			session.CollectionsState.CurrentPage--
			HandleCollectionsCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackCollectionsLastPage:
		if session.CollectionsState.CurrentPage != session.CollectionsState.LastPage {
			session.CollectionsState.CurrentPage = session.CollectionsState.LastPage
			HandleCollectionsCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackCollectionsFirstPage:
		if session.CollectionsState.CurrentPage != 1 {
			session.CollectionsState.CurrentPage = 1
			HandleCollectionsCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackCollectionsNew:
		HandleNewCollectionCommand(app, session)

	case callback == states.CallbackCollectionsManage:
		HandleManageCollectionCommand(app, session)

	case callback == states.CallbackCollectionsFind:
		handleCollectionsFindByName(app, session)

	case callback == states.CallbackCollectionsSorting:
		HandleSortingCollectionsCommand(app, session)

	case callback == states.CallbackCollectionsFavorite:
		handleFavoriteCollection(app, session)

	case strings.HasPrefix(callback, "select_collection_"):
		HandleCollectionSelect(app, session)
	}
}

func HandleCollectionProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearAllStates()
		HandleCollectionsCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessFindCollectionsAwaitingName:
		parseCollectionsFindName(app, session)
	}
}

func HandleCollectionSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	idStr := strings.TrimPrefix(callback, "select_collection_")
	id, err := strconv.Atoi(idStr)

	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "getCollectionFailure", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack("").Build(session.Lang)
		app.SendMessage(msg, keyboard)
		sl.Log.Error("failed to parse collection ID", slog.Any("error", err), slog.String("callback", callback))
		return
	}

	session.CollectionDetailState.ObjectID = id
	session.FilmsState.CurrentPage = 1

	session.SetContext(states.ContextCollection)

	films.HandleFilmsCommand(app, session)
}

func handleCollectionsFindByName(app models.App, session *models.Session) {
	msg := "❓" + translator.Translate(session.Lang, "collectionRequestName", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessFindCollectionsAwaitingName)
}

func parseCollectionsFindName(app models.App, session *models.Session) {
	name := utils.ParseMessageString(app.Upd)

	session.CollectionsState.Name = name
	session.CollectionsState.CurrentPage = 1

	session.ClearState()

	HandleFindCollectionsCommand(app, session)
}

func handleFavoriteCollection(app models.App, session *models.Session) {
	session.CollectionDetailState.Collection.IsFavorite = !session.CollectionDetailState.Collection.IsFavorite

	updateCollection(app, session)

	session.ClearAllStates()

	films.HandleFilmsCommand(app, session)
}

func GetCollections(app models.App, session *models.Session) (*filters.Metadata, error) {
	collectionsResponse, err := watchlist.GetCollections(app, session)
	if err != nil {
		return nil, err
	}

	UpdateSessionWithCollections(session, collectionsResponse.Collections, &collectionsResponse.Metadata)

	return &collectionsResponse.Metadata, nil
}

func UpdateSessionWithCollections(session *models.Session, collections []apiModels.Collection, metadata *filters.Metadata) {
	session.CollectionsState.Collections = collections
	session.CollectionsState.LastPage = metadata.LastPage
	//session.CollectionsState.LastPage = metadata.LastPage
}
