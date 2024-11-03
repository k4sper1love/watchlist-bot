package postgres

import "github.com/k4sper1love/watchlist-bot/internal/models"

func GetSessionByTelegramID(telegramID int) (*models.Session, error) {
	var session models.Session

	if err := GetDB().
		Preload("CollectionsState").
		Preload("FilmState").
		Preload("CollectionDetailState").
		Preload("CollectionFilmState").
		FirstOrInit(&session, models.Session{TelegramID: telegramID}).Error; err != nil {
		return nil, err
	}

	if session.CollectionsState == nil {
		session.CollectionsState = &models.CollectionsState{}
	}
	if session.FilmState == nil {
		session.FilmState = &models.FilmState{}
		session.FilmState.ObjectID = -1
	}
	if session.CollectionDetailState == nil {
		session.CollectionDetailState = &models.CollectionDetailState{}
		session.CollectionDetailState.ObjectID = -1
	}
	if session.CollectionFilmState == nil {
		session.CollectionFilmState = &models.CollectionFilmState{}
		session.CollectionFilmState.Index = -1
	}

	return &session, nil
}

func SaveSessionWihDependencies(session *models.Session) {
	Save(session, session.CollectionsState, session.FilmState, session.CollectionDetailState, session.CollectionFilmState)
}
