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
	IsDocument   bool   `json:"is_document,omitempty"` // True if sent as document (HEIC, etc.)
}

// StateData is a helper struct for storing state data
type StateData struct {
	PhoneNumber        string      `json:"phone_number,omitempty"`
	ChildName          string      `json:"child_name,omitempty"`
	ChildClass         string      `json:"child_class,omitempty"`
	Language           string      `json:"language,omitempty"`
	ComplaintText      string      `json:"complaint_text,omitempty"`
	ProposalText       string      `json:"proposal_text,omitempty"`
	AnnouncementTitle  string      `json:"announcement_title,omitempty"`
	AnnouncementText   string      `json:"announcement_text,omitempty"`
	AnnouncementImage  *ImageData  `json:"announcement_image,omitempty"` // Single image for announcement
	Images             []ImageData `json:"images,omitempty"`             // Array of images for the complaint or proposal
}

// State constants
const (
	StateStart                      = "start"
	StateAwaitingLanguage           = "awaiting_language"
	StateAwaitingPhone              = "awaiting_phone"
	StateAwaitingChildName          = "awaiting_child_name"
	StateAwaitingChildClass         = "awaiting_child_class"
	StateRegistered                 = "registered"
	StateAwaitingComplaint          = "awaiting_complaint"
	StateAwaitingImages             = "awaiting_images"          // State for collecting complaint images
	StateConfirmingComplaint        = "confirming_complaint"
	StateAwaitingProposal           = "awaiting_proposal"
	StateAwaitingProposalImages     = "awaiting_proposal_images" // State for collecting proposal images
	StateConfirmingProposal         = "confirming_proposal"
	StateAwaitingAdminPhone         = "awaiting_admin_phone"
	StateAwaitingClassName          = "awaiting_class_name"
	StateAwaitingAnnouncementTitle  = "awaiting_announcement_title"
	StateAwaitingAnnouncementText   = "awaiting_announcement_text"
	StateAwaitingAnnouncementImage  = "awaiting_announcement_image"
)
