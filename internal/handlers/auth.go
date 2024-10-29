package handlers

import (
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
)

func handleLogoutCommand(app config.App, session *models.Session) {
	sendMessage(app, "Вы точно хотите из аккаунта? Отправьте \"+\" для выхода")
	setState(session, ProcessLogoutStateAwaitingConfirm)
}

func handleLogoutProcess(app config.App, session *models.Session) {
	switch session.State {
	case ProcessLogoutStateAwaitingConfirm:
		if parseMessageText(app.Upd) == "+" {
			if err := watchlist.Logout(app, session); err != nil {
				sendMessage(app, "Неудачный выход из системы")
			} else {
				sendMessage(app, "Успешно вышли из системы")
				clearSession(session)
			}
		} else {
			sendMessage(app, "Отмена выхода из системы")
		}
		resetState(session)
	}
}

func handleAuthProcess(app config.App, session *models.Session) error {
	if session.AccessToken != "" {
		if watchlist.IsTokenValid(app, session.AccessToken) {
			return nil
		} else if session.RefreshToken != "" {
			if err := watchlist.RefreshAccessToken(app, session); err == nil {
				return nil
			}
		}
	}

	err := watchlist.Register(app, session)
	if err == nil {
		return nil
	}
	err = watchlist.Login(app, session)
	if err == nil {
		return nil
	}

	return err
}
