package postgres

import (
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func GetSessionByTelegramID(app models.App) (*models.Session, error) {
	var session models.Session

	telegramID := utils.ParseTelegramID(app.Upd)
	lang := utils.ParseLanguageCode(app.Upd)

	if err := GetDB().
		Preload("ProfileState").
		Preload("FeedbackState").
		Preload("CollectionsState").
		Preload("CollectionDetailState").
		Preload("FilmsState").
		Preload("FilmDetailState").
		Preload("CollectionFilmsState").
		FirstOrInit(&session, models.Session{TelegramID: telegramID}).Error; err != nil {
		return nil, err
	}

	if session.TelegramID == app.Vars.AdminID {
		session.Role = config.SuperAdminRole
	}

	if session.Lang == "" {
		session.Lang = lang
	}

	if session.ProfileState == nil {
		session.ProfileState = &models.ProfileState{}
	}

	if session.FeedbackState == nil {
		session.FeedbackState = &models.FeedbackState{}
	}

	if session.CollectionsState == nil {
		session.CollectionsState = &models.CollectionsState{}
	}

	if session.CollectionDetailState == nil {
		session.CollectionDetailState = &models.CollectionDetailState{}
		session.CollectionDetailState.ObjectID = -1
	}

	if session.FilmsState == nil {
		session.FilmsState = &models.FilmsState{}
	}

	if session.FilmDetailState == nil {
		session.FilmDetailState = &models.FilmDetailState{}
		session.FilmDetailState.Index = -1
	}

	if session.CollectionFilmsState == nil {
		session.CollectionFilmsState = &models.CollectionFilmsState{}
	}

	return &session, nil
}

func SaveSessionWihDependencies(session *models.Session) {
	Save(session, session.ProfileState, session.FeedbackState, session.CollectionsState, session.CollectionDetailState, session.FilmsState, session.FilmDetailState, session.CollectionFilmsState)
}
