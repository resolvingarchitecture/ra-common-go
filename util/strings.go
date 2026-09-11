package util

import (
	"strings"
	"unicode"
)

// CapitalizeFirst uppercases the first character, leaving the rest untouched.
func CapitalizeFirst(text string) string {
	if text == "" {
		return ""
	}
	return strings.ToUpper(text[:1]) + text[1:]
}

// Capitalize uppercases the first character and every character following a space.
func Capitalize(text string) string {
	var out strings.Builder
	capitalizeNext := true
	for _, ch := range text {
		if capitalizeNext && ch != ' ' {
			out.WriteRune(unicode.ToUpper(ch))
			capitalizeNext = false
		} else {
			out.WriteRune(ch)
			capitalizeNext = ch == ' '
		}
	}
	return out.String()
}
