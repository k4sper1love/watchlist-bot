package models

import (
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	TelegramID            int            `gorm:"unique"`
	User                  apiModels.User `json:"user" gorm:"serializer:json"`
	Lang                  string
	AccessToken           string `json:"access_token"`
	RefreshToken          string `json:"refresh_token"`
	IsAdmin               bool   `gorm:"default:false"`
	IsBanned              bool   `gorm:"default:false"`
	State                 string
	Context               string
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
	s.FilmDetailState.Clear()
	s.CollectionDetailState.Clear()
}

func (s *Session) Logout() {
	s.AccessToken = ""
	s.RefreshToken = ""
	s.ClearUser()
	s.ClearContext()
	s.ClearAllStates()
}
