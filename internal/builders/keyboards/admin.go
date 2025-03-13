package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
)

var superAdminButtons = []Button{
	{"🛡️", "admins", states.CallAdminAdmins, "", true},
}

var adminButtons = []Button{
	{"👥", "profile", states.CallAdminUsers, "", true},
	{"📢", "broadcast", states.CallAdminBroadcast, "", true},
}

var helperButtons = []Button{
	{"💬", "feedback", states.CallAdminFeedback, "", true},
}

var rolesButtons = []Button{
	{"👤", "user", states.CallUserDetailRoleSelectUser, "", true},
	{"👷🏼", "helper", states.CallUserDetailRoleSelectHelper, "", true},
	{"👨🏻‍💼", "admin", states.CallUserDetailRoleSelectAdmin, "", true},
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
			k.AddSearch(states.CallEntitiesFind)
		}).
		AddIf(session.AdminState.IsAdmin, func(k *Keyboard) {
			k.AddAdminSelect(session, admins)
		}).
		AddIf(!session.AdminState.IsAdmin, func(k *Keyboard) {
			k.AddUserSelect(session, admins)
		}).
		AddNavigation(session.AdminState.CurrentPage, session.AdminState.LastPage, states.EntitiesPage, true).
		AddBack(states.CallEntitiesBack).
		Build(session.Lang)
}

func FeedbackList(session *models.Session, feedbacks []models.Feedback) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddFeedbackSelect(session, feedbacks).
		AddNavigation(session.AdminState.CurrentPage, session.AdminState.LastPage, states.FeedbacksPage, true).
		AddBack(states.CallFeedbacksBack).
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
		AddBack(states.CallUserDetailBack).
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
		AddBack(states.CallAdminDetailBack).
		Build(session.Lang)
}

func FeedbackDetail(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddFeedbackDelete().
		AddBack(states.CallFeedbackDetailBack).
		Build(session.Lang)
}

func UserRoleSelect(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddButtons(rolesButtons...).
		AddIf(session.Role.HasAccess(roles.Root), func(k *Keyboard) {
			k.AddSuperRole()
		}).
		AddBack(states.CallUserDetailAgain).
		Build(session.Lang)
}

func BroadcastConfirm(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddSendBroadcast().
		AddCancel().
		Build(session.Lang)
}
