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

type FilterRangeConfig struct {
	MinValue float64
	MaxValue float64
}

func IsValidNumberRange[T int | float64](value T, minValue T, maxValue T) bool {
	return value >= minValue && value <= maxValue
}

func IsValidStringLength(value string, minLength int, maxLength int) bool {
	return IsValidNumberRange(utf8.RuneCountInString(value), minLength, maxLength)
}

func IsValidURL(u string, minLength int, maxLength int) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil && IsValidStringLength(u, minLength, maxLength)
}

func IsValidEmail(email string, minLength int, maxLength int) bool {
	_, err := mail.ParseAddress(email)
	return err == nil && IsValidStringLength(email, minLength, maxLength)
}

func ValidateFiltersRange(input string, config FilterRangeConfig) (string, error) {
	singleValuePattern := `^\d+(\.\d+)?$`
	rangeValuePattern := `^\d+(\.\d+)?-\d+(\.\d+)?$`
	incompleteRangePattern := `^(-?\d+(\.\d+)?-|-?\d+(\.\d+)?)$`

	input = strings.TrimSpace(input)

	if match, _ := regexp.MatchString(singleValuePattern, input); match {
		return validateSingleValue(input, config)
	}

	if match, _ := regexp.MatchString(rangeValuePattern, input); match {
		return validateRangeValue(input, config)
	}

	if match, _ := regexp.MatchString(incompleteRangePattern, input); match {
		return validateIncompleteRange(input, config)
	}

	return "", fmt.Errorf("invalid filter format: %s", input)
}

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
