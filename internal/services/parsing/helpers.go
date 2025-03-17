package parsing

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

// getIntFromMap extracts an integer value from a map using the specified key.
// If the key is not found or the value is not a float64, it returns the default value.
func getIntFromMap(data map[string]interface{}, key string, defaultValue int) int {
	if value, ok := data[key].(float64); ok {
		return int(value)
	}
	return defaultValue
}

// getStringFromMap extracts a string value from a map using the specified key.
// If the key is not found or the value is not a string, it returns the default value.
func getStringFromMap(data map[string]interface{}, key, defaultValue string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return defaultValue
}

// getFloatFromNestedMap extracts a float64 value from a nested map using the specified keys.
// If the keys are not found or the value is not a float64, it returns the default value.
func getFloatFromNestedMap(data map[string]interface{}, key, nestedKey string, defaultValue float64) float64 {
	if nestedMap, ok := data[key].(map[string]interface{}); ok {
		if value, ok := nestedMap[nestedKey].(float64); ok {
			return value
		}
	}
	return defaultValue
}

// getStringFromNestedMap extracts a string value from a nested map using the specified keys.
// If the keys are not found or the value is not a string, it returns the default value.
func getStringFromNestedMap(data map[string]interface{}, key, nestedKey, defaultValue string) string {
	if nestedMap, ok := data[key].(map[string]interface{}); ok {
		if value, ok := nestedMap[nestedKey].(string); ok {
			return value
		}
	}
	return defaultValue
}

// getIntFromStringMap extracts an integer value from a map where the value is stored as a string.
// If the key is not found or the value cannot be converted to an integer, it returns the default value.
func getIntFromStringMap(data map[string]interface{}, key string, defaultValue int) int {
	if value, ok := data[key].(string); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

// getFloatFromStringMap extracts a float64 value from a map where the value is stored as a string.
// If the key is not found or the value cannot be converted to a float64, it returns the default value.
func getFloatFromStringMap(data map[string]interface{}, key string, defaultValue float64) float64 {
	if value, ok := data[key].(string); ok {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

// getTextOrDefault extracts text from an HTML document using the specified CSS selector.
// If no text is found or the text is empty, it returns the default value.
func getTextOrDefault(doc *goquery.Document, selector, defaultValue string) string {
	text := strings.TrimSpace(doc.Find(selector).First().Text())
	if text == "" {
		return defaultValue
	}
	return text
}
