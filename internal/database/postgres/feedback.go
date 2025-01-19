package postgres

import (
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func SaveFeedbackToDatabase(telegramID int, telegramUsername string, category, feedback string) error {
	feedbackEntry := models.Feedback{
		TelegramID:       telegramID,
		TelegramUsername: telegramUsername,
		Category:         category,
		Message:          feedback,
	}

	return GetDB().Create(&feedbackEntry).Error
}

func GetFeedbackCounts() (int64, error) {
	var count int64

	err := GetDB().Model(&models.Feedback{}).Count(&count).Error
	if err != nil {
		return 0, err
	}

	return count, nil
}

func GetAllFeedbacks() ([]models.Feedback, error) {
	var feedbacks []models.Feedback

	err := GetDB().Order("created_at DESC").Find(&feedbacks).Error
	return feedbacks, err
}

func GetAllFeedbacksWithPagination(page, pageSize int) ([]models.Feedback, error) {
	var feedbacks []models.Feedback

	offset := (page - 1) * pageSize

	err := GetDB().Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&feedbacks).Error

	return feedbacks, err
}

func GetFeedbackByID(id int) (*models.Feedback, error) {
	var feedback models.Feedback

	err := GetDB().Model(&models.Feedback{}).Where("id = ?", id).First(&feedback).Error
	if err != nil {
		return nil, err
	}

	return &feedback, nil
}
func DeleteFeedbackByID(id int) error {
	return GetDB().Delete(&models.Feedback{}, id).Error
}
