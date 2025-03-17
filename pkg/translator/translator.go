package translator

import (
	"encoding/json"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

var (
	bundle     *i18n.Bundle // Central i18n bundle that stores all loaded translation messages.
	localizers sync.Map     // Thread-safe map to cache localizers for different languages.
	once       sync.Once    // Ensures the initialization of the bundle happens only once.
)

// Init initializes the translator by loading translation files from the specified directory.
// It sets up the `i18n.Bundle` with English as the default language and registers the JSON unmarshal function.
// Returns an error if reading the directory or loading a translation file fails.
func Init(localeDir string) error {
	var err error
	once.Do(func() {
		// Initialize the i18n bundle with English as the default language.
		bundle = i18n.NewBundle(language.English)
		bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

		// Load all translation files from the specified directory.
		err = loadLocales(bundle, localeDir)
	})
	return err
}

// loadLocales loads all JSON translation files from the specified directory into the i18n bundle.
// Ignores non-JSON files and logs warnings for files that fail to load.
func loadLocales(bundle *i18n.Bundle, dir string) error {
	// Read the directory containing translation files.
	files, err := os.ReadDir(dir)
	if err != nil {
		sl.Log.Error("failed to read directory", slog.Any("error", err), slog.String("dir", dir))
		return err
	}

	for _, file := range files {
		// Skip files that are not in JSON format.
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filePath := filepath.Join(dir, file.Name())
		sl.Log.Info("loading translation file", slog.String("file", filePath))

		// Attempt to load the translation file into the bundle.
		if _, err = bundle.LoadMessageFile(filePath); err != nil {
			sl.Log.Warn("failed to load translation file", slog.Any("error", err), slog.String("file", filePath))
		}
	}

	return nil
}

// getLocalizer retrieves or creates a localizer for the specified language code.
// If a localizer already exists for the language, it is returned from the cache.
// Otherwise, a new localizer is created and stored in the cache.
func getLocalizer(languageCode string) *i18n.Localizer {
	// Check if a localizer for the language already exists in the cache.
	if localizer, ok := localizers.Load(languageCode); ok {
		return localizer.(*i18n.Localizer)
	}

	// Create a new localizer with the specified language and English as the fallback.
	localizer := i18n.NewLocalizer(bundle, languageCode, "en")
	localizers.Store(languageCode, localizer)
	return localizer
}

// Translate translates a message into the specified language using the provided message ID, template data, and plural count.
// If the translation is missing for the requested language, it falls back to English (`en`).
// If the translation is also missing in English, the message ID is returned as the fallback value.
func Translate(languageCode string, messageID string, templateData map[string]interface{}, pluralCount interface{}) string {
	// Retrieve the localizer for the specified language.
	localizer := getLocalizer(languageCode)

	// Attempt to localize the message using the provided configuration.
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
		PluralCount:  pluralCount,
	})

	// If localization fails or the message is empty, fall back to English.
	if err != nil || msg == "" {
		sl.Log.Warn("translation missing, falling back to 'en'", slog.String("lang", languageCode), slog.String("message", messageID))
		fallbackLocalizer := getLocalizer("en")

		msg, err = fallbackLocalizer.Localize(&i18n.LocalizeConfig{
			MessageID:    messageID,
			TemplateData: templateData,
			PluralCount:  pluralCount,
		})

		// If fallback localization also fails, return the message ID as the fallback value.
		if err != nil || msg == "" {
			sl.Log.Warn("translation missing in fallback language", slog.String("lang", "en"), slog.String("message", messageID))
			return messageID
		}
	}

	return msg
}
