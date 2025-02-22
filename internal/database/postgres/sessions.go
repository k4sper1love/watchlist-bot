package postgres

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"gorm.io/gorm"
	"log/slog"
)

func GetSessionByTelegramID(app models.App) (*models.Session, error) {
	var session models.Session
	telegramID := utils.ParseTelegramID(app.Update)

	if err := GetDatabase().
		Preload("ProfileState").
		Preload("FeedbackState").
		Preload("CollectionsState").
		Preload("CollectionsState.Sorting").
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
		sl.Log.Warn(
			"failed to get session by Telegram ID",
			slog.Any("error", err),
			slog.Int("telegram_id", telegramID),
		)
		return nil, err
	}

	initializeSessionDefaults(app, &session)
	return &session, nil
}

func initializeSessionDefaults(app models.App, session *models.Session) {
	if session.TelegramID == app.Config.RootID {
		session.Role = roles.Root
	}

	if session.TelegramUsername == "" {
		session.TelegramUsername = utils.ParseTelegramUsername(app.Update)
	}

	if session.Lang == "" {
		session.Lang = utils.ParseLanguageCode(app.Update)
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

	if session.CollectionsState.Sorting == nil {
		session.CollectionsState.Sorting = &models.Sorting{}
	}

	if session.CollectionDetailState == nil {
		session.CollectionDetailState = &models.CollectionDetailState{ObjectID: -1}
	}

	if session.FilmsState == nil {
		session.FilmsState = &models.FilmsState{}
	}

	if session.FilmsState.FilmFilters == nil {
		session.FilmsState.FilmFilters = &models.FiltersFilm{}
	}

	if session.FilmsState.CollectionFilters == nil {
		session.FilmsState.CollectionFilters = &models.FiltersFilm{}
	}

	if session.FilmsState.FilmSorting == nil {
		session.FilmsState.FilmSorting = &models.Sorting{}
	}

	if session.FilmsState.CollectionSorting == nil {
		session.FilmsState.CollectionSorting = &models.Sorting{}
	}

	if session.FilmDetailState == nil {
		session.FilmDetailState = &models.FilmDetailState{Index: -1}
	}

	if session.CollectionFilmsState == nil {
		session.CollectionFilmsState = &models.CollectionFilmsState{}
	}
}

func SaveSessionWithDependencies(session *models.Session) {
	GetDatabase().Session(&gorm.Session{FullSaveAssociations: true}).Save(session)
}
