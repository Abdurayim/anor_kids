package models

import "time"

// Announcement represents an announcement posted by admin
type Announcement struct {
	ID                  int       `json:"id" db:"id"`
	AdminID             int       `json:"admin_id" db:"admin_id"`
	Title               string    `json:"title" db:"title"`
	AnnouncementText    string    `json:"announcement_text" db:"announcement_text"`
	ImageTelegramFileID string    `json:"image_telegram_file_id" db:"image_telegram_file_id"`
	ImageFileUniqueID   string    `json:"image_file_unique_id" db:"image_file_unique_id"`
	ImageFileSize       int       `json:"image_file_size" db:"image_file_size"`
	ImageMimeType       string    `json:"image_mime_type" db:"image_mime_type"`
	IsDocument          bool      `json:"is_document" db:"is_document"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
}

// AnnouncementWithAdmin represents an announcement with admin information (from view)
type AnnouncementWithAdmin struct {
	ID                  int       `json:"id" db:"id"`
	AdminID             int       `json:"admin_id" db:"admin_id"`
	Title               string    `json:"title" db:"title"`
	AnnouncementText    string    `json:"announcement_text" db:"announcement_text"`
	ImageTelegramFileID string    `json:"image_telegram_file_id" db:"image_telegram_file_id"`
	ImageFileUniqueID   string    `json:"image_file_unique_id" db:"image_file_unique_id"`
	ImageFileSize       int       `json:"image_file_size" db:"image_file_size"`
	ImageMimeType       string    `json:"image_mime_type" db:"image_mime_type"`
	IsDocument          bool      `json:"is_document" db:"is_document"`
	CreatedAt           time.Time `json:"created_at" db:"created_at"`
	UpdatedAt           time.Time `json:"updated_at" db:"updated_at"`
	AdminPhone          string    `json:"admin_phone" db:"admin_phone"`
	AdminTelegramID     int64     `json:"admin_telegram_id" db:"admin_telegram_id"`
	AdminUsername       string    `json:"admin_username" db:"admin_username"`
}

// CreateAnnouncementRequest is the request to create a new announcement
type CreateAnnouncementRequest struct {
	AdminID             int    `json:"admin_id" validate:"required"`
	Title               string `json:"title" validate:"required,min=3,max=200"`
	AnnouncementText    string `json:"announcement_text" validate:"required,min=10,max=5000"`
	ImageTelegramFileID string `json:"image_telegram_file_id" validate:"required"`
	ImageFileUniqueID   string `json:"image_file_unique_id"`
	ImageFileSize       int    `json:"image_file_size"`
	ImageMimeType       string `json:"image_mime_type"`
	IsDocument          bool   `json:"is_document"`
}

// UpdateAnnouncementRequest is the request to update an announcement
type UpdateAnnouncementRequest struct {
	ID                  int    `json:"id" validate:"required"`
	Title               string `json:"title" validate:"required,min=3,max=200"`
	AnnouncementText    string `json:"announcement_text" validate:"required,min=10,max=5000"`
	ImageTelegramFileID string `json:"image_telegram_file_id"`
	ImageFileUniqueID   string `json:"image_file_unique_id"`
	ImageFileSize       int    `json:"image_file_size"`
	ImageMimeType       string `json:"image_mime_type"`
}
