package handlers

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/config"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
)

func handleProfileCommand(app config.App, session *models.Session) {
	user, err := watchlist.GetUser(app, session)
	if err != nil {
		sendMessage(app, err.Error())
		return
	}

	msg := fmt.Sprintf("Вот ваш профиль, %s:\n", user.Username) +
		fmt.Sprintf("Ваш ID в системе API: %d\n", user.ID) +
		fmt.Sprintf("Ваш email: %s\n", user.Email) +
		fmt.Sprintf("Аккаунт был создан %v", user.CreatedAt)

	sendMessage(app, msg)
}
