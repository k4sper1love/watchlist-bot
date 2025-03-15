package postgres

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

var db *gorm.DB

func ConnectDatabase(config *models.Config) error {
	var err error
	db, err = gorm.Open(postgres.Open(config.DatabaseURL), &gorm.Config{})
	if err != nil {
		sl.Log.Error("failed to open database connection", slog.Any("error", err))
		return err
	}

	return autoMigrate()
}

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

func GetDatabase() *gorm.DB {
	return db
}

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
