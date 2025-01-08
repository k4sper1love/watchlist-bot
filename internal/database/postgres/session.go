package postgres

import (
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"gorm.io/gorm"
)

func GetSessionByTelegramID(app models.App) (*models.Session, error) {
	var session models.Session

	telegramID := utils.ParseTelegramID(app.Upd)
	lang := utils.ParseLanguageCode(app.Upd)
	username := utils.ParseTelegramUsername(app.Upd)

	if err := GetDB().
		Preload("ProfileState").
		Preload("FeedbackState").
		Preload("CollectionsState").
		Preload("CollectionDetailState").
		Preload("FilmsState").
		Preload("FilmsState.FilmFilters").
		Preload("FilmsState.CollectionFilters").
		Preload("FilmsState.FilmSorting").
		Preload("FilmsState.CollectionSorting").
		Preload("FilmDetailState").
		Preload("CollectionFilmsState").
		Preload("AdminState").
		FirstOrInit(&session, models.Session{TelegramID: telegramID}).Error; err != nil {
		return nil, err
	}

	session = initializeSessionDefaults(session, lang, username, app.Vars.RootID)

	return &session, nil
}

func initializeSessionDefaults(session models.Session, lang, username string, rootID int) models.Session {
	if session.TelegramID == rootID {
		session.Role = roles.Root
	}

	if session.TelegramUsername == "" && username != "" {
		session.TelegramUsername = username
	}

	if session.Lang == "" {
		session.Lang = lang
	}

	if session.AdminState == nil {
		session.AdminState = &models.AdminState{}
	}

	if session.ProfileState == nil {
		session.ProfileState = &models.ProfileState{}
	}

	if session.FeedbackState == nil {
		session.FeedbackState = &models.FeedbackState{}
	}

	if session.CollectionsState == nil {
		session.CollectionsState = &models.CollectionsState{}
	}

	if session.CollectionDetailState == nil {
		session.CollectionDetailState = &models.CollectionDetailState{ObjectID: -1}
	}

	if session.FilmsState == nil {
		session.FilmsState = &models.FilmsState{}
	}

	if session.FilmsState.FilmFilters == nil {
		session.FilmsState.FilmFilters = &models.FilmsFilters{}
	}

	if session.FilmsState.CollectionFilters == nil {
		session.FilmsState.CollectionFilters = &models.FilmsFilters{}
	}

	if session.FilmsState.FilmSorting == nil {
		session.FilmsState.FilmSorting = &models.FilmsSorting{}
	}

	if session.FilmsState.CollectionSorting == nil {
		session.FilmsState.CollectionSorting = &models.FilmsSorting{}
	}

	if session.FilmDetailState == nil {
		session.FilmDetailState = &models.FilmDetailState{Index: -1}
	}

	if session.CollectionFilmsState == nil {
		session.CollectionFilmsState = &models.CollectionFilmsState{}
	}

	return session
}

func SaveSessionWithDependencies(session *models.Session) {
	//Save(
	//	session,
	//	session.ProfileState,
	//	session.FeedbackState,
	//	session.CollectionsState,
	//	session.CollectionDetailState,
	//	session.FilmDetailState,
	//	session.CollectionFilmsState,
	//	session.AdminState,
	//	session.FilmsState,
	//	session.FilmsState.FilmFilters,
	//	session.FilmsState.CollectionFilters,
	//	session.FilmsState.FilmSorting,
	//	session.FilmsState.CollectionSorting,
	//)
	db.Session(&gorm.Session{FullSaveAssociations: true}).Save(session)
}
