package models

import (
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"gorm.io/gorm"
)

type CollectionsState struct {
	gorm.Model  `json:"-"`
	SessionID   uint                   `json:"-"`
	Object      []apiModels.Collection `json:"collections" gorm:"serializer:json"`
	LastPage    int                    `json:"-"`
	PageSize    int                    `json:"-" gorm:"default:5"`
	CurrentPage int                    `json:"-"`
}

type CollectionDetailState struct {
	gorm.Model
	SessionID    uint
	ObjectID     int                       `json:"-"`
	Object       apiModels.CollectionFilms `json:"collection_films" gorm:"serializer:json"`
	CurrentPage  int
	LastPage     int
	PageSize     int `json:"-" gorm:"default:5"`
	TotalRecords int
	Name         string `json:"name"`
	Description  string `json:"description"`
}

type CollectionFilmState struct {
	gorm.Model  `json:"-"`
	SessionID   uint           `json:"-"`
	Index       int            `json:"-"`
	Object      apiModels.Film `json:"film" gorm:"serializer:json"`
	Title       string         `json:"title"`
	Year        int            `json:"year,omitempty"`
	Genre       string         `json:"genre,omitempty"`
	Description string         `json:"description,omitempty"`
	Rating      float64        `json:"rating,omitempty"`
	ImageURL    string         `json:"image_url,omitempty"`
	Comment     string         `json:"comment,omitempty"`
	IsViewed    bool           `json:"is_viewed,omitempty"`
	UserRating  float64        `json:"user_rating,omitempty"`
	Review      string         `json:"review,omitempty"`
}

type FilmState struct {
	gorm.Model
	SessionID   uint
	ObjectID    int `json:"-" gorm:"default:-1"`
	CurrentPage int
	LastPage    int
	PageSize    int `json:"-" gorm:"default:1"`
}
