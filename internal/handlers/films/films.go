package films

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
	"log/slog"
	"strconv"
	"strings"
)

func HandleFilmsCommand(app models.App, session *models.Session) {
	session.FilmsState.Title = ""

	metadata, err := GetFilms(app, session)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "getFilmsFailure", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack("").Build(session.Lang)
		app.SendMessage(msg, keyboard)
		return
	}

	msg := messages.BuildFilmsMessage(session, metadata)

	keyboard := keyboards.BuildFilmsKeyboard(session, metadata.CurrentPage, metadata.LastPage)

	app.SendMessage(msg, keyboard)
}

func HandleFilmsButtons(app models.App,
	session *models.Session,
	backFunc func(models.App, *models.Session),
) {
	callback := utils.ParseCallback(app.Upd)

	switch {
	case callback == states.CallbackFilmsBack:
		if session.Context == states.ContextCollection {
			session.CollectionsState.CurrentPage = 1
		}
		backFunc(app, session)

	case callback == states.CallbackFilmsNextPage:
		if session.FilmsState.CurrentPage < session.FilmsState.LastPage {
			session.FilmsState.CurrentPage++
			HandleFilmsCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackFilmsPrevPage:
		if session.FilmsState.CurrentPage > 1 {
			session.FilmsState.CurrentPage--
			HandleFilmsCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackFilmsLastPage:
		if session.FilmsState.CurrentPage != session.FilmsState.LastPage {
			session.FilmsState.CurrentPage = session.FilmsState.LastPage
			HandleFilmsCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackFilmsFirstPage:
		if session.FilmsState.CurrentPage != 1 {
			session.FilmsState.CurrentPage = 1
			HandleFilmsCommand(app, session)
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			app.SendMessage(msg, nil)
		}

	case callback == states.CallbackFilmsNew:
		HandleNewFilmCommand(app, session)

	case callback == states.CallbackFilmsManage:
		HandleManageFilmCommand(app, session)

	case callback == states.CallbackFilmsFind:
		handleFilmsFindByTitle(app, session)

	case callback == states.CallbackFilmsFilters:
		HandleFiltersFilmsCommand(app, session)

	case callback == states.CallbackFilmsSorting:
		HandleSortingFilmsCommand(app, session)

	case strings.HasPrefix(callback, "select_film_"):
		handleFilmSelect(app, session)
	}
}

func HandleFilmsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Upd) {
		session.ClearAllStates()
		HandleFilmsCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessFindFilmsAwaitingTitle:
		parseFilmsFindTitle(app, session)
	}
}

func handleFilmSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Upd)
	indexStr := strings.TrimPrefix(callback, "select_film_")
	index, err := strconv.Atoi(indexStr)
	if err != nil {
		msg := "🚨 " + translator.Translate(session.Lang, "getFilmFailure", nil, nil)
		keyboard := keyboards.NewKeyboard().AddBack(states.CallbackMenuSelectCollections).Build(session.Lang)
		app.SendMessage(msg, keyboard)
		sl.Log.Error("failed to parse film index", slog.Any("error", err), slog.String("callback", callback))
		return
	}

	session.FilmDetailState.Index = index

	HandleFilmsDetailCommand(app, session)
}

func handleFilmsFindByTitle(app models.App, session *models.Session) {
	msg := "❓ " + translator.Translate(session.Lang, "filmRequestTitle", nil, nil)

	keyboard := keyboards.NewKeyboard().AddCancel().Build(session.Lang)

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessFindFilmsAwaitingTitle)
}

func parseFilmsFindTitle(app models.App, session *models.Session) {
	title := utils.ParseMessageString(app.Upd)

	session.FilmsState.Title = title
	session.FilmsState.CurrentPage = 1

	session.ClearState()

	HandleFindFilmsCommand(app, session)
}

func UpdateFilmsList(app models.App, session *models.Session, next bool) error {
	currentPage := session.FilmsState.CurrentPage
	lastPage := session.FilmsState.LastPage

	switch next {
	case true:
		if currentPage < lastPage {
			session.FilmsState.CurrentPage++
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "lastPageAlert", nil, nil)
			return fmt.Errorf(msg)
		}
	case false:
		if currentPage > 1 {
			session.FilmsState.CurrentPage--
		} else {
			msg := "❗️" + translator.Translate(session.Lang, "firstPageAlert", nil, nil)
			return fmt.Errorf(msg)
		}
	}

	_, err := GetFilms(app, session)
	if err != nil {
		return err
	}

	return nil
}

func GetFilms(app models.App, session *models.Session) (*filters.Metadata, error) {
	var films []apiModels.Film
	var metadata *filters.Metadata
	var err error

	switch session.Context {
	case states.ContextFilm:
		films, metadata, err = fetchFilmsFromUser(app, session)
		if err != nil {
			return nil, err
		}

	case states.ContextCollection:
		films, metadata, err = fetchFilmsFromCollection(app, session)
		if err != nil {
			return nil, err
		}

	default:
		return nil, fmt.Errorf("unsupported session context: %v", session.Context)
	}

	UpdateSessionWithFilms(session, films, metadata)

	return metadata, nil
}

func fetchFilmsFromUser(app models.App, session *models.Session) ([]apiModels.Film, *filters.Metadata, error) {
	filmsResponse, err := watchlist.GetFilms(app, session)
	if err != nil {
		return nil, nil, err
	}
	return filmsResponse.Films, &filmsResponse.Metadata, nil
}

func fetchFilmsFromCollection(app models.App, session *models.Session) ([]apiModels.Film, *filters.Metadata, error) {
	collectionResponse, err := watchlist.GetCollectionFilms(app, session)
	if err != nil {
		return nil, nil, err
	}

	session.CollectionDetailState.Collection = collectionResponse.CollectionFilms.Collection

	return collectionResponse.CollectionFilms.Films, &collectionResponse.Metadata, nil
}

func UpdateSessionWithFilms(session *models.Session, films []apiModels.Film, metadata *filters.Metadata) {
	session.FilmsState.Films = films
	session.FilmsState.LastPage = metadata.LastPage
	session.FilmsState.TotalRecords = metadata.TotalRecords
	session.FilmsState.PageSize = metadata.PageSize
}
