package models

import (
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"gorm.io/gorm"
)

type ProfileState struct {
	gorm.Model
	SessionID uint `json:"-"`
	Username  string
	Email     string
}

type CollectionsState struct {
	gorm.Model
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
	Description  string `json:"description,omitempty"`
}

type CollectionFilmState struct {
	gorm.Model   `json:"-"`
	SessionID    uint           `json:"-"`
	Index        int            `json:"-"`
	Object       apiModels.Film `json:"film" gorm:"serializer:json"`
	Title        string         `json:"title"`
	Year         int            `json:"year,omitempty"`
	Genre        string         `json:"genre,omitempty"`
	Description  string         `json:"description,omitempty"`
	Rating       float64        `json:"rating,omitempty"`
	ImageURL     string         `json:"image_url,omitempty"`
	Comment      string         `json:"comment,omitempty"`
	IsViewed     bool           `json:"is_viewed"`
	IsEditViewed bool           `json:"-"`
	UserRating   float64        `json:"user_rating,omitempty"`
	Review       string         `json:"review,omitempty"`
}

type FilmDetailState struct {
	gorm.Model   `json:"-"`
	SessionID    uint           `json:"-"`
	Index        int            `json:"-"`
	Object       apiModels.Film `json:"film" gorm:"serializer:json"`
	Title        string         `json:"title"`
	Year         int            `json:"year,omitempty"`
	Genre        string         `json:"genre,omitempty"`
	Description  string         `json:"description,omitempty"`
	Rating       float64        `json:"rating,omitempty"`
	ImageURL     string         `json:"image_url,omitempty"`
	Comment      string         `json:"comment,omitempty"`
	IsViewed     bool           `json:"is_viewed"`
	IsEditViewed bool           `json:"-"`
	UserRating   float64        `json:"user_rating,omitempty"`
	Review       string         `json:"review,omitempty"`
}

type FilmsState struct {
	gorm.Model
	SessionID    uint             `json:"-"`
	Object       []apiModels.Film `json:"films" gorm:"serializer:json"`
	LastPage     int              `json:"-"`
	PageSize     int              `json:"-" gorm:"default:5"`
	CurrentPage  int              `json:"-"`
	TotalRecords int
}

func (s *ProfileState) Clear() {
	s.Username = ""
	s.Email = ""
}

func (s *CollectionDetailState) Clear() {
	s.Name = ""
	s.Description = ""
}

func (s *CollectionFilmState) Clear() {
	s.Title = ""
	s.Year = 0
	s.Genre = ""
	s.Description = ""
	s.Rating = 0
	s.ImageURL = ""
	s.Comment = ""
	s.IsViewed = false
	s.IsEditViewed = false
	s.UserRating = 0
	s.Review = ""
}

func (s *FilmDetailState) Clear() {
	s.Title = ""
	s.Year = 0
	s.Genre = ""
	s.Description = ""
	s.Rating = 0
	s.ImageURL = ""
	s.Comment = ""
	s.IsViewed = false
	s.IsEditViewed = false
	s.UserRating = 0
	s.Review = ""
}
