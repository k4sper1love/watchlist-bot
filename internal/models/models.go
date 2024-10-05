package models

import "gorm.io/gorm"

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
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}
