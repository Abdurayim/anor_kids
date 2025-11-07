package utils

import (
	"testing"
)

func TestTruncateText(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		maxLen   int
		expected string
	}{
		{
			name:     "Short text no truncation",
			input:    "Hello",
			maxLen:   10,
			expected: "Hello",
		},
		{
			name:     "Exact length no truncation",
			input:    "Hello",
			maxLen:   5,
			expected: "Hello",
		},
		{
			name:     "Simple truncation",
			input:    "Hello World",
			maxLen:   5,
			expected: "Hello...",
		},
		{
			name:     "Cyrillic text no truncation",
			input:    "Привет мир",
			maxLen:   15,
			expected: "Привет мир",
		},
		{
			name:     "Cyrillic text truncation",
			input:    "Привет мир",
			maxLen:   6,
			expected: "Привет...",
		},
		{
			name:     "Uzbek Cyrillic truncation",
			input:    "Ўзбекистон Республикаси",
			maxLen:   10,
			expected: "Ўзбекистон...",
		},
		{
			name:     "Mixed Latin and Cyrillic",
			input:    "Hello Привет Ўзбек",
			maxLen:   12,
			expected: "Hello Привет...",
		},
		{
			name:     "Text with emojis",
			input:    "Hello 👋 World 🌍",
			maxLen:   8,
			expected: "Hello 👋 ...",
		},
		{
			name:     "Long Uzbek complaint preview",
			input:    "Мактабда ўқитувчилар билан муаммо бор, улар болаларга яхши муносабатда бўлмаяпти",
			maxLen:   50,
			expected: "Мактабда ўқитувчилар билан муаммо бор, улар болала...",
		},
		{
			name:     "Empty string",
			input:    "",
			maxLen:   10,
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := TruncateText(tt.input, tt.maxLen)
			if result != tt.expected {
				t.Errorf("TruncateText(%q, %d) = %q, want %q", tt.input, tt.maxLen, result, tt.expected)
				t.Logf("Result length: %d runes, Expected length: %d runes", len([]rune(result)), len([]rune(tt.expected)))
			}
		})
	}
}

// TestTruncateText_NoUnicodeSplit verifies that we don't split Unicode characters
func TestTruncateText_NoUnicodeSplit(t *testing.T) {
	// This text has multi-byte Unicode characters
	text := "Ўзбекистон 🇺🇿"

	// Truncate to various lengths and ensure no panics or broken characters
	for i := 1; i <= len([]rune(text))+5; i++ {
		result := TruncateText(text, i)
		// Should not panic and should be valid UTF-8
		_ = result

		// Verify the result is valid UTF-8
		if !isValidUTF8(result) {
			t.Errorf("TruncateText produced invalid UTF-8 for length %d: %q", i, result)
		}
	}
}

// isValidUTF8 checks if a string is valid UTF-8
func isValidUTF8(s string) bool {
	// Try to convert to runes and back
	return s == string([]rune(s))
}

func TestFormatPhoneNumber(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid Uzbek phone",
			input:    "+998901234567",
			expected: "+998 90 123 45 67",
		},
		{
			name:     "Invalid length",
			input:    "+99890123",
			expected: "+99890123",
		},
		{
			name:     "Empty string",
			input:    "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := FormatPhoneNumber(tt.input)
			if result != tt.expected {
				t.Errorf("FormatPhoneNumber(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestSanitizeClassName(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "Valid class name",
			input:    "9A",
			expected: "9A",
		},
		{
			name:     "Class with extra spaces",
			input:    "  9A  ",
			expected: "9A",
		},
		{
			name:     "Class with multiple spaces",
			input:    "9   A",
			expected: "9 A",
		},
		{
			name:     "Empty string",
			input:    "   ",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := SanitizeClassName(tt.input)
			if result != tt.expected {
				t.Errorf("SanitizeClassName(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}
