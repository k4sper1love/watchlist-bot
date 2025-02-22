package models

import (
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"gorm.io/gorm"
)

type BaseState struct {
	gorm.Model `json:"-"`
	SessionID  uint `json:"-"`
}

type ProfileState struct {
	BaseState
	Username string
	Email    string
}

type FeedbackState struct {
	BaseState
	Category string
	Message  string
}

type FilmsState struct {
	BaseState
	Films             []apiModels.Film `json:"films" gorm:"serializer:json"`
	LastPage          int              `json:"-"`
	PageSize          int              `json:"-" gorm:"default:4"`
	CurrentPage       int              `json:"-"`
	TotalRecords      int
	Title             string       `json:"-"`
	FilmFilters       *FiltersFilm `gorm:"polymorphic:Filterable;polymorphicValue:FilmFilters"`
	CollectionFilters *FiltersFilm `gorm:"polymorphic:Filterable;polymorphicValue:CollectionFilters"`
	FilmSorting       *Sorting     `gorm:"polymorphic:Sortable;polymorphicValue:FilmSorting"`
	CollectionSorting *Sorting     `gorm:"polymorphic:Sortable;polymorphicValue:CollectionSorting"`
}

type FilmDetailState struct {
	BaseState
	Index        int            `json:"-"`
	Film         apiModels.Film `json:"film" gorm:"serializer:json"`
	IsFavorite   bool           `json:"is_favorite"`
	Title        string         `json:"title,omitempty"`
	Year         int            `json:"year,omitempty"`
	Genre        string         `json:"genre,omitempty"`
	Description  string         `json:"description,omitempty"`
	Rating       float64        `json:"rating,omitempty"`
	ImageURL     string         `json:"image_url,omitempty"`
	Comment      string         `json:"comment,omitempty"`
	IsViewed     bool           `json:"is_viewed"`
	IsEditViewed bool           `json:"-"`
	UserRating   float64        `json:"user_rating"`
	Review       string         `json:"review"`
	URL          string         `json:"url,omitempty"`
}

type CollectionsState struct {
	BaseState
	Collections []apiModels.Collection `json:"collections" gorm:"serializer:json"`
	LastPage    int                    `json:"-"`
	PageSize    int                    `json:"-" gorm:"default:4"`
	CurrentPage int                    `json:"-"`
	Name        string                 `json:"-"`
	Sorting     *Sorting               `gorm:"polymorphic:Sortable;polymorphicValue:Sorting"`
}

type CollectionDetailState struct {
	BaseState
	ObjectID    int                  `json:"-"`
	Collection  apiModels.Collection `json:"collection" gorm:"serializer:json"`
	IsFavorite  bool                 `json:"is_favorite"`
	Name        string               `json:"name,omitempty"`
	Description string               `json:"description,omitempty"`
}

type CollectionFilmsState struct {
	BaseState
	LastPage    int `json:"-"`
	PageSize    int `json:"-" gorm:"default:4"`
	CurrentPage int `json:"-"`
}

type AdminState struct {
	BaseState
	UserID           int        `json:"-"`
	UserLang         string     `json:"-"`
	UserRole         roles.Role `json:"-"`
	FeedbackID       int        `json:"-"`
	LastPage         int        `json:"-"`
	PageSize         int        `json:"-" gorm:"default:4"`
	CurrentPage      int        `json:"-"`
	TotalRecords     int        `json:"-"`
	FeedbackMessage  string     `json:"-"`
	FeedbackImageURL string     `json:"-"`
	NeedFeedbackPin  bool       `json:"-"`
}

func (s *FilmsState) Clear() {
	s.Title = ""
	s.FilmSorting.Clear()
	s.CollectionSorting.Clear()
}

func (s *CollectionsState) Clear() {
	s.Name = ""
	s.Sorting.Clear()
}

func (s *AdminState) Clear() {
	s.FeedbackMessage, s.FeedbackImageURL = "", ""
	s.NeedFeedbackPin = false
}

func (s *ProfileState) Clear() {
	s.Username, s.Email = "", ""
}

func (s *FeedbackState) Clear() {
	s.Category, s.Message = "", ""
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
	s.IsFavorite = false
	s.IsEditViewed = false
	s.UserRating = 0
	s.Review = ""
	s.URL = ""
}

func (s *FilmDetailState) HasIndex() bool {
	return s.Index != -1
}

func (s *FilmDetailState) ClearIndex() {
	s.Index = -1
}

func (s *CollectionDetailState) Clear() {
	s.IsFavorite = false
	s.Name, s.Description = "", ""
}

func (s *FilmDetailState) SetImageURL(url string) {
	s.ImageURL = url
}

func (s *FilmDetailState) SetFromFilm(film *apiModels.Film) {
	s.Title = film.Title
	s.Description = film.Description
	s.Genre = film.Genre
	s.Year = film.Year
	s.Rating = film.Rating
	s.URL = film.URL
}
