package translator

import (
	"encoding/json"
	"fmt"
	"github.com/k4sper1love/watchlist-api/pkg/logger/sl"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
	"log"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
)

var (
	bundle      *i18n.Bundle
	localizers  sync.Map
	initialized bool
)

func InitTranslator(localeDir string) error {
	if initialized {
		return nil
	}

	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)

	if err := loadLocales(bundle, localeDir); err != nil {
		sl.Log.Error("failed to load locale files", slog.Any("error", err))
		return err
	}

	initialized = true
	return nil
}

func loadLocales(bundle *i18n.Bundle, dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("failed to read directory: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) == ".json" {
			filePath := filepath.Join(dir, file.Name())
			log.Printf("Loading translation file: %s", filePath)

			_, err := bundle.LoadMessageFile(filePath)
			if err != nil {
				log.Printf("Failed to load translation file %s: %v", filePath, err)
				continue
			}
		}
	}

	return nil
}

func getLocalizer(languageCode string) *i18n.Localizer {
	if localizer, ok := localizers.Load(languageCode); ok {
		return localizer.(*i18n.Localizer)
	}

	localizer := i18n.NewLocalizer(bundle, languageCode) // default language
	localizers.Store(languageCode, localizer)
	return localizer
}

func Translate(languageCode string, messageID string, templateData map[string]interface{}, pluralCount interface{}) string {
	localizer := getLocalizer(languageCode)

	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    messageID,
		TemplateData: templateData,
		PluralCount:  pluralCount,
	})

	if err != nil || msg == "" {
		log.Printf("Translation missing for '%s' in language '%s', falling back to 'en'", messageID, languageCode)

		fallbackLocalizer := getLocalizer("en")

		msg, err = fallbackLocalizer.Localize(&i18n.LocalizeConfig{
			MessageID:    messageID,
			TemplateData: templateData,
			PluralCount:  pluralCount,
		})

		if err != nil || msg == "" {
			log.Printf("Translation missing for '%s' in fallback language 'en'", messageID)
			return messageID
		}
	}

	return msg
}
