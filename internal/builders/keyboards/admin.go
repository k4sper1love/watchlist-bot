package keyboards

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/internal/utils"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
)

var superAdminMenuButtons = []Button{
	{"🛡️", "admins", states.CallbackAdminSelectAdmins, "", true},
}

var adminMenuButtons = []Button{
	{"👥", "users", states.CallbackAdminSelectUsers, "", true},
	{"📢", "broadcast", states.CallbackAdminSelectBroadcast, "", true},
}

var helperMenuButtons = []Button{
	{"💬", "feedback", states.CallbackAdminSelectFeedback, "", true},
}

var rolesButtons = []Button{
	{"👤", "user", states.CallbackAdminUserRoleSelectUser, "", true},
	{"👷🏼", "helper", states.CallbackAdminUserRoleSelectHelper, "", true},
	{"👨🏻‍💼", "admin", states.CallbackAdminUserRoleSelectAdmin, "", true},
}

var superRoleButton = Button{"🦸", "superAdmin", states.CallbackAdminUserRoleSelectSuper, "", true}

func BuildAdminMenuKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	if session.Role.HasAccess(roles.SuperAdmin) {
		keyboard.AddButtons(superAdminMenuButtons...)
	}

	if session.Role.HasAccess(roles.Admin) {
		keyboard.AddButtons(adminMenuButtons...)
	}

	if session.Role.HasAccess(roles.Helper) {
		keyboard.AddButtons(helperMenuButtons...)
	}

	keyboard.AddBack("")

	return keyboard.Build(session.Lang)
}

func BuildAdminUserListKeyboard(session *models.Session, users []models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	if len(users) > 0 {
		keyboard.AddSearch(states.CallbackAdminManageUsersSelectFind)
	}

	keyboard.AddAdminUserSelect(session, users)

	keyboard.AddNavigation(
		session.AdminState.CurrentPage,
		session.AdminState.LastPage,
		states.CallbackAdminUsersListPrevPage,
		states.CallbackAdminUsersListNextPage,
		states.CallbackAdminUsersListFirstPage,
		states.CallbackAdminUsersListLastPage,
	)

	keyboard.AddBack(states.CallbackAdminManageUsersSelectBack)

	return keyboard.Build(session.Lang)
}

func BuildAdminListKeyboard(session *models.Session, admins []models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	if len(admins) > 0 {
		keyboard.AddSearch(states.CallbackAdminListSelectFind)
	}

	keyboard.AddAdminSelect(session, admins)

	keyboard.AddNavigation(
		session.AdminState.CurrentPage,
		session.AdminState.LastPage,
		states.CallbackAdminListPrevPage,
		states.CallbackAdminListNextPage,
		states.CallbackAdminListFirstPage,
		states.CallbackAdminListLastPage,
	)

	keyboard.AddBack(states.CallbackAdminListBack)

	return keyboard.Build(session.Lang)
}

func BuildFeedbackListKeyboard(session *models.Session, feedbacks []models.Feedback) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddFeedbackSelect(session, feedbacks)

	keyboard.AddNavigation(
		session.AdminState.CurrentPage,
		session.AdminState.LastPage,
		states.CallbackAdminFeedbackListPrevPage,
		states.CallbackAdminFeedbackListNextPage,
		states.CallbackAdminFeedbackListFirstPage,
		states.CallbackAdminFeedbackListLastPage,
	)

	keyboard.AddBack(states.CallbackAdminFeedbackListBack)

	return keyboard.Build(session.Lang)
}

func BuildAdminUserDetailKeyboard(session *models.Session, user *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	if session.Role.HasAccess(roles.Root) ||
		session.Role.HasAccess(roles.SuperAdmin) && !user.Role.HasAccess(roles.SuperAdmin) {
		keyboard.AddLogs()
	}

	if user.IsBanned && !user.Role.HasAccess(roles.Root) {
		keyboard.AddUnbanUser()
	} else if !user.Role.HasAccess(roles.Root) {
		keyboard.AddBanUser()
	}

	if session.Role.HasAccess(roles.SuperAdmin) && !user.Role.HasAccess(roles.Root) {
		keyboard.AddUserManageRole()
	}

	keyboard.AddBack(states.CallbackAdminUserDetailBack)

	return keyboard.Build(session.Lang)
}

func BuildAdminDetailKeyboard(session *models.Session, user *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	if !user.Role.HasAccess(session.Role.PrevRole()) && !user.Role.HasAccess(roles.Root) {
		keyboard.AddRaiseRank()
	}

	if user.Role.HasAccess(roles.Helper) && !user.Role.HasAccess(roles.Root) {
		keyboard.AddLowerRank()
	}

	if !user.Role.HasAccess(session.Role) && user.Role.HasAccess(roles.Helper) && !user.Role.HasAccess(roles.Root) {
		keyboard.AddRemoveAdminRole()
	}

	keyboard.AddBack(states.CallbackAdminDetailBack)

	return keyboard.Build(session.Lang)
}

func BuildAdminFeedbackDetailKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddFeedbackDelete()

	keyboard.AddBack(states.CallbackAdminFeedbackDetailBack)

	return keyboard.Build(session.Lang)
}

func BuildAdminUserRoleKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.AddButtons(rolesButtons...)

	if session.Role.HasAccess(roles.Root) {
		keyboard.AddButtons(superRoleButton)
	}

	keyboard.AddBack(states.CallbackAdminUserRoleSelectBack)

	return keyboard.Build(session.Lang)
}

func BuildBroadcastConfirmKeyboard(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

	keyboard.addSendBroadcast()

	keyboard.AddCancel()

	return keyboard.Build(session.Lang)
}

func (k *Keyboard) AddAdminUserSelect(session *models.Session, users []models.Session) *Keyboard {
	buttons := userSelectButtons(session, users, "select_admin_user")

	k.AddButtonsWithRowSize(2, buttons...)

	return k
}

func (k *Keyboard) AddAdminSelect(session *models.Session, admins []models.Session) *Keyboard {
	buttons := userSelectButtons(session, admins, "select_admin")

	k.AddButtonsWithRowSize(2, buttons...)

	return k
}

func (k *Keyboard) AddFeedbackSelect(session *models.Session, feedbacks []models.Feedback) *Keyboard {
	var buttons []Button

	for i, feedback := range feedbacks {
		itemID := utils.GetItemID(i, session.AdminState.CurrentPage, session.AdminState.PageSize)

		text := fmt.Sprintf("%s", utils.NumberToEmoji(itemID))

		buttons = append(buttons, Button{"", text, fmt.Sprintf("select_admin_feedback_%d", feedback.ID), "", false})
	}

	k.AddButtonsWithRowSize(2, buttons...)

	return k
}

func (k *Keyboard) AddLogs() *Keyboard {
	return k.AddButton("💾", "logs", states.CallbackAdminUserDetailLogs, "", true)
}

func (k *Keyboard) AddUserManageRole() *Keyboard {
	return k.AddButton("🔄", "manageUserRole", states.CallbackAdminUserDetailRole, "", true)
}

func (k *Keyboard) AddUnbanUser() *Keyboard {
	return k.AddButton("✅", "unban", states.CallbackAdminUserDetailUnban, "", true)
}

func (k *Keyboard) AddBanUser() *Keyboard {
	return k.AddButton("❌", "ban", states.CallbackAdminUserDetailBan, "", true)
}

func (k *Keyboard) AddViewUserFeedback() *Keyboard {
	return k.AddButton("📩", "viewFeedback", states.CallbackAdminUserDetailFeedback, "", true)
}

func (k *Keyboard) AddRaiseRank() *Keyboard {
	return k.AddButton("⬆️", "raiseRole", states.CallbackAdminDetailRaiseRole, "", true)
}

func (k *Keyboard) AddLowerRank() *Keyboard {
	return k.AddButton("⬇️", "lowerRole", states.CallbackAdminDetailLowerRole, "", true)
}

func (k *Keyboard) AddRemoveAdminRole() *Keyboard {
	return k.AddButton("⚠️", "removeAdminRole", states.CallbackAdminDetailRemoveRole, "", true)
}

func (k *Keyboard) AddFeedbackDelete() *Keyboard {
	return k.AddButton("🗑️", "delete", states.CallbackAdminFeedbackDetailDelete, "", true)
}

func (k *Keyboard) addSendBroadcast() *Keyboard {
	return k.AddButton("➤", "send", states.CallbackAdminBroadcastSend, "", true)
}

func userSelectButtons(session *models.Session, users []models.Session, callback string) []Button {
	var buttons []Button

	for i, user := range users {
		itemID := utils.GetItemID(i, session.AdminState.CurrentPage, session.AdminState.PageSize)

		text := fmt.Sprintf("%s %d", utils.NumberToEmoji(itemID), user.TelegramID)

		if user.TelegramUsername != "" {
			text += fmt.Sprintf(" (@%s)", user.TelegramUsername)
		}

		buttons = append(buttons, Button{"", text, fmt.Sprintf("%s_%d", callback, user.TelegramID), "", false})
	}

	return buttons
}
