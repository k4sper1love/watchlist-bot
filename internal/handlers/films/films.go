package films

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strconv"
	"strings"
)

func HandleFilmsCommand(app models.App, session *models.Session) {
	session.FilmsState.Clear()

	if metadata, err := GetFilms(app, session); err != nil {
		app.SendMessage(messages.BuildFilmsFailureMessage(session), keyboards.BuildKeyboardWithBack(session, ""))
	} else {
		app.SendMessage(messages.BuildFilmsMessage(session, metadata), keyboards.BuildFilmsKeyboard(session, metadata.CurrentPage, metadata.LastPage))
	}
}

func HandleFilmsButtons(app models.App, session *models.Session, back func(models.App, *models.Session)) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallbackFilmsBack:
		if session.Context == states.ContextCollection {
			session.CollectionsState.CurrentPage = 1
		}
		back(app, session)

	case states.CallbackFilmsNextPage, states.CallbackFilmsPrevPage,
		states.CallbackFilmsLastPage, states.CallbackFilmsFirstPage:
		handleFilmsPagination(app, session, callback)

	case states.CallbackFilmsNew:
		HandleNewFilmCommand(app, session)

	case states.CallbackFilmsManage:
		HandleManageFilmCommand(app, session)

	case states.CallbackFilmsFind:
		handleFilmsFindByTitle(app, session)

	case states.CallbackFilmsFilters:
		HandleFiltersFilmsCommand(app, session)

	case states.CallbackFilmsSorting:
		HandleSortingFilmsCommand(app, session)

	default:
		if strings.HasPrefix(callback, states.PrefixSelectFilm) {
			handleFilmSelect(app, session)
		}
	}
}

func handleFilmsPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallbackFilmsNextPage:
		if session.FilmsState.CurrentPage >= session.FilmsState.LastPage {
			app.SendMessage(messages.BuildLastPageAlertMessage(session), nil)
			return
		}
		session.FilmsState.CurrentPage++

	case states.CallbackFilmsPrevPage:
		if session.FilmsState.CurrentPage <= 1 {
			app.SendMessage(messages.BuildFirstPageAlertMessage(session), nil)
			return
		}
		session.FilmsState.CurrentPage--

	case states.CallbackFilmsLastPage:
		if session.FilmsState.CurrentPage == session.FilmsState.LastPage {
			app.SendMessage(messages.BuildLastPageAlertMessage(session), nil)
			return
		}
		session.FilmsState.CurrentPage = session.FilmsState.LastPage

	case states.CallbackFilmsFirstPage:
		if session.FilmsState.CurrentPage == 1 {
			app.SendMessage(messages.BuildFirstPageAlertMessage(session), nil)
			return
		}
		session.FilmsState.CurrentPage = 1
	}

	HandleFilmsCommand(app, session)
}

func HandleFilmsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleFilmsCommand(app, session)
		return
	}

	switch session.State {
	case states.ProcessFindFilmsAwaitingTitle:
		parseFilmFindTitle(app, session, HandleFindFilmsCommand)
	}
}

func handleFilmSelect(app models.App, session *models.Session) {
	callback := utils.ParseCallback(app.Update)
	indexStr := strings.TrimPrefix(callback, states.PrefixSelectFilm)

	if index, err := strconv.Atoi(indexStr); err != nil {
		utils.LogParseSelectError(err, callback)
		app.SendMessage(messages.BuildFilmsFailureMessage(session), keyboards.BuildKeyboardWithBack(session, states.CallbackMenuSelectCollections))
	} else {
		session.FilmDetailState.Index = index
		HandleFilmsDetailCommand(app, session)
	}
}

func handleFilmsFindByTitle(app models.App, session *models.Session) {
	app.SendMessage(messages.BuildFilmRequestTitleMessage(session), keyboards.BuildKeyboardWithCancel(session))
	session.SetState(states.ProcessFindFilmsAwaitingTitle)
}

func UpdateFilmsList(app models.App, session *models.Session, next bool) error {
	if next {
		if session.FilmsState.CurrentPage >= session.FilmsState.LastPage {
			return fmt.Errorf(messages.BuildLastPageAlertMessage(session))
		}
		session.FilmsState.CurrentPage++
	} else {
		if session.FilmsState.CurrentPage <= 1 {
			return fmt.Errorf(messages.BuildFirstPageAlertMessage(session))
		}
		session.FilmsState.CurrentPage--
	}

	_, err := GetFilms(app, session)
	return err
}

func GetFilms(app models.App, session *models.Session) (*filters.Metadata, error) {
	var (
		films    []apiModels.Film
		metadata *filters.Metadata
		err      error
	)

	switch session.Context {
	case states.ContextFilm:
		films, metadata, err = fetchFilmsFromUser(app, session)
	case states.ContextCollection:
		films, metadata, err = fetchFilmsFromCollection(app, session)
	default:
		return nil, fmt.Errorf("unsupported session context: %v", session.Context)
	}

	if err != nil {
		return nil, err
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
