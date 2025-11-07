package validator

import (
	"testing"
)

func TestValidateName(t *testing.T) {
	tests := []struct {
		name      string
		input     string
		wantError bool
		desc      string
	}{
		// Valid names with straight apostrophe and hyphen
		{"valid_latin", "John", false, "Simple Latin name"},
		{"valid_cyrillic", "Иван", false, "Simple Cyrillic name"},
		{"valid_uzbek", "Otabek", false, "Uzbek name"},
		{"valid_apostrophe", "O'tkir", false, "Name with straight apostrophe"},
		{"valid_hyphen", "Anne-Marie", false, "Name with straight hyphen"},

		// Valid names with curly quotes (smart quotes from mobile keyboards)
		{"valid_left_curly", "O'tkir", false, "Name with left curly apostrophe"},
		{"valid_right_curly", "O'tkir", false, "Name with right curly apostrophe"},

		// Valid names with special hyphens
		{"valid_en_dash", "Anne–Marie", false, "Name with en-dash"},
		{"valid_em_dash", "Anne—Marie", false, "Name with em-dash"},

		// Valid names with Uzbek Cyrillic characters
		{"valid_uzbek_cyrillic", "Ўлуғбек", false, "Name with Uzbek Ў"},
		{"valid_uzbek_cyrillic2", "Қобил", false, "Name with Uzbek Қ"},
		{"valid_uzbek_cyrillic3", "Ҳасан", false, "Name with Uzbek Ҳ"},

		// Valid names with spaces
		{"valid_with_space", "Anna Maria", false, "Name with space"},
		{"valid_with_multiple_spaces", "John Paul Jones", false, "Name with multiple spaces"},

		// Invalid names with numbers
		{"invalid_number", "John123", true, "Name with numbers"},
		{"invalid_number2", "123John", true, "Name starting with number"},

		// Invalid names with special characters
		{"invalid_special", "John@Doe", true, "Name with @"},
		{"invalid_special2", "John_Doe", true, "Name with underscore"},
		{"invalid_special3", "John+Doe", true, "Name with plus"},
		{"invalid_special4", "John%Doe", true, "Name with percent"},

		// Invalid names - too short
		{"invalid_short", "A", true, "Name too short (1 character)"},

		// Empty/whitespace
		{"invalid_empty", "", true, "Empty name"},
		{"invalid_whitespace", "   ", true, "Only whitespace"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ValidateName(tt.input)

			if tt.wantError {
				if err == nil {
					t.Errorf("ValidateName(%q) expected error, got nil. Desc: %s", tt.input, tt.desc)
				}
			} else {
				if err != nil {
					t.Errorf("ValidateName(%q) unexpected error: %v. Desc: %s", tt.input, err, tt.desc)
				}
				// Allow for trimming, just check result is not empty
				if len(result) == 0 && len(tt.input) > 0 {
					t.Errorf("ValidateName(%q) returned empty result. Desc: %s", tt.input, tt.desc)
				}
			}
		})
	}
}

// TestValidateName_RealWorldExamples tests real-world examples that might fail
func TestValidateName_RealWorldExamples(t *testing.T) {
	// These are examples that commonly fail due to smart quotes from mobile keyboards
	realWorldNames := []string{
		"O'tkir",      // Right curly apostrophe (most common from iOS/Android)
		"O'tkir",      // Left curly apostrophe
		"O'tkir",      // Straight apostrophe
		"Anne–Marie",  // En-dash
		"Anne—Marie",  // Em-dash
		"Anne-Marie",  // Straight hyphen
		"Ўлуғбек",     // Uzbek Cyrillic
		"Қобил Ҳасан", // Multiple Uzbek characters with space
	}

	for _, name := range realWorldNames {
		t.Run(name, func(t *testing.T) {
			_, err := ValidateName(name)
			if err != nil {
				t.Errorf("ValidateName(%q) should be valid but got error: %v", name, err)
			}
		})
	}
}
