package postgres

import (
	"errors"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
)

const (
	TelegramIDField       = "telegram_id"
	TelegramUsernameField = "telegram_username"
)

func GetUserCount(isAdmin bool) (int64, error) {
	var count int64

	query := GetDatabase().Model(&models.Session{})
	if isAdmin {
		query = query.Where("role > 0")
	}

	err := query.Count(&count).Error
	return count, err
}

func GetTelegramIDs() ([]int, error) {
	var telegramIDs []int
	err := GetDatabase().Model(&models.Session{}).Pluck("telegram_id", &telegramIDs).Error
	return telegramIDs, err
}

func GetUsersWithPagination(page, pageSize int, isAdmin bool) ([]models.Session, error) {
	var sessions []models.Session
	offset := utils.CalculateOffset(page, pageSize)

	query := GetDatabase().Order("created_at DESC").Limit(pageSize).Offset(offset)
	if isAdmin {
		query = query.Where("role > 0")
	}

	err := query.Find(&sessions).Error
	return sessions, err
}

func SetUserBanStatus(telegramID int, isBanned bool) error {
	return GetDatabase().Model(&models.Session{}).Where("telegram_id = ?", telegramID).Update("is_banned", isBanned).Error
}

func IsUserBanned(telegramID int) bool {
	var session models.Session
	if err := GetDatabase().Model(&models.Session{}).Where("telegram_id = ?", telegramID).First(&session).Error; err != nil {
		return false
	}
	return session.IsBanned
}

func SetUserRole(telegramID int, role roles.Role) error {
	return GetDatabase().Model(&models.Session{}).Where("telegram_id = ?", telegramID).Update("role", role).Error
}

func GetUserByField(field string, value any, isAdmin bool) (*models.Session, error) {
	var session models.Session

	query := GetDatabase().Model(&models.Session{}).Where(field+" = ?", value)
	if isAdmin {
		query = query.Where("role > 0")
	}

	err := query.First(&session).Error
	return &session, err
}

func GetUserByAPIUserID(id int) (*models.Session, error) {
	var sessions []models.Session
	if err := GetDatabase().Find(&sessions).Error; err != nil {
		return nil, err
	}
	for _, session := range sessions {
		if session.User.ID == id {
			return &session, nil
		}
	}

	return nil, errors.New("session not found")
}
