package models

import "gorm.io/gorm"

type Session struct {
	gorm.Model
	TelegramID            int `gorm:"unique"`
	UserID                int
	AccessToken           string
	RefreshToken          string
	State                 string
	CollectionsState      *CollectionsState      `gorm:"foreignKey:SessionID"`
	CollectionDetailState *CollectionDetailState `gorm:"foreignKey:SessionID"`
	CollectionFilmState   *CollectionFilmState   `gorm:"foreignKey:SessionID"`
	FilmState             *FilmState             `gorm:"foreignKey:SessionID"`
}

func (s *Session) SetState(state string) {
	s.State = state
}

func (s *Session) ResetState() {
	s.State = ""
}
