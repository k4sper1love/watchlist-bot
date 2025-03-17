package postgres

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

var db *gorm.DB // Global database connection instance

// ConnectDatabase establishes a connection to the PostgreSQL database using the provided configuration.
// It also performs automatic migration to ensure the database schema is up-to-date.
func ConnectDatabase(config *models.Config) error {
	var err error
	db, err = gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		sl.Log.Error("failed to open database connection", slog.Any("error", err))
		return err
	}

	// Perform automatic migration to update the database schema
	return autoMigrate()
}

// autoMigrate performs automatic migration for all defined models.
// Ensures that the database schema matches the structure of the models.
func autoMigrate() error {
	return db.AutoMigrate(
		&models.Feedback{},
		&models.Session{},
		&models.FilmsState{},
		&models.FilmFilters{},
		&models.Sorting{},
		&models.ProfileState{},
		&models.FeedbackState{},
		&models.CollectionsState{},
		&models.CollectionDetailState{},
		&models.FilmDetailState{},
		&models.CollectionFilmsState{},
		&models.AdminState{},
	)
}

// GetDatabase returns the global database connection instance.
// This function is used to access the database throughout the application.
func GetDatabase() *gorm.DB {
	return db
}

// SaveRecords saves multiple records to the database.
func SaveRecords(records ...interface{}) {
	for _, record := range records {
		if err := db.Save(record).Error; err != nil {
			sl.Log.Warn(
				"failed to save record",
				slog.Any("error", err),
				slog.Any("record", record))
		}
	}
}
