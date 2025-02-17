package postgres

import (
	"github.com/k4sper1love/watchlist-bot/internal/models"
)

func SaveFeedback(telegramID int, telegramUsername string, category, feedback string) error {
	feedbackEntry := models.Feedback{
		TelegramID:       telegramID,
		TelegramUsername: telegramUsername,
		Category:         category,
		Message:          feedback,
	}

	if err := GetDB().Create(&feedbackEntry).Error; err != nil {
		return err
	}

	return nil
}

func GetFeedbackCounts() (int64, error) {
	var count int64

	if err := GetDB().Model(&models.Feedback{}).Count(&count).Error; err != nil {
		return 0, err
	}

	return count, nil
}

func GetAllFeedbacks() ([]models.Feedback, error) {
	var feedbacks []models.Feedback

	if err := GetDB().Order("created_at DESC").Find(&feedbacks).Error; err != nil {
		return nil, err
	}

	return feedbacks, nil
}

func GetAllFeedbacksWithPagination(page, pageSize int) ([]models.Feedback, error) {
	var feedbacks []models.Feedback
	offset := (page - 1) * pageSize

	if err := GetDB().Order("created_at DESC").
		Limit(pageSize).
		Offset(offset).
		Find(&feedbacks).
		Error; err != nil {
		return nil, err
	}

	return feedbacks, nil
}

func GetFeedbackByID(id int) (*models.Feedback, error) {
	var feedback models.Feedback

	if err := GetDB().Model(&models.Feedback{}).Where("id = ?", id).First(&feedback).Error; err != nil {
		return nil, err
	}

	return &feedback, nil
}
func DeleteFeedbackByID(id int) error {
	if err := GetDB().Delete(&models.Feedback{}, id).Error; err != nil {
		return err
	}

	return nil
}
