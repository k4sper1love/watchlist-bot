package postgres

import (
	"errors"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
)

func GetUserCounts() (int64, error) {
	var count int64

	err := GetDB().Model(&models.Session{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetAdminCounts() (int64, error) {
	var count int64

	err := GetDB().Model(&models.Session{}).Where("role > 0").Count(&count).Error
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

func GetAllUsers() ([]models.Session, error) {
	var sessions []models.Session
	err := GetDB().Order("created_at DESC").Find(&sessions).Error
	return sessions, err
}

func GetAllUsersWithPagination(page, pageSize int) ([]models.Session, error) {
	var sessions []models.Session

	offset := (page - 1) * pageSize

	err := GetDB().Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&sessions).Error

	return sessions, err
}

func GetAllAdminsWithPagination(page, pageSize int) ([]models.Session, error) {
	var sessions []models.Session

	offset := (page - 1) * pageSize

	err := GetDB().
		Model(&models.Session{}).
		Where("role > 0").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&sessions).Error

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

func SetUserRole(telegramID int, role roles.Role) (bool, error) {
	err := GetDB().
		Model(&models.Session{}).
		Where("telegram_id = ?", telegramID).
		Update("role", role).Error

	if err != nil {
		return false, err
	}

	return true, nil
}

func GetUserByTelegramID(telegramID int) (*models.Session, error) {
	var session models.Session

	err := GetDB().Model(&models.Session{}).Where("telegram_id = ?", telegramID).First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func GetUserByTelegramUsername(username string) (*models.Session, error) {
	var session models.Session

	err := GetDB().Model(&models.Session{}).Where("telegram_username = ?", username).First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func GetUserByAPIUserID(id int) (*models.Session, error) {
	var sessions []models.Session

	err := GetDB().Find(&sessions).Error
	if err != nil {
		return nil, err
	}

	for _, session := range sessions {
		if session.User.ID == id {
			return &session, nil
		}
	}

	return nil, errors.New("session not found")
}

func GetAdminByTelegramID(telegramID int) (*models.Session, error) {
	var session models.Session

	err := GetDB().Model(&models.Session{}).Where("telegram_id = ? AND role > 0", telegramID).First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}

func GetAdminByTelegramUsername(username string) (*models.Session, error) {
	var session models.Session

	err := GetDB().Model(&models.Session{}).Where("telegram_username = ? AND role > 0", username).First(&session).Error
	if err != nil {
		return nil, err
	}

	return &session, nil
}
