package keyboards

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/k4sper1love/watchlist-bot/internal/handlers/states"
	"github.com/k4sper1love/watchlist-bot/internal/models"
	"github.com/k4sper1love/watchlist-bot/pkg/roles"
)

// Predefined buttons for the main menu.
var menuButtons = []Button{
	{"👤", "profile", states.CallMenuProfile, "", true},
	{"🎥", "films", states.CallMenuFilms, "", true},
	{"📚", "collections", states.CallMenuCollections, "", true},
	{"⚙️", "settings", states.CallMenuSettings, "", true},
	{"💬", "feedback", states.CallMenuFeedback, "", true},
	{"🚪", "logout", states.CallMenuLogout, "", true},
}

// Predefined buttons for feedback categories.
var feedbackCategoryButtons = []Button{
	{"💡", "suggestions", states.CallFeedbackCategorySuggestions, "", true},
	{"🐞", "bugs", states.CallFeedbackCategoryBugs, "", true},
	{"❓", "issues", states.CallFeedbackCategoryIssues, "", true},
}

// Cancel creates an inline keyboard with a cancel button.
func Cancel(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().AddCancel().Build(session.Lang)
}

// Back creates an inline keyboard with a back button.
func Back(session *models.Session, callback string) *tgbotapi.InlineKeyboardMarkup {
	return New().AddBack(callback).Build(session.Lang)
}

// Survey creates an inline keyboard with yes/no buttons for surveys or confirmations.
func Survey(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().AddSurvey().Build(session.Lang)
}

// SkipAndCancel creates an inline keyboard with skip and cancel buttons.
func SkipAndCancel(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddSkip().
		AddCancel().
		Build(session.Lang)
}

// SurveyAndCancel creates an inline keyboard with yes/no buttons and a cancel button.
func SurveyAndCancel(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddSurvey().
		AddCancel().
		Build(session.Lang)
}

// LanguageSelect creates an inline keyboard for selecting a language.
func LanguageSelect(languages []string) *tgbotapi.InlineKeyboardMarkup {
	return New().AddLanguageSelect(languages, states.SelectStartLang).Build("")
}

// Menu creates an inline keyboard for the main menu.
// Includes admin panel button if the user has helper access, followed by predefined menu buttons.
func Menu(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddIf(session.Role.HasAccess(roles.Helper), func(k *Keyboard) {
			k.AddAdminPanel()
		}).
		AddButtons(menuButtons...).
		Build(session.Lang)
}

// Feedback creates an inline keyboard for selecting feedback categories.
func Feedback(session *models.Session) *tgbotapi.InlineKeyboardMarkup {
	return New().
		AddButtons(feedbackCategoryButtons...).
		AddBack("").
		Build(session.Lang)
}
