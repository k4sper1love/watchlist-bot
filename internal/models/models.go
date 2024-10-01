package models

import "gorm.io/gorm"

type User struct {
	gorm.Model
	TelegramID   int `gorm:"unique"`
	AccessToken  string
	RefreshToken string
}
