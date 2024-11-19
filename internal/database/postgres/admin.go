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

func FetchAllUsers() ([]models.Session, error) {
	var sessions []models.Session
	err := GetDB().Order("created_at DESC").Find(&sessions).Error
	return sessions, err
}

func BanUser(telegramID int) error {
	return GetDB().Model(&models.Session{}).Where("telegram_id = ?", telegramID).Update("is_banned", true).Error
}

func UnbanUser(telegramID int) error {
	return GetDB().Model(&models.Session{}).Where("telegram_id = ?", telegramID).Update("is_banned", false).Error
}

func IsUserBanned(telegramID int) (bool, error) {
	var session models.Session

	err := GetDB().Model(&models.Session{}).Where("telegram_id = ?", telegramID).First(&session).Error
	if err != nil {
		return false, err
	}

	return session.IsBanned, nil
}
