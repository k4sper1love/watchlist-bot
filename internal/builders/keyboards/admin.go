package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
)

var superAdminButtons = []Button{
	{"🛡️", "admins", states.CallbackAdminSelectAdmins, "", true},
}

var adminButtons = []Button{
	{"👥", "users", states.CallbackAdminSelectUsers, "", true},
	{"📢", "broadcast", states.CallbackAdminSelectBroadcast, "", true},
}

var helperButtons = []Button{
	{"💬", "feedback", states.CallbackAdminSelectFeedback, "", true},
}

var rolesButtons = []Button{
	{"👤", "user", states.CallbackAdminUserRoleSelectUser, "", true},
	{"👷🏼", "helper", states.CallbackAdminUserRoleSelectHelper, "", true},
	{"👨🏻‍💼", "admin", states.CallbackAdminUserRoleSelectAdmin, "", true},
}

func AdminMenu(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(session.Role.HasAccess(roles.SuperAdmin), func(k *Keyboard) {
			k.AddButtons(superAdminButtons...)
		}).
		AddIf(session.Role.HasAccess(roles.Admin), func(k *Keyboard) {
			k.AddButtons(adminButtons...)
		}).
		AddIf(session.Role.HasAccess(roles.Helper), func(k *Keyboard) {
			k.AddButtons(helperButtons...)
		}).
		AddBack("").
		Build(session.Lang)
}

func AdminList(session *models.Session, admins []models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(len(admins) > 0, func(k *Keyboard) {
			k.AddSearch(states.CallbackEntitiesSelectFind)
		}).
		AddIf(session.AdminState.IsAdmin, func(k *Keyboard) {
			k.AddAdminSelect(session, admins)
		}).
		AddIf(!session.AdminState.IsAdmin, func(k *Keyboard) {
			k.AddUserSelect(session, admins)
		}).
		AddNavigation(session.AdminState.CurrentPage, session.AdminState.LastPage, states.PrefixEntitiesListPage, true).
		AddBack(states.CallbackEntitiesListBack).
		Build(session.Lang)
}

func FeedbackList(session *models.Session, feedbacks []models.Feedback) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddFeedbackSelect(session, feedbacks).
		AddNavigation(session.AdminState.CurrentPage, session.AdminState.LastPage, states.PrefixAdminFeedbackListPage, true).
		AddBack(states.CallbackAdminFeedbackListBack).
		Build(session.Lang)
}

func UserDetail(session *models.Session, user *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(session.Role.HasAccess(roles.SuperAdmin) && !user.Role.HasAccess(session.Role), func(k *Keyboard) {
			k.AddLogs()
		}).
		AddIf(user.IsBanned, func(k *Keyboard) {
			k.AddUnban()
		}).
		AddIf(!user.IsBanned && !user.Role.HasAccess(session.Role), func(k *Keyboard) {
			k.AddBan()
		}).
		AddIf(session.Role.HasAccess(roles.SuperAdmin) && !user.Role.HasAccess(session.Role), func(k *Keyboard) {
			k.AddManageRole()
		}).
		AddBack(states.CallbackAdminUserDetailBack).
		Build(session.Lang)
}

func AdminDetail(session *models.Session, user *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(!user.Role.HasAccess(session.Role.PrevRole()), func(k *Keyboard) {
			k.AddRaiseRank()
		}).
		AddIf(user.Role.HasAccess(roles.Admin) && !user.Role.HasAccess(session.Role), func(k *Keyboard) {
			k.AddLowerRank()
		}).
		AddIf(user.Role.HasAccess(roles.Helper) &&
			session.Role.HasAccess(roles.SuperAdmin) && !user.Role.HasAccess(session.Role), func(k *Keyboard) {
			k.AddRemoveAdminRole()
		}).
		AddBack(states.CallbackAdminDetailBack).
		Build(session.Lang)
}

func FeedbackDetail(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddFeedbackDelete().
		AddBack(states.CallbackAdminFeedbackDetailBack).
		Build(session.Lang)
}

func UserRoleSelect(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddButtons(rolesButtons...).
		AddIf(session.Role.HasAccess(roles.Root), func(k *Keyboard) {
			k.AddSuperRole()
		}).
		AddBack(states.CallbackAdminUserDetailAgain).
		Build(session.Lang)
}

func BroadcastConfirm(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddSendBroadcast().
		AddCancel().
		Build(session.Lang)
}
