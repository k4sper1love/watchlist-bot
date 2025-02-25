package collections

import (
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/films"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strconv"
	"strings"
)

func HandleCollectionsCommand(app models.App, session *models.Session) {
	session.CollectionsState.Clear()

	if metadata, err := getCollections(app, session); err != nil {
		app.SendMessage(messages.BuildCollectionsFailureMessage(session), keyboards.BuildKeyboardWithBack(session, ""))
	} else {
		app.SendMessage(messages.BuildCollectionsMessage(session, metadata, false), keyboards.BuildCollectionsKeyboard(session, metadata.CurrentPage, metadata.LastPage))
	}
}

func HandleCollectionsButtons(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallbackCollectionsBack:
		general.HandleMenuCommand(app, session)

	case states.CallbackCollectionsNextPage, states.CallbackCollectionsPrevPage,
		states.CallbackCollectionsLastPage, states.CallbackCollectionsFirstPage:
		handleCollectionsPagination(app, session, callback)

	case states.CallbackCollectionsNew:
		HandleNewCollectionCommand(app, session)

	case states.CallbackCollectionsManage:
		HandleManageCollectionCommand(app, session)

	case states.CallbackCollectionsFind:
		handleCollectionsFindByName(app, session)

	case states.CallbackCollectionsSorting:
		HandleSortingCollectionsCommand(app, session)

	case states.CallbackCollectionsFavorite:
		handleFavoriteCollection(app, session)

	default:
		if strings.HasPrefix(callback, states.PrefixSelectCollection) {
			handleCollectionSelect(app, session, callback)
		}
	}
}

func HandleCollectionProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleCollectionsCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessFindCollectionsAwaitingName:
		parseCollectionFindName(app, session, HandleFindCollectionsCommand)
	}
}

func handleCollectionsPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallbackCollectionsNextPage:
		if session.CollectionsState.CurrentPage >= session.CollectionsState.LastPage {
			app.SendMessage(messages.BuildLastPageAlertMessage(session), nil)
			return
		}
		session.CollectionsState.CurrentPage++

	case states.CallbackCollectionsPrevPage:
		if session.CollectionsState.CurrentPage <= 1 {
			app.SendMessage(messages.BuildFirstPageAlertMessage(session), nil)
			return
		}
		session.CollectionsState.CurrentPage--

	case states.CallbackCollectionsLastPage:
		if session.CollectionsState.CurrentPage == session.CollectionsState.LastPage {
			app.SendMessage(messages.BuildLastPageAlertMessage(session), nil)
			return
		}
		session.CollectionsState.CurrentPage = session.CollectionsState.LastPage

	case states.CallbackCollectionsFirstPage:
		if session.CollectionsState.CurrentPage == 1 {
			app.SendMessage(messages.BuildFirstPageAlertMessage(session), nil)
			return
		}
		session.CollectionsState.CurrentPage = 1
	}

	HandleCollectionsCommand(app, session)
}

func handleCollectionSelect(app models.App, session *models.Session, callback string) {
	if id, err := strconv.Atoi(strings.TrimPrefix(callback, states.PrefixSelectCollection)); err != nil {
		utils.LogParseSelectError(err, callback)
		app.SendMessage(messages.BuildCollectionsFailureMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackCollectionsBack))
	} else {
		session.CollectionDetailState.ObjectID = id
		setContextAndHandleFilms(app, session)
	}
}

func handleCollectionsFindByName(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildCollectionRequestNameMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessFindCollectionsAwaitingName)
}

func handleFavoriteCollection(app models.App, session *models.Session) {
	session.CollectionDetailState.SetFavorite(!session.CollectionDetailState.Collection.IsFavorite)
	HandleUpdateCollection(app, session, films.HandleFilmsCommand)
}

func getCollections(app models.App, session *models.Session) (*filters.Metadata, error) {
	collectionsResponse, err := watchlist.GetCollections(app, session)
	if err != nil {
		return nil, err
	}

	updateSessionWithCollections(session, collectionsResponse.Collections, &collectionsResponse.Metadata)
	return &collectionsResponse.Metadata, nil
}

func updateSessionWithCollections(session *models.Session, collections []apiModels.Collection, metadata *filters.Metadata) {
	session.CollectionsState.Collections = collections
	session.CollectionsState.LastPage = metadata.LastPage
}

func setContextAndHandleFilms(app models.App, session *models.Session) {
	session.SetContext(states.ContextCollection)
	session.FilmsState.CurrentPage = 1
	films.HandleFilmsCommand(app, session)
}
