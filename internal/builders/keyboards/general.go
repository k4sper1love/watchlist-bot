package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
)

var menuButtons = []Button{
	{"👤", "profile", states.CallMenuProfile, "", true},
	{"🎥", "films", states.CallMenuFilms, "", true},
	{"📚", "collections", states.CallMenuCollections, "", true},
	{"⚙️", "settings", states.CallMenuSettings, "", true},
	{"💬", "feedback", states.CallMenuFeedback, "", true},
	{"🚪", "logout", states.CallMenuLogout, "", true},
}

var feedbackCategoryButtons = []Button{
	{"💡", "suggestions", states.CallFeedbackCategorySuggestions, "", true},
	{"🐞", "bugs", states.CallFeedbackCategoryBugs, "", true},
	{"❓", "issues", states.CallFeedbackCategoryIssues, "", true},
}

func Cancel(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().AddCancel().Build(session.Lang)
}

func Back(session *models.Session, callback string) *tgbotapi.InlineKeyboardMarkup {
	return New().AddBack(callback).Build(session.Lang)
}

func Survey(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().AddSurvey().Build(session.Lang)
}

func SkipAndCancel(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddSkip().
		AddCancel().
		Build(session.Lang)
}

func SurveyAndCancel(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddSurvey().
		AddCancel().
		Build(session.Lang)
}

func LanguageSelect(languages []string) *tgbotapi.InlineKeyboardMarkup {
	return New().AddLanguageSelect(languages, states.SelectStartLang).Build("")
}

func Menu(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(session.Role.HasAccess(roles.Helper), func(k *Keyboard) {
			k.AddAdminPanel()
		}).
		AddButtons(menuButtons...).
		Build(session.Lang)
}

func Feedback(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddButtons(feedbackCategoryButtons...).
		AddBack("").
		Build(session.Lang)
}
