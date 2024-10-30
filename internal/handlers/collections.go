package handlers

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"log"
	"strconv"
	"strings"
)

func handleCollectionsCommand(app config.App, session *models.Session) {
	collectionsResponse, err := watchlist.GetCollections(app, session)
	if err != nil {
		sendMessage(app, err.Error())
		return
	}

	if collectionsResponse.Metadata.TotalRecords == 0 {
		sendMessage(app, "Не найдено коллекций. Создайте новую коллекцию с помощью /new_collection")
		return
	}

	currentPage := collectionsResponse.Metadata.CurrentPage
	lastPage := collectionsResponse.Metadata.LastPage

	session.CollectionState.CurrentPage = currentPage
	session.CollectionState.LastPage = lastPage

	msg := builders.BuildCollectionsMessage(collectionsResponse)

	buttons := builders.BuildCollectionsSelectButtons(collectionsResponse)
	buttons = append(buttons, builders.BuildNavigationButtons(currentPage, lastPage, CallbackCollectionsPrevPage, CallbackCollectionsNextPage)...)

	keyboard := builders.BuildButtonKeyboard(buttons, 1)

	sendMessageWithKeyboard(app, msg, keyboard)
}

func handleCollectionsButtons(app config.App, session *models.Session) {
	switch {
	case session.State == CallbackCollectionsNextPage:
		if session.CollectionState.CurrentPage < session.CollectionState.LastPage {
			session.CollectionState.CurrentPage++
			handleCollectionsCommand(app, session)
		} else {
			sendMessage(app, "Вы уже на последней странице")
		}
		resetState(session)

	case session.State == CallbackCollectionsPrevPage:
		if session.CollectionState.CurrentPage > 1 {
			session.CollectionState.CurrentPage--
			handleCollectionsCommand(app, session)
		} else {
			sendMessage(app, "Вы уже на первой странице")
		}
		resetState(session)

	case strings.HasPrefix(session.State, "select_collection_"):
		handleCollectionSelect(app, session)
		resetState(session)
	}
}

func handleCollectionSelect(app config.App, session *models.Session) {
	collectionIDStr := strings.TrimPrefix(session.State, "select_collection_")
	collectionID, err := strconv.Atoi(collectionIDStr)

	if err != nil {
		sendMessage(app, "Ошибка при получении ID коллекции.")
		log.Printf("error parsing collection ID: %v", err)
		return
	}

	session.CollectionState.ObjectID = collectionID
	session.CollectionFilmState.CurrentPage = 1
	handleCollectionFilmsCommand(app, session)
}

func handleNewCollectionCommand(app config.App, session *models.Session) {
	sendMessage(app, "Введите название коллекции")
	setState(session, ProcessNewCollectionAwaitingName)
}

func handleNewCollectionProcess(app config.App, session *models.Session) {
	switch session.State {
	case ProcessNewCollectionAwaitingName:
		session.CollectionState.Name = parseMessageText(app.Upd)
		sendMessage(app, "Введите описание коллекции (-, если хотите оставить пустым)")
		setState(session, ProcessNewCollectionAwaitingDescription)

	case ProcessNewCollectionAwaitingDescription:
		if app.Upd.Message.Text == "-" {
			session.CollectionState.Description = ""
		} else {
			session.CollectionState.Description = parseMessageText(app.Upd)
		}

		collection, err := watchlist.CreateCollection(app, session)
		fmt.Println(err)
		if err != nil {
			sendMessage(app, "Не удалось создать коллекцию")
		} else {
			sendMessage(app, "Новая коллекция успешно создана!")
			msg := fmt.Sprintf("ID: %d\n", collection.ID) +
				fmt.Sprintf("Name: %s\n", collection.Name) +
				fmt.Sprintf("Description: %s\n", collection.Description) +
				fmt.Sprintf("Last updated: %s", collection.UpdatedAt.String()) +
				fmt.Sprintf("Created: %s\n", collection.CreatedAt.String())
			sendMessage(app, msg)
		}

		resetState(session)
	}
}
