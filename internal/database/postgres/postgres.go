package postgres

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log/slog"
)

var db *gorm.DB

func OpenDB(vars *models.Vars) error {
	var err error

	db, err = gorm.Open(postgres.Open(vars.DSN), &gorm.Config{})
	if err != nil {
		sl.Log.Error("failed to open database connection", slog.Any("error", err))
		return err
	}
	return db.AutoMigrate(
		&models.Feedback{},
		&models.Session{},
		&models.FilmsState{},
		&models.FiltersFilm{},
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

func GetDB() *gorm.DB {
	return db
}

func Save(values ...interface{}) {
	for _, value := range values {
		if err := db.Save(value).Error; err != nil {
			sl.Log.Warn(
				"failed to save data",
				slog.Any("error", err),
				slog.Any("value", value),
			)
		}
	}
}
