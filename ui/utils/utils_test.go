package utils

import "testing"

func TestInitials(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty", "", ""},
		{"single word", "Alice", "A"},
		{"two words", "Alice Bob", "AB"},
		{"three words", "Alice Bob Carol", "AB"},
		{"spaces", "  Alice   Bob  ", "AB"},
		{"unicode", "Élodie Brûlé", "ÉB"},
		{"hyphenated", "Mary-Jane Watson", "MW"},
		{"lowercase", "john doe", "JD"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Initials(tt.input)
			if got != tt.expected {
				t.Errorf("Initials(%q) = %q, want %q", tt.input, got, tt.expected)
			}
		})
	}
}
