package models

import (
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	TelegramID            int            `gorm:"unique"`
	IsAdmin               bool           `gorm:"default:false"`
	User                  apiModels.User `json:"user" gorm:"serializer:json"`
	AccessToken           string         `json:"access_token"`
	RefreshToken          string         `json:"refresh_token"`
	State                 string
	ProfileState          *ProfileState          `gorm:"foreignKey:SessionID"`
	CollectionsState      *CollectionsState      `gorm:"foreignKey:SessionID"`
	CollectionDetailState *CollectionDetailState `gorm:"foreignKey:SessionID"`
	CollectionFilmState   *CollectionFilmState   `gorm:"foreignKey:SessionID"`
	FilmsState            *FilmsState            `gorm:"foreignKey:SessionID"`
	FilmDetailState       *FilmDetailState       `gorm:"foreignKey:SessionID"`
}

func (s *Session) SetState(state string) {
	s.State = state
}

func (s *Session) ClearState() {
	s.State = ""
}

func (s *Session) ClearUser() {
	s.User = apiModels.User{}
}

func (s *Session) ClearFull() {
	s.AccessToken = ""
	s.RefreshToken = ""
	s.ClearState()
	s.ClearUser()
	s.ProfileState.Clear()
	s.FilmDetailState.Clear()
	s.CollectionDetailState.Clear()
	s.CollectionFilmState.Clear()
}
