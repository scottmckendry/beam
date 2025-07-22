// Package utils is primarily for utility functions and types for the UI package
// plus anything else that doesn't fit elsewhere.
package utils

import "strings"

// Initials returns the uppercase initials of a name (max 2 letters)
func Initials(name string) string {
	parts := strings.Fields(name)
	initials := ""
	for i, part := range parts {
		if i > 1 {
			break
		}
		if len(part) > 0 {
			initials += string([]rune(part)[0])
		}
	}
	return strings.ToUpper(initials)
}
