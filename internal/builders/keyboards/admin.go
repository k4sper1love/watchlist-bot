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
	{"🛡️", "admins", states.CallbackAdminSelectAdmins, ""},
}

var adminMenuButtons = []Button{
	{"👥", "users", states.CallbackAdminSelectUsers, ""},
	{"📢", "broadcast", states.CallbackAdminSelectBroadcast, ""},
}

var helperMenuButtons = []Button{
	{"💬", "feedback", states.CallbackAdminSelectFeedback, ""},
}

var rolesButtons = []Button{
	{"👤", "user", states.CallbackAdminUserRoleSelectUser, ""},
	{"👷🏼", "helper", states.CallbackAdminUserRoleSelectHelper, ""},
	{"👨🏻‍💼", "admin", states.CallbackAdminUserRoleSelectAdmin, ""},
}

var superRoleButton = Button{"🦸", "superAdmin", states.CallbackAdminUserRoleSelectSuper, ""}

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
	)

	keyboard.AddBack(states.CallbackAdminFeedbackListBack)

	return keyboard.Build(session.Lang)
}

func BuildAdminUserDetailKeyboard(session *models.Session, user *models.Session) *tgbotapi.InlineKeyboardMarkup {
	keyboard := NewKeyboard()

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

		text := fmt.Sprintf("%d - ID: %d", itemID, feedback.ID)

		buttons = append(buttons, Button{"", text, fmt.Sprintf("select_admin_feedback_%d", feedback.ID), ""})
	}

	k.AddButtonsWithRowSize(2, buttons...)

	return k
}

func (k *Keyboard) AddUserManageRole() *Keyboard {
	return k.AddButton("🎭", "manageUserRole", states.CallbackAdminUserDetailRole, "")
}

func (k *Keyboard) AddUnbanUser() *Keyboard {
	return k.AddButton("🟢", "unban", states.CallbackAdminUserDetailUnban, "")
}

func (k *Keyboard) AddBanUser() *Keyboard {
	return k.AddButton("🔴️", "ban", states.CallbackAdminUserDetailBan, "")
}

func (k *Keyboard) AddViewUserFeedback() *Keyboard {
	return k.AddButton("📩", "viewFeedback", states.CallbackAdminUserDetailFeedback, "")
}

func (k *Keyboard) AddRaiseRank() *Keyboard {
	return k.AddButton("⬆️", "raiseRole", states.CallbackAdminDetailRaiseRole, "")
}

func (k *Keyboard) AddLowerRank() *Keyboard {
	return k.AddButton("⬇️", "lowerRole", states.CallbackAdminDetailLowerRole, "")
}

func (k *Keyboard) AddRemoveAdminRole() *Keyboard {
	return k.AddButton("❌", "removeAdminRole", states.CallbackAdminDetailRemoveRole, "")
}

func (k *Keyboard) AddFeedbackDelete() *Keyboard {
	return k.AddButton("🗑️", "delete", states.CallbackAdminFeedbackDetailDelete, "")
}

func (k *Keyboard) addSendBroadcast() *Keyboard {
	return k.AddButton("➤", "send", states.CallbackAdminBroadcastSend, "")
}

func userSelectButtons(session *models.Session, users []models.Session, callback string) []Button {
	var buttons []Button

	for i, user := range users {
		itemID := utils.GetItemID(i, session.AdminState.CurrentPage, session.AdminState.PageSize)

		text := fmt.Sprintf("%d - %d", itemID, user.TelegramID)

		if user.TelegramUsername != "" {
			text += fmt.Sprintf(" (@%s)", user.TelegramUsername)
		}

		buttons = append(buttons, Button{"", text, fmt.Sprintf("%s_%d", callback, user.TelegramID), ""})
	}

	return buttons
}
