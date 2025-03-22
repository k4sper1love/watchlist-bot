// Package translator provides utilities for internationalization (i18n) and localization.
//
// Features:
// - Translation File Loading: Loads JSON-based translation files.
// - Language Support: Supports multiple languages with fallback to English ("en").
// - Thread-Safe Localizer Management: Uses a sync.Map to store localizers.
// - Message Translation: Uses message IDs, template data, and pluralization.
// - Error Handling: Logs warnings when translations are missing.
//
// Usage:
// To initialize the package, use `Init(localeDir)`, then call `Translate`:
//
//	err := translator.Init("locales")
//	if err != nil {
//	    log.Fatalf("Failed to initialize translator: %v", err)
//	}
//
//	message := translator.Translate("fr", "welcome_message", map[string]interface{}{"Name": "John"}, nil)
//	fmt.Println(message) // Output: "Bienvenue, John!"
//
// Translation Files:
// Files must be in JSON format, named as "<language_code>.json" (e.g., "en.json", "fr.json").
//
//	{
//	  "welcome_message": {
//	    "one": "Welcome, {{.Name}}!",
//	    "other": "Welcome, everyone!",
//	  }
//	}
//
// Fallback Mechanism:
// If a translation is missing, it falls back to English ("en"). If no translation exists, the message ID is returned.
//
// This package simplifies multilingual support in applications.
package translator
