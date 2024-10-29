package models

import (
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	TelegramID          int `gorm:"unique"`
	UserID              int `gorm:"default:-1"`
	AccessToken         string
	RefreshToken        string
	State               string
	CollectionState     CollectionState     `gorm:"foreignKey:SessionID"`
	FilmState           FilmState           `gorm:"foreignKey:SessionID"`
	CollectionFilmState CollectionFilmState `gorm:"foreignKey:SessionID"`
}

type CollectionState struct {
	gorm.Model
	SessionID   uint
	ObjectID    int    `json:"-" gorm:"default:-1"`
	Name        string `json:"name"`
	Description string `json:"description"`
	CurrentPage int    `json:"-"`
	LastPage    int    `json:"-"`
	PageSize    int    `json:"-" gorm:"default:1"`
}

type FilmState struct {
	gorm.Model
	SessionID   uint
	ObjectID    int `json:"-" gorm:"default:-1"`
	CurrentPage int
	LastPage    int
	PageSize    int `json:"-" gorm:"default:1"`
}

type CollectionFilmState struct {
	gorm.Model
	SessionID   uint
	ObjectID    int `json:"-" gorm:"default:-1"`
	CurrentPage int
	LastPage    int
	PageSize    int `json:"-" gorm:"default:1"`
}

type CollectionsResponse struct {
	Collections []apiModels.Collection `json:"collections"`
	Metadata    filters.Metadata       `json:"metadata"`
}

type FilmsResponse struct {
	Films    []apiModels.Film `json:"films"`
	Metadata filters.Metadata `json:"metadata"`
}

type CollectionFilmsResponse struct {
	CollectionFilms apiModels.CollectionFilms `json:"collection_films"`
	Metadata        filters.Metadata          `json:"metadata"`
}
