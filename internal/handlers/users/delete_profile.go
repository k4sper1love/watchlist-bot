package users

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/services/watchlist"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
)

func HandleDeleteProfileCommand(app models.App, session *models.Session) {
	msg := fmt.Sprintf("Вы уверены, что хотите удалить свой аккаунт %q?", session.User.Username)

	keyboard := keyboards.NewKeyboard().
		AddSurvey().
		Build()

	app.SendMessage(msg, keyboard)

	session.SetState(states.ProcessDeleteProfileAwaitingConfirm)
}

func HandleDeleteProfileProcess(app models.App, session *models.Session) {
	switch session.State {
	case states.ProcessDeleteProfileAwaitingConfirm:
		parseDeleteProfileConfirm(app, session)
	}
}

func parseDeleteProfileConfirm(app models.App, session *models.Session) {
	session.ClearState()

	switch utils.IsAgree(app.Upd) {
	case true:
		if err := watchlist.DeleteUser(app, session); err != nil {
			app.SendMessage("Не удалось удалить профиль", nil)
			HandleProfileCommand(app, session)
			return
		}
		app.SendMessage("Профиль успешно удален!", nil)
		session.Logout()
		general.HandleMenuCommand(app, session)

	case false:
		app.SendMessage("Действие отменено", nil)
		HandleProfileCommand(app, session)
	}
}
