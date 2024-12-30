package postgres

import (
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func OpenDB(vars *models.Vars) error {
	var err error

	db, err = gorm.Open(postgres.Open(vars.DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	return db.Debug().AutoMigrate(
		&models.Feedback{},
		&models.Session{},
		&models.ProfileState{},
		&models.FeedbackState{},
		&models.CollectionsState{},
		&models.CollectionDetailState{},
		&models.FilmsState{},
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
		db.Save(value)
	}
}
