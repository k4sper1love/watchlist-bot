package postgres

import "github.com/k4sper1love/watchlist-bot/internal/models"

func GetSessionByTelegramID(telegramID int) (*models.Session, error) {
	var session models.Session

	if err := GetDB().
		Preload("ProfileState").
		Preload("CollectionsState").
		Preload("FilmsState").
		Preload("FilmDetailState").
		Preload("CollectionDetailState").
		Preload("CollectionFilmState").
		FirstOrInit(&session, models.Session{TelegramID: telegramID}).Error; err != nil {
		return nil, err
	}

	if session.ProfileState == nil {
		session.ProfileState = &models.ProfileState{}
	}

	if session.CollectionsState == nil {
		session.CollectionsState = &models.CollectionsState{}
	}
	if session.FilmsState == nil {
		session.FilmsState = &models.FilmsState{}
	}

	if session.FilmDetailState == nil {
		session.FilmDetailState = &models.FilmDetailState{}
		session.FilmDetailState.Index = -1
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
	Save(session, session.ProfileState, session.CollectionsState, session.FilmsState, session.FilmDetailState, session.CollectionDetailState, session.CollectionFilmState)
}
