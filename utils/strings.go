package utils

import "unicode"

// IsBlank checks if string is empty ("") or whitespace only.
func IsBlank(s string) bool {
	if s == "" {
		return true
	}

	// checks whitespace only
	for _, v := range s {
		if !unicode.IsSpace(v) {
			return false
		}
	}

	return true
}
