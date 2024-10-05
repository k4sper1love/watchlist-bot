package models

import (
	"github.com/k4sper1love/watchlist-api/pkg/filters"
	apiModels "github.com/k4sper1love/watchlist-api/pkg/models"
	"gorm.io/gorm"
)

type Session struct {
	gorm.Model
	TelegramID   int `gorm:"unique"`
	IsLogged     bool
	AccessToken  string
	RefreshToken string
	State        string
	AuthState    AuthState `gorm:"embedded" `
}

type AuthState struct {
	Username string `json:"username,omitempty"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type CollectionsResponse struct {
	Collections []apiModels.Collection `json:"collections"`
	Metadata    filters.Metadata       `json:"metadata"`
}
