package collections

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleCollectionFilmsDetailCommand(app models.App, session *models.Session) {
	index := session.CollectionFilmState.Index
	films := session.CollectionDetailState.Object.Films
	film := films[index]

	if len(films) == 0 {
		app.SendMessage("Фильмы не найдены. Начните с начала", nil)
		return
	}

	if index == -1 || index >= len(films) {
		app.SendMessage("Неизвестный фильм в коллекции. Начните с начала", nil)
		return
	}

	session.CollectionFilmState.Object = films[index]

	itemID := utils.GetItemID(index, session.CollectionDetailState.CurrentPage, session.CollectionDetailState.PageSize)

	msg := fmt.Sprintf("🎞️ <b>№</b>: %d\n", itemID)
	msg += builders.BuildCollectionFilmDetailMessage(&films[index])

	var buttons []builders.Button

	if !film.IsViewed {
		buttons = append(buttons, builders.Button{"Просмотрено✅", states.CallbackCollectionFilmDetailViewed})
	}

	keyboard := builders.NewKeyboard(2).
		AddSeveral(buttons).
		AddCollectionFilmsManage().
		AddNavigation(itemID, session.CollectionDetailState.TotalRecords, states.CallbackCollectionFilmDetailPrevPage, states.CallbackCollectionFilmDetailNextPage).
		AddBack(states.CallbackCollectionFilmDetailBack).
		Build()

	imageURL := films[index].ImageURL

	app.SendImage(imageURL, msg, keyboard)
}

func HandleCollectionFilmsDetailButtons(app models.App, session *models.Session) {
	currentIndex := session.CollectionFilmState.Index
	lastIndex := getCollectionFilmsLastIndex(session)

	switch utils.ParseCallback(app.Upd) {
	case states.CallbackCollectionFilmDetailNextPage:
		if currentIndex < lastIndex {
			session.CollectionFilmState.Index++
			HandleCollectionFilmsDetailCommand(app, session)
		} else {
			if err := updateCollectionFilmsList(app, session, true); err != nil {
				app.SendMessage(err.Error(), nil)
				return
			}
			session.CollectionFilmState.Index = 0
			HandleCollectionFilmsDetailCommand(app, session)
		}

	case states.CallbackCollectionFilmDetailPrevPage:
		if currentIndex > 0 {
			session.CollectionFilmState.Index--
			HandleCollectionFilmsDetailCommand(app, session)
		} else {
			if err := updateCollectionFilmsList(app, session, false); err != nil {
				app.SendMessage(err.Error(), nil)
				return
			}

			session.CollectionFilmState.Index = getCollectionFilmsLastIndex(session)
			HandleCollectionFilmsDetailCommand(app, session)
		}

	case states.CallbackCollectionFilmDetailBack:
		HandleCollectionFilmsCommand(app, session)

	case states.CallbackCollectionFilmDetailViewed:
		HandleViewedCollectionFilmCommand(app, session)
	}
}

func getCollectionFilmsLastIndex(session *models.Session) int {
	return len(session.CollectionDetailState.Object.Films) - 1
}
