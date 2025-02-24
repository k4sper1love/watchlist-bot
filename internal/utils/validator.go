package utils

import (
	"net/url"
	"unicode/utf8"
)

func IsValidNumberRange[T int | float64](value T, minValue T, maxValue T) bool {
	return value >= minValue && value <= maxValue
}

func IsValidStringLength(value string, minLength int, maxLength int) bool {
	return IsValidNumberRange(utf8.RuneCountInString(value), minLength, maxLength)
}

func IsValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
