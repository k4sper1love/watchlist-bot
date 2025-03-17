package utils

import (
	"errors"
	"fmt"
	"net/mail"
	"net/url"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

// FilterRangeConfig defines the configuration for validating numeric ranges.
type FilterRangeConfig struct {
	MinValue float64 // The minimum allowed value in the range.
	MaxValue float64 // The maximum allowed value in the range.
}

// IsValidNumberRange checks if a numeric value is within the specified range [minValue, maxValue].
func IsValidNumberRange[T int | float64](value T, minValue T, maxValue T) bool {
	return value >= minValue && value <= maxValue
}

// IsValidStringLength checks if a string's length (in runes) is within the specified range [minLength, maxLength].
func IsValidStringLength(value string, minLength int, maxLength int) bool {
	return IsValidNumberRange(utf8.RuneCountInString(value), minLength, maxLength)
}

// IsValidURL validates a URL by checking its format and ensuring its length is within the specified range.
func IsValidURL(u string, minLength int, maxLength int) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil && IsValidStringLength(u, minLength, maxLength)
}

// IsValidEmail validates an email address by checking its format and ensuring its length is within the specified range.
func IsValidEmail(email string, minLength int, maxLength int) bool {
	_, err := mail.ParseAddress(email)
	return err == nil && IsValidStringLength(email, minLength, maxLength)
}

// ValidateFiltersRange validates a filter input string based on the provided configuration.
// It supports single values, ranges, and incomplete ranges.
func ValidateFiltersRange(input string, config FilterRangeConfig) (string, error) {
	// Regular expressions for matching different filter formats.
	singleValuePattern := `^\d+(\.\d+)?$`
	rangeValuePattern := `^\d+(\.\d+)?-\d+(\.\d+)?$`
	incompleteRangePattern := `^(-?\d+(\.\d+)?-|-?\d+(\.\d+)?)$`

	input = strings.TrimSpace(input)

	// Check if the input matches a single value pattern.
	if match, _ := regexp.MatchString(singleValuePattern, input); match {
		return validateSingleValue(input, config)
	}

	// Check if the input matches a range value pattern.
	if match, _ := regexp.MatchString(rangeValuePattern, input); match {
		return validateRangeValue(input, config)
	}

	// Check if the input matches an incomplete range pattern.
	if match, _ := regexp.MatchString(incompleteRangePattern, input); match {
		return validateIncompleteRange(input, config)
	}

	// If no pattern matches, return an error.
	return "", fmt.Errorf("invalid filter format: %s", input)
}

// validateSingleValue validates a single numeric value against the provided configuration.
func validateSingleValue(input string, config FilterRangeConfig) (string, error) {
	value, err := strconv.ParseFloat(input, 64)
	if err != nil {
		return "", fmt.Errorf("invalid number: %s", input)
	}
	if value < config.MinValue || value > config.MaxValue {
		return "", fmt.Errorf("value out of range: %s (must be between %.2f and %.2f)", input, config.MinValue, config.MaxValue)
	}
	return input, nil
}

// validateRangeValue validates a numeric range against the provided configuration.
func validateRangeValue(input string, config FilterRangeConfig) (string, error) {
	parts := strings.Split(input, "-")
	start, err1 := strconv.ParseFloat(parts[0], 64)
	end, err2 := strconv.ParseFloat(parts[1], 64)
	if err1 != nil || err2 != nil {
		return "", errors.New("invalid range format")
	}
	if start >= end {
		return "", fmt.Errorf("invalid range: start (%s) must be less than end (%s)", parts[0], parts[1])
	}
	if start < config.MinValue || end > config.MaxValue {
		return "", fmt.Errorf("range out of bounds: %s (must be between %.2f and %.2f)", input, config.MinValue, config.MaxValue)
	}
	return input, nil
}

// validateIncompleteRange validates an incomplete numeric range against the provided configuration.
func validateIncompleteRange(input string, config FilterRangeConfig) (string, error) {
	parts := strings.Split(input, "-")
	if parts[0] != "" && parts[1] == "" {
		start, err := strconv.ParseFloat(parts[0], 64)
		if err != nil {
			return "", fmt.Errorf("invalid start value: %s", parts[0])
		}
		if start < config.MinValue || start > config.MaxValue {
			return "", fmt.Errorf("value out of range: %s (must be between %.2f and %.2f)", parts[0], config.MinValue, config.MaxValue)
		}
		return fmt.Sprintf("%s-%.f", parts[0], config.MaxValue), nil
	}
	if parts[0] == "" && parts[1] != "" {
		end, err := strconv.ParseFloat(parts[1], 64)
		if err != nil {
			return "", fmt.Errorf("invalid end value: %s", parts[1])
		}
		if end < config.MinValue || end > config.MaxValue {
			return "", fmt.Errorf("value out of range: %s (must be between %.2f and %.2f)", parts[1], config.MinValue, config.MaxValue)
		}
		return fmt.Sprintf("%.f-%s", config.MinValue, parts[1]), nil
	}
	return "", fmt.Errorf("invalid incomplete range format: %s", input)
}
