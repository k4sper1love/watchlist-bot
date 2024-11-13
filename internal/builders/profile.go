package builders

import (
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/models"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
)

func BuildProfileMessage(user *models.User) string {
	msg := fmt.Sprintf("Вот ваш профиль, %s:\n", user.Username) +
		fmt.Sprintf("Ваш ID в системе API: %d\n", user.ID) +
		fmt.Sprintf("Ваш email: %s\n", user.Email) +
		fmt.Sprintf("Аккаунт был создан %v", user.CreatedAt)

	return msg
}

func (k *Keyboard) AddProfileUpdate() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Редактировать", states.CallbackProfileSelectUpdate})

	return k
}

func (k *Keyboard) AddProfileDelete() *Keyboard {
	k.Buttons = append(k.Buttons, Button{"Удалить (!)", states.CallbackProfileSelectDelete})

	return k
}
