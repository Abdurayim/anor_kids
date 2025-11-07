package models

import "time"

// Complaint represents a user complaint with PDF file
type Complaint struct {
	ID                 int       `json:"id" db:"id"`
	UserID             int       `json:"user_id" db:"user_id"`
	ComplaintText      string    `json:"complaint_text" db:"complaint_text"`
	PDFTelegramFileID  string    `json:"pdf_telegram_file_id" db:"pdf_telegram_file_id"`
	PDFFilename        string    `json:"pdf_filename" db:"pdf_filename"`
	CreatedAt          time.Time `json:"created_at" db:"created_at"`
	Status             string    `json:"status" db:"status"`
}

// ComplaintImage represents an image attached to a complaint
type ComplaintImage struct {
	ID             int       `json:"id" db:"id"`
	ComplaintID    int       `json:"complaint_id" db:"complaint_id"`
	TelegramFileID string    `json:"telegram_file_id" db:"telegram_file_id"`
	FileUniqueID   string    `json:"file_unique_id" db:"file_unique_id"`
	FileSize       int       `json:"file_size" db:"file_size"`
	MimeType       string    `json:"mime_type" db:"mime_type"`
	OrderIndex     int       `json:"order_index" db:"order_index"`
	CreatedAt      time.Time `json:"created_at" db:"created_at"`
}

// ComplaintWithUser represents a complaint with user information (from view)
type ComplaintWithUser struct {
	ID                 int       `json:"id" db:"id"`
	UserID             int       `json:"user_id" db:"user_id"`
	ComplaintText      string    `json:"complaint_text" db:"complaint_text"`
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

// ComplaintWithImages represents a complaint with its associated images
type ComplaintWithImages struct {
	Complaint
	Images []ComplaintImage `json:"images"`
}

// CreateComplaintRequest is the request to create a new complaint
type CreateComplaintRequest struct {
	UserID            int      `json:"user_id" validate:"required"`
	ComplaintText     string   `json:"complaint_text" validate:"required,min=10,max=5000"`
	ImageFileIDs      []string `json:"image_file_ids"` // Array of Telegram file IDs for images
	PDFTelegramFileID string   `json:"pdf_telegram_file_id" validate:"required"`
	PDFFilename       string   `json:"pdf_filename" validate:"required"`
}

// ComplaintStatus constants
const (
	StatusPending  = "pending"
	StatusReviewed = "reviewed"
	StatusArchived = "archived"
)
