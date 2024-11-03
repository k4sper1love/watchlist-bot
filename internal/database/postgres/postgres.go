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
		&models.Session{},
		&models.CollectionsState{},
		&models.FilmState{},
		&models.CollectionDetailState{},
		&models.CollectionFilmState{},
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
