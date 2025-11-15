package utils

import (
	"fmt"
	"strings"
	"time"

	"anor-kids/internal/validator"
)

// GenerateComplaintFilename generates a filename for complaint document (legacy DOCX)
// Format: Shikoyat_ParentName_ClassName_Date.docx
func GenerateComplaintFilename(childName, childClass string) string {
	date := time.Now().Format("2006-01-02")

	// Sanitize name
	safeName := validator.SanitizeFilename(childName)
	safeName = strings.ReplaceAll(safeName, " ", "_")

	// Create filename
	filename := fmt.Sprintf("Shikoyat_%s_%s_sinf_%s.docx", safeName, childClass, date)

	return filename
}

// GeneratePDFFilename generates a filename for complaint PDF document
// Format: Shikoyat_ChildName_ClassName_Date.pdf
func GeneratePDFFilename(childName, childClass string) string {
	date := time.Now().Format("2006-01-02")

	// Sanitize name
	safeName := validator.SanitizeFilename(childName)
	safeName = strings.ReplaceAll(safeName, " ", "_")

	// Create filename
	filename := fmt.Sprintf("Shikoyat_%s_%s_sinf_%s.pdf", safeName, childClass, date)

	return filename
}

// GenerateProposalPDFFilename generates a filename for proposal PDF document
// Format: Taklif_ChildName_ClassName_Date.pdf
func GenerateProposalPDFFilename(childName, childClass string) string {
	date := time.Now().Format("2006-01-02")

	// Sanitize name
	safeName := validator.SanitizeFilename(childName)
	safeName = strings.ReplaceAll(safeName, " ", "_")

	// Create filename
	filename := fmt.Sprintf("Taklif_%s_%s_sinf_%s.pdf", safeName, childClass, date)

	return filename
}

// GenerateComplaintCaption generates caption for complaint document
func GenerateComplaintCaption(childName, childClass, phoneNumber string) string {
	return fmt.Sprintf(
		"YANGI SHIKOYAT / НОВАЯ ЖАЛОБА\n\n"+
			"Ota-ona / Родитель: %s\n"+
			"Sinf / Класс: %s\n"+
			"Telefon / Телефон: %s\n"+
			"Sana / Дата: %s",
		childName,
		childClass,
		phoneNumber,
		time.Now().Format("02.01.2006 15:04"),
	)
}

// TruncateText truncates text to specified length (character count, not bytes)
// Properly handles Unicode characters (Cyrillic, emojis, etc.)
func TruncateText(text string, maxLen int) string {
	runes := []rune(text)
	if len(runes) <= maxLen {
		return text
	}

	return string(runes[:maxLen]) + "..."
}

// FormatPhoneNumber formats phone number for display
func FormatPhoneNumber(phone string) string {
	// +998 90 123 45 67
	if len(phone) != 13 {
		return phone
	}

	return fmt.Sprintf("%s %s %s %s %s",
		phone[:4],   // +998
		phone[4:6],  // 90
		phone[6:9],  // 123
		phone[9:11], // 45
		phone[11:],  // 67
	)
}

// EscapeMarkdown escapes special characters for Telegram Markdown
func EscapeMarkdown(text string) string {
	replacer := strings.NewReplacer(
		"_", "\\_",
		"*", "\\*",
		"[", "\\[",
		"]", "\\]",
		"(", "\\(",
		")", "\\)",
		"~", "\\~",
		"`", "\\`",
		">", "\\>",
		"#", "\\#",
		"+", "\\+",
		"-", "\\-",
		"=", "\\=",
		"|", "\\|",
		"{", "\\{",
		"}", "\\}",
		".", "\\.",
		"!", "\\!",
	)

	return replacer.Replace(text)
}

// EscapeHTML escapes special HTML characters for Telegram HTML mode
func EscapeHTML(text string) string {
	replacer := strings.NewReplacer(
		"&", "&amp;",
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
	)

	return replacer.Replace(text)
}

// FormatDateTime formats datetime for display
func FormatDateTime(t time.Time) string {
	return t.Format("02.01.2006 15:04")
}

// FormatDate formats date for display
func FormatDate(t time.Time) string {
	return t.Format("02.01.2006")
}

// SanitizeClassName sanitizes class name
func SanitizeClassName(className string) string {
	// Remove extra spaces and trim
	className = strings.TrimSpace(className)
	// Remove multiple spaces
	className = strings.Join(strings.Fields(className), " ")
	return className
}

// StripEmojis removes emoji characters from text for PDF compatibility
// DejaVu fonts don't support emojis, so we need to strip them
func StripEmojis(text string) string {
	var result strings.Builder

	for _, r := range text {
		// Keep characters in these ranges (basic Latin, Cyrillic, common punctuation)
		// Basic Latin: 0x0000-0x007F
		// Latin Extended: 0x0080-0x00FF
		// Cyrillic: 0x0400-0x04FF
		// General Punctuation: 0x2000-0x206F
		// Currency: 0x20A0-0x20CF
		if (r >= 0x0020 && r <= 0x007E) || // ASCII printable
		   (r >= 0x00A0 && r <= 0x00FF) || // Latin Extended
		   (r >= 0x0400 && r <= 0x04FF) || // Cyrillic
		   (r >= 0x2000 && r <= 0x206F) || // General Punctuation
		   (r >= 0x20A0 && r <= 0x20CF) || // Currency
		   r == '\n' || r == '\r' || r == '\t' { // Whitespace
			result.WriteRune(r)
		} else {
			// Skip emoji and other unsupported characters
			// Optionally, you could replace with a space or placeholder
			// For now, we just skip them
		}
	}

	return result.String()
}
