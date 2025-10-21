package utils

import (
	"strings"
)

// CleanStringForAnalysis removes extra spaces and normalizes for analysis
func CleanStringForAnalysis(s string) string {
	return strings.TrimSpace(s)
}

// ExtractFirstRune extracts the first character from a string
func ExtractFirstRune(s string) string {
	if s == "" {
		return ""
	}
	for _, r := range s {
		return string(r)
	}
	return ""
}

// IsValidString checks if the string is valid for analysis
func IsValidString(s string) bool {
	return len(strings.TrimSpace(s)) > 0
}
