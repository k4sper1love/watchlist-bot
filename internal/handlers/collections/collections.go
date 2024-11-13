package collections

import (
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"log"
	"strconv"
	"strings"
)

func HandleCollectionsCommand(app models.App, session *models.Session) {
	collectionsResponse, err := watchlist.GetCollections(app, session)
	if err != nil {
		app.SendMessage(err.Error(), nil)
		return
	}

	session.CollectionsState.Object = collectionsResponse.Collections

	currentPage := collectionsResponse.Metadata.CurrentPage
	lastPage := collectionsResponse.Metadata.LastPage

	session.CollectionsState.CurrentPage = currentPage
	session.CollectionsState.LastPage = lastPage

	msg := builders.BuildCollectionsMessage(collectionsResponse)

	if collectionsResponse.Metadata.TotalRecords == 0 {
		msg = "Не найдено коллекций"
	}

	keyboard := builders.NewKeyboard(1).
		AddCollectionsSelect(collectionsResponse).
		AddCollectionsNew().
		AddNavigation(currentPage, lastPage, states.CallbackCollectionsPrevPage, states.CallbackCollectionsNextPage).
		AddBack(states.CallbackCollectionsBack).
		Build()

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
			app.SendMessage("Вы уже на последней странице", nil)
		}

	case callback == states.CallbackCollectionsPrevPage:
		if session.CollectionsState.CurrentPage > 1 {
			session.CollectionsState.CurrentPage--
			HandleCollectionsCommand(app, session)
		} else {
			app.SendMessage("Вы уже на первой странице", nil)
		}

	case callback == states.CallbackCollectionsNew:
		HandleNewCollectionCommand(app, session)

	case callback == states.CallbackCollectionsManage:
		HandleManageCollectionCommand(app, session)

	case strings.HasPrefix(callback, "select_collection_"):
		HandleCollectionSelect(app, session)
	}
}

func HandleCollectionSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	collectionIDStr := strings.TrimPrefix(callback, "select_collection_")
	collectionID, err := strconv.Atoi(collectionIDStr)

	if err != nil {
		app.SendMessage("Ошибка при получении ID коллекции.", nil)
		log.Printf("error parsing collection ID: %v", err)
		return
	}

	session.CollectionDetailState.ObjectID = collectionID
	session.CollectionDetailState.CurrentPage = 1
	HandleCollectionFilmsCommand(app, session)
}
