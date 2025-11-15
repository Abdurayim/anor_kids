package models

import "time"

// Proposal represents a user proposal with PDF file
type Proposal struct {
	ID                int       `json:"id" db:"id"`
	UserID            int       `json:"user_id" db:"user_id"`
	ProposalText      string    `json:"proposal_text" db:"proposal_text"`
	PDFTelegramFileID string    `json:"pdf_telegram_file_id" db:"pdf_telegram_file_id"`
	PDFFilename       string    `json:"pdf_filename" db:"pdf_filename"`
	CreatedAt         time.Time `json:"created_at" db:"created_at"`
	Status            string    `json:"status" db:"status"`
}

// ProposalImage represents an image attached to a proposal
type ProposalImage struct {
	ID             int       `json:"id" db:"id"`
	ProposalID     int       `json:"proposal_id" db:"proposal_id"`
	TelegramFileID string    `json:"telegram_file_id" db:"telegram_file_id"`
	FileUniqueID   string    `json:"file_unique_id" db:"file_unique_id"`
	FileSize       int       `json:"file_size" db:"file_size"`
	MimeType       string    `json:"mime_type" db:"mime_type"`
	OrderIndex     int       `json:"order_index" db:"order_index"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// ProposalWithUser represents a proposal with user information (from view)
type ProposalWithUser struct {
	ID                 int       `json:"id" db:"id"`
	UserID             int       `json:"user_id" db:"user_id"`
	ProposalText       string    `json:"proposal_text" db:"proposal_text"`
	PDFTelegramFileID  string    `json:"pdf_telegram_file_id" db:"pdf_telegram_file_id"`
	PDFFilename        string    `json:"pdf_filename" db:"pdf_filename"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	Status             string    `json:"status" db:"status"`
	UserTelegramID     int64     `json:"user_telegram_id" db:"user_telegram_id"`
	TelegramUsername   string    `json:"telegram_username" db:"telegram_username"`
	PhoneNumber        string    `json:"phone_number" db:"phone_number"`
	ChildName          string    `json:"child_name" db:"child_name"`
	ChildClass         string    `json:"child_class" db:"child_class"`
}

// ProposalWithImages represents a proposal with its associated images
type ProposalWithImages struct {
	Proposal
	Images []ProposalImage `json:"images"`
}

// CreateProposalRequest is the request to create a new proposal
type CreateProposalRequest struct {
	UserID            int      `json:"user_id" validate:"required"`
	ProposalText      string   `json:"proposal_text" validate:"required,min=10,max=5000"`
	ImageFileIDs      []string `json:"image_file_ids"` // Array of Telegram file IDs for images
	PDFTelegramFileID string   `json:"pdf_telegram_file_id" validate:"required"`
	PDFFilename       string   `json:"pdf_filename" validate:"required"`
}

// ProposalStatus constants (reusing the same status constants)
const (
	ProposalStatusPending  = "pending"
	ProposalStatusReviewed = "reviewed"
	ProposalStatusArchived = "archived"
)
