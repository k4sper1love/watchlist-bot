package postgres

import (
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"log/slog"
)

func SaveFeedback(telegramID int, telegramUsername string, category, feedback string) error {
	feedbackEntry := models.Feedback{
		TelegramID:       telegramID,
		TelegramUsername: telegramUsername,
		Category:         category,
		Message:          feedback,
	}

	if err := GetDB().Create(&feedbackEntry).Error; err != nil {
		sl.Log.Warn(
			"failed to save feedback",
			slog.Any("error", err),
			slog.Int("telegram_id", telegramID),
		)
		return err
	}

	return nil
}

func GetFeedbackCounts() (int64, error) {
	var count int64

	if err := GetDB().Model(&models.Feedback{}).Count(&count).Error; err != nil {
		sl.Log.Warn("failed to get feedback counts", slog.Any("error", err))
		return 0, err
	}

	return count, nil
}

func GetAllFeedbacks() ([]models.Feedback, error) {
	var feedbacks []models.Feedback

	if err := GetDB().Order("created_at DESC").Find(&feedbacks).Error; err != nil {
		sl.Log.Warn("failed to get all feedbacks", slog.Any("error", err))
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
		sl.Log.Warn("failed to get all feedbacks with pagination", slog.Any("error", err))
		return nil, err
	}

	return feedbacks, nil
}

func GetFeedbackByID(id int) (*models.Feedback, error) {
	var feedback models.Feedback

	if err := GetDB().Model(&models.Feedback{}).Where("id = ?", id).First(&feedback).Error; err != nil {
		sl.Log.Warn("failed to get feedback by ID", slog.Any("error", err), slog.Int("feedback_id", id))
		return nil, err
	}

	return &feedback, nil
}
func DeleteFeedbackByID(id int) error {
	if err := GetDB().Delete(&models.Feedback{}, id).Error; err != nil {
		sl.Log.Warn("failed to delete feedback by ID", slog.Any("error", err), slog.Int("feedback_id", id))
		return err
	}

	return nil
}
