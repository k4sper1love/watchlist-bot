package admin

import (
	"fmt"
	"github.com/k4sper1love/watchlist-bot/internal/builders/keyboards"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/general"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
	"github.com/k4sper1love/watchlist-bot/pkg/translator"
)

func HandleMenuCommand(app models.App, session *models.Session) {
	part1 := translator.Translate(session.Lang, "adminPanel", nil, nil)
	part2 := translator.Translate(session.Lang, "choiceAction", nil, nil)
	msg := fmt.Sprintf("🛠️ <b>%s</b>\n\n%s", part1, part2)

	keyboard := keyboards.BuildAdminMenuKeyboard(session)

	app.SendMessage(msg, keyboard)

}

func HandleMenuButton(app models.App, session *models.Session) {
	switch utils.ParseCallback(app.Upd) {
	case states.CallbackAdminSelectAdmins:
		session.AdminState.CurrentPage = 1
		general.RequireRole(app, session, HandleAdminsCommand, roles.SuperAdmin)

	case states.CallbackAdminSelectUsers:
		session.AdminState.CurrentPage = 1
		general.RequireRole(app, session, HandleUsersCommand, roles.Admin)

	case states.CallbackAdminSelectBroadcast:
		general.RequireRole(app, session, HandleBroadcastCommand, roles.Admin)

	case states.CallbackAdminSelectFeedback:
		session.AdminState.CurrentPage = 1
		general.RequireRole(app, session, HandleFeedbacksCommand, roles.Helper)
	}
}

//func HandleMenuProcess(app models.App, session *models.Session) {
//
//}
