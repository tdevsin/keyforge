package utils

import "strings"

// IsEmpty checks if string contains only spaces
func IsEmpty(s string) bool {
	return len(strings.TrimSpace(s)) == 0
}
