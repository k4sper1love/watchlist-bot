package films

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/builders/messages"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/parser"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"strconv"
	"strings"
)

func HandleFilmsCommand(app models.App, session *models.Session) {
	session.FilmsState.Clear()

	if metadata, err := getFilms(app, session); err != nil {
		app.SendMessage(messages.FilmsFailure(session), keyboards.Back(session, ""))
	} else {
		app.SendMessage(messages.Films(session, metadata), keyboards.Films(session, metadata.CurrentPage, metadata.LastPage))
	}
}

func HandleFilmsButtons(app models.App, session *models.Session, back func(models.App, *models.Session)) {
	callback := utils.ParseCallback(app.Update)

	switch callback {
	case states.CallFilmsBack:
		if session.Context == states.CtxCollection {
			session.CollectionsState.CurrentPage = 1
		}
		back(app, session)

	case states.CallFilmsNew:
		HandleNewFilmCommand(app, session)

	case states.CallFilmsManage:
		HandleManageFilmCommand(app, session)

	case states.CallFilmsFind:
		requestFindFilmsTitle(app, session)

	case states.CallFilmsFilters:
		HandleFilmFiltersCommand(app, session)

	case states.CallFilmsSorting:
		HandleSortingFilmsCommand(app, session)

	default:
		if strings.HasPrefix(callback, states.FilmsPage) {
			handleFilmsPagination(app, session, callback)
		}

		if strings.HasPrefix(callback, states.SelectFilm) {
			handleFilmSelect(app, session, callback)
		}
	}
}

func HandleFilmsProcess(app models.App, session *models.Session) {
	if utils.IsCancel(app.Update) {
		session.ClearAllStates()
		HandleFilmsCommand(app, session)
		return
	}

	switch session.State {
	case states.AwaitFilmsTitle:
		parser.ParseFindFilmsTitle(app, session, HandleFindFilmsCommand)
	}
}

func handleFilmsPagination(app models.App, session *models.Session, callback string) {
	switch callback {
	case states.CallFilmsPageNext:
		if session.FilmsState.CurrentPage >= session.FilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage++

	case states.CallFilmsPagePrev:
		if session.FilmsState.CurrentPage <= 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage--

	case states.CallFilmsPageLast:
		if session.FilmsState.CurrentPage == session.FilmsState.LastPage {
			app.SendMessage(messages.LastPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage = session.FilmsState.LastPage

	case states.CallFilmsPageFirst:
		if session.FilmsState.CurrentPage == 1 {
			app.SendMessage(messages.FirstPageAlert(session), nil)
			return
		}
		session.FilmsState.CurrentPage = 1
	}

	HandleFilmsCommand(app, session)
}

func handleFilmSelect(app models.App, session *models.Session, callback string) {
	if index, err := strconv.Atoi(strings.TrimPrefix(callback, states.SelectFilm)); err != nil {
		utils.LogParseSelectError(err, callback)
		app.SendMessage(messages.FilmsFailure(session), keyboards.Back(session, states.CallFilmsBack))
	} else {
		session.FilmDetailState.Index = index
		HandleFilmDetailCommand(app, session)
	}
}

func requestFindFilmsTitle(app models.App, session *models.Session) {
	app.SendMessage(messages.RequestFilmTitle(session), keyboards.Cancel(session))
	session.SetState(states.AwaitFilmsTitle)
}

func updateFilmList(app models.App, session *models.Session, next bool) error {
	if next {
		if session.FilmsState.CurrentPage >= session.FilmsState.LastPage {
			return fmt.Errorf(messages.LastPageAlert(session))
		}
		session.FilmsState.CurrentPage++

	} else {
		if session.FilmsState.CurrentPage <= 1 {
			return fmt.Errorf(messages.FirstPageAlert(session))
		}
		session.FilmsState.CurrentPage--
	}

	_, err := getFilms(app, session)
	return err
}

func getFilms(app models.App, session *models.Session) (*filters.Metadata, error) {
	var (
		films    []apiModels.Film
		metadata *filters.Metadata
		err      error
	)

	switch session.Context {
	case states.CtxFilm:
		films, metadata, err = fetchFilmsFromUser(app, session)
	case states.CtxCollection:
		films, metadata, err = fetchFilmsFromCollection(app, session)
	default:
		return nil, fmt.Errorf("unsupported session context: %v", session.Context)
	}

	if err != nil {
		return nil, err
	}

	updateSessionWithFilms(session, films, metadata)
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

func updateSessionWithFilms(session *models.Session, films []apiModels.Film, metadata *filters.Metadata) {
	session.FilmsState.Films = films
	session.FilmsState.LastPage = metadata.LastPage
	session.FilmsState.TotalRecords = metadata.TotalRecords
	session.FilmsState.PageSize = metadata.PageSize
}
