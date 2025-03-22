package postgres

import (
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"gorm.io/gorm"
	"log/slog"
)

// GetSessionByTelegramID retrieves a session from the database by Telegram ID.
// It preloads all related states and initializes default values if necessary.
func GetSessionByTelegramID(app models.App) (*models.Session, error) {
	var session models.Session
	telegramID := utils.ParseTelegramID(app.Update)

	// Query the database with preloading of all related states
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
		slog.Warn(
			"failed to get session by Telegram ID",
			slog.Any("error", err),
			slog.Int("telegram_id", telegramID),
		)
		return nil, err
	}

	// Initialize default values for the session
	initializeSessionDefaults(app, &session)
	return &session, nil
}

// initializeSessionDefaults sets default values for a session if they are not already set.
// This ensures that all necessary fields and related states are properly initialized.
func initializeSessionDefaults(app models.App, session *models.Session) {
	// Set role to Root if Telegram ID matches the root ID
	if session.TelegramID == app.Config.RootID {
		session.Role = roles.Root
	}

	// Initialize Telegram username if not set
	if session.TelegramUsername == "" {
		session.TelegramUsername = utils.ParseTelegramUsername(app.Update)
	}

	// Initialize language code if not set
	if session.Lang == "" {
		session.Lang = utils.ParseLanguageCode(app.Update)
	}

	// Ensure all state objects are initialized
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
		session.FilmsState.FilmFilters = &models.FilmFilters{
			FilterableID:   session.FilmsState.ID,
			FilterableType: "FilmsState",
		}
	}

	if session.FilmsState.CollectionFilters == nil {
		session.FilmsState.CollectionFilters = &models.FilmFilters{
			FilterableID:   session.FilmsState.ID,
			FilterableType: "CollectionsState",
		}
	}

	if session.FilmsState.FilmSorting == nil {
		session.FilmsState.FilmSorting = &models.Sorting{
			SortableID:   session.FilmsState.ID,
			SortableType: "FilmsState",
		}
	}

	if session.FilmsState.CollectionSorting == nil {
		session.FilmsState.CollectionSorting = &models.Sorting{
			SortableID:   session.FilmsState.ID,
			SortableType: "CollectionsState",
		}
	}

	if session.CollectionsState.Sorting == nil {
		session.CollectionsState.Sorting = &models.Sorting{
			SortableID:   session.CollectionsState.ID,
			SortableType: "CollectionsState",
		}
	}

	if session.FilmDetailState == nil {
		session.FilmDetailState = &models.FilmDetailState{Index: -1}
	}

	if session.CollectionFilmsState == nil {
		session.CollectionFilmsState = &models.CollectionFilmsState{}
	}
}

// SaveSessionWithDependencies saves a session and all its associated dependencies to the database.
// Uses FullSaveAssociations to ensure all related states are saved.
func SaveSessionWithDependencies(session *models.Session) {
	GetDatabase().Session(&gorm.Session{FullSaveAssociations: true}).Save(session)
}
