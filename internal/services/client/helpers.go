package client

import (
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

func GetIntFromMap(data map[string]interface{}, key string, defaultValue int) int {
	if value, ok := data[key].(float64); ok {
		return int(value)
	}
	return defaultValue
}

func GetStringFromMap(data map[string]interface{}, key, defaultValue string) string {
	if value, ok := data[key].(string); ok {
		return value
	}
	return defaultValue
}

func GetFloatFromNestedMap(data map[string]interface{}, key, nestedKey string, defaultValue float64) float64 {
	if nestedMap, ok := data[key].(map[string]interface{}); ok {
		if value, ok := nestedMap[nestedKey].(float64); ok {
			return value
		}
	}
	return defaultValue
}

func GetStringFromNestedMap(data map[string]interface{}, key, nestedKey, defaultValue string) string {
	if nestedMap, ok := data[key].(map[string]interface{}); ok {
		if value, ok := nestedMap[nestedKey].(string); ok {
			return value
		}
	}
	return defaultValue
}

func GetIntFromStringMap(data map[string]interface{}, key string, defaultValue int) int {
	if value, ok := data[key].(string); ok {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func GetFloatFromStringMap(data map[string]interface{}, key string, defaultValue float64) float64 {
	if value, ok := data[key].(string); ok {
		if floatValue, err := strconv.ParseFloat(value, 64); err == nil {
			return floatValue
		}
	}
	return defaultValue
}

func GetTextOrDefault(doc *goquery.Document, selector, defaultValue string) string {
	text := strings.TrimSpace(doc.Find(selector).First().Text())
	if text == "" {
		return defaultValue
	}
	return text
}
