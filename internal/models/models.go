package models

import "gorm.io/gorm"

type Feedback struct {
	gorm.Model
	TelegramID int    `gorm:"not null"`
	Category   string `gorm:"not null"`
	Message    string `gorm:"not null"`
}
