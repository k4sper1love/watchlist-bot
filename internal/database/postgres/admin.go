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

	if err := GetDB().Order("created_at DESC").Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

func GetAllUsersWithPagination(page, pageSize int) ([]models.Session, error) {
	var sessions []models.Session
	offset := (page - 1) * pageSize

	if err := GetDB().Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

func GetAllAdminsWithPagination(page, pageSize int) ([]models.Session, error) {
	var sessions []models.Session
	offset := (page - 1) * pageSize

	if err := GetDB().
		Model(&models.Session{}).
		Where("role > 0").
		Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&sessions).Error; err != nil {
		return nil, err
	}

	return sessions, nil
}

func BanUser(telegramID int) error {
	if err := GetDB().
		Model(&models.Session{}).
		Where("telegram_id = ?", telegramID).
		Update("is_banned", true).
		Error; err != nil {
		return err
	}

	return nil
}

func UnbanUser(telegramID int) error {
	if err := GetDB().
		Model(&models.Session{}).
		Where("telegram_id = ?", telegramID).
		Update("is_banned", false).
		Error; err != nil {
		return err
	}

	return nil
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
	if err := GetDB().
		Model(&models.Session{}).
		Where("telegram_id = ?", telegramID).
		Update("role", role).
		Error; err != nil {
		return false, nil
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

	if err := GetDB().
		Model(&models.Session{}).
		Where("telegram_username = ?", username).
		First(&session).
		Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func GetUserByAPIUserID(id int) (*models.Session, error) {
	var sessions []models.Session

	if err := GetDB().Find(&sessions).Error; err != nil {
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

	if err := GetDB().
		Model(&models.Session{}).
		Where("telegram_id = ? AND role > 0", telegramID).
		First(&session).
		Error; err != nil {
		return nil, err
	}

	return &session, nil
}

func GetAdminByTelegramUsername(username string) (*models.Session, error) {
	var session models.Session

	if err := GetDB().
		Model(&models.Session{}).
		Where("telegram_username = ? AND role > 0", username).
		First(&session).
		Error; err != nil {
		return nil, err
	}

	return &session, nil
}
