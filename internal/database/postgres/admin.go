package postgres

import (
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func GetUserCounts() (int64, error) {
	var count int64

	err := GetDB().Model(&models.Session{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetAllTelegramID() ([]int, error) {
	var telegramIDs []int

	err := GetDB().Model(&models.Session{}).Pluck("telegram_id", &telegramIDs).Error
	if err != nil {
		return nil, err
	}

	return telegramIDs, nil
}
