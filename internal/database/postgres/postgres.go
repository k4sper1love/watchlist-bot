package postgres

import (
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var db *gorm.DB

func OpenDB(vars *config.Vars) error {
	var err error

	db, err = gorm.Open(postgres.Open(vars.DSN), &gorm.Config{})
	if err != nil {
		return err
	}

	err = db.Debug().AutoMigrate(&models.Session{})
	if err != nil {
		return err
	}

	return nil
}

func GetDB() *gorm.DB {
	return db
}

func Save(value interface{}) {
	db.Save(value)
}
