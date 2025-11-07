package models

import (
	"encoding/json"
	"time"
)

// UserState represents the conversation state of a user
type UserState struct {
	TelegramID int64           `json:"telegram_id" db:"telegram_id"`
	State      string          `json:"state" db:"state"`
	Data       json.RawMessage `json:"data" db:"data"`
	UpdatedAt  time.Time       `json:"updated_at" db:"updated_at"`
}

// ImageData represents a single image in the complaint
type ImageData struct {
	FileID       string `json:"file_id"`
	FileUniqueID string `json:"file_unique_id"`
	FileSize     int    `json:"file_size"`
	MimeType     string `json:"mime_type,omitempty"`
}

// StateData is a helper struct for storing state data
type StateData struct {
	PhoneNumber   string      `json:"phone_number,omitempty"`
	ChildName     string      `json:"child_name,omitempty"`
	ChildClass    string      `json:"child_class,omitempty"`
	Language      string      `json:"language,omitempty"`
	ComplaintText string      `json:"complaint_text,omitempty"`
	Images        []ImageData `json:"images,omitempty"` // Array of images for the complaint
}

// State constants
const (
	StateStart               = "start"
	StateAwaitingLanguage    = "awaiting_language"
	StateAwaitingPhone       = "awaiting_phone"
	StateAwaitingChildName   = "awaiting_child_name"
	StateAwaitingChildClass  = "awaiting_child_class"
	StateRegistered          = "registered"
	StateAwaitingComplaint   = "awaiting_complaint"
	StateAwaitingImages      = "awaiting_images"       // New state for collecting images
	StateConfirmingComplaint = "confirming_complaint"
	StateAwaitingAdminPhone  = "awaiting_admin_phone"
	StateAwaitingClassName   = "awaiting_class_name"
)
