package validator

import (
	"fmt"
	"strings"
)

// Common validation functions

// IsEmpty checks if string is empty or contains only whitespace
func IsEmpty(s string) bool {
	return strings.TrimSpace(s) == ""
}

// ValidateLanguage validates language code
func ValidateLanguage(lang string) (string, error) {
	lang = strings.ToLower(strings.TrimSpace(lang))

	if lang != "uz" && lang != "ru" {
		return "", fmt.Errorf("faqat uz yoki ru tili qo'llab-quvvatlanadi / поддерживаются только uz или ru языки")
	}

	return lang, nil
}

// ContainsDangerousPatterns checks for dangerous patterns
func ContainsDangerousPatterns(text string) bool {
	dangerous := []string{
		"<script", "javascript:", "onerror=", "onclick=",
		"eval(", "exec(", "../", "..\\",
	}

	lowerText := strings.ToLower(text)
	for _, pattern := range dangerous {
		if strings.Contains(lowerText, pattern) {
			return true
		}
	}

	return false
}

// SanitizeFilename sanitizes filename to prevent path traversal
func SanitizeFilename(filename string) string {
	// Remove path separators
	filename = strings.ReplaceAll(filename, "/", "_")
	filename = strings.ReplaceAll(filename, "\\", "_")
	filename = strings.ReplaceAll(filename, "..", "_")

	// Remove dangerous characters
	dangerous := []string{"<", ">", ":", "\"", "|", "?", "*"}
	for _, char := range dangerous {
		filename = strings.ReplaceAll(filename, char, "_")
	}

	return strings.TrimSpace(filename)
}

// ValidateAnnouncementImageType validates image type for announcements
// Rejects GIF and other non-image formats
// Supports: HEIC, JPG, JPEG, PNG, WEBP, BMP, TIFF
func ValidateAnnouncementImageType(mimeType string) error {
	if mimeType == "" {
		return nil // Allow empty mime type for backward compatibility
	}

	mimeType = strings.ToLower(strings.TrimSpace(mimeType))

	// Reject GIF explicitly
	if strings.Contains(mimeType, "gif") {
		return fmt.Errorf("GIF formati qo'llab-quvvatlanmaydi / Формат GIF не поддерживается")
	}

	// List of supported image MIME types
	supportedTypes := []string{
		"image/jpeg",
		"image/jpg",
		"image/png",
		"image/heic",
		"image/heif",
		"image/webp",
		"image/bmp",
		"image/tiff",
		"image/x-icon",
	}

	// Check if mime type is in supported list
	for _, supported := range supportedTypes {
		if mimeType == supported {
			return nil
		}
	}

	// If starts with "image/" but not in our list, reject
	if strings.HasPrefix(mimeType, "image/") {
		return fmt.Errorf("Bu rasm formati qo'llab-quvvatlanmaydi. Faqat JPG, PNG, HEIC formatlarini yuklang / Этот формат изображения не поддерживается. Загрузите только JPG, PNG, HEIC")
	}

	// Not an image at all
	return fmt.Errorf("Iltimos, rasm yuboring (JPG, PNG, HEIC). Video yoki boshqa fayllar qabul qilinmaydi / Пожалуйста, отправьте изображение (JPG, PNG, HEIC). Видео или другие файлы не принимаются")
}
