package models

import (
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	TelegramID            int `gorm:"unique"`
	TelegramUsername      string
	Role                  roles.Role `json:"role" gorm:"serializer:json"`
	Lang                  string
	IsBanned              bool           `gorm:"default:false"`
	User                  apiModels.User `json:"user" gorm:"serializer:json"`
	AccessToken           string         `json:"access_token"`
	RefreshToken          string         `json:"refresh_token"`
	KinopoiskAPIToken     string         `json:"kinopoisk_api_token"`
	State                 string
	Context               string
	AdminState            *AdminState            `gorm:"foreignKey:SessionID"`
	ProfileState          *ProfileState          `gorm:"foreignKey:SessionID"`
	FeedbackState         *FeedbackState         `gorm:"foreignKey:SessionID"`
	CollectionsState      *CollectionsState      `gorm:"foreignKey:SessionID"`
	CollectionDetailState *CollectionDetailState `gorm:"foreignKey:SessionID"`
	FilmsState            *FilmsState            `gorm:"foreignKey:SessionID"`
	FilmDetailState       *FilmDetailState       `gorm:"foreignKey:SessionID"`
	CollectionFilmsState  *CollectionFilmsState  `gorm:"foreignKey:SessionID"`
}

func (s *Session) SetContext(context string) {
	s.Context = context
}

func (s *Session) SetState(state string) {
	s.State = state
}

func (s *Session) ClearState() {
	s.State = ""
}

func (s *Session) ClearContext() {
	s.Context = ""
}

func (s *Session) ClearUser() {
	s.User = apiModels.User{}
}

func (s *Session) ClearAllStates() {
	s.ClearState()
	s.ProfileState.Clear()
	s.FeedbackState.Clear()
	s.FilmsState.Clear()
	s.FilmDetailState.Clear()
	s.CollectionsState.Clear()
	s.CollectionDetailState.Clear()
	s.AdminState.Clear()
}

func (s *Session) Logout() {
	s.AccessToken, s.RefreshToken, s.KinopoiskAPIToken = "", "", ""
	s.ClearUser()
	s.ClearContext()
	s.ClearAllStates()
}

func (s *Session) GetFilmFiltersByCtx() *FilmFilters {
	switch s.Context {
	case states.CtxFilm:
		return s.FilmsState.FilmFilters
	case states.CtxCollection:
		return s.FilmsState.CollectionFilters
	default:
		return &FilmFilters{}
	}
}

func (s *Session) GetFilmSortingByCtx() *Sorting {
	switch s.Context {
	case states.CtxFilm:
		return s.FilmsState.FilmSorting
	case states.CtxCollection:
		return s.FilmsState.CollectionSorting
	default:
		return &Sorting{}
	}
}
