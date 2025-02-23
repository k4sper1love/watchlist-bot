package utils

import (
	"net/url"
	"unicode/utf8"
)

func ValidNumberRange[T int | float64](value T, minValue T, maxValue T) bool {
	return value >= minValue && value <= maxValue
}

func ValidStringLength(value string, minLength int, maxLength int) bool {
	return ValidNumberRange(utf8.RuneCountInString(value), minLength, maxLength)
}

func ValidURL(u string) bool {
	_, err := url.ParseRequestURI(u)
	return err == nil
}
