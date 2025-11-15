package repository

import (
	"database/sql"
	"fmt"
	"time"

	"anor-kids/internal/models"
)

type AnnouncementRepository struct {
	db *sql.DB
}

func NewAnnouncementRepository(db *sql.DB) *AnnouncementRepository {
	return &AnnouncementRepository{db: db}
}

// Create creates a new announcement
func (r *AnnouncementRepository) Create(req *models.CreateAnnouncementRequest) (*models.Announcement, error) {
	query := `
		INSERT INTO announcements (admin_id, title, announcement_text, image_telegram_file_id, image_file_unique_id, image_file_size, image_mime_type, is_document)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		RETURNING id, admin_id, title, announcement_text, image_telegram_file_id, image_file_unique_id, image_file_size, image_mime_type, is_document, created_at, updated_at
	`

	var announcement models.Announcement
	err := r.db.QueryRow(
		query,
		req.AdminID,
		req.Title,
		req.AnnouncementText,
		req.ImageTelegramFileID,
		req.ImageFileUniqueID,
		req.ImageFileSize,
		req.ImageMimeType,
		req.IsDocument,
	).Scan(
		&announcement.ID,
		&announcement.AdminID,
		&announcement.Title,
		&announcement.AnnouncementText,
		&announcement.ImageTelegramFileID,
		&announcement.ImageFileUniqueID,
		&announcement.ImageFileSize,
		&announcement.ImageMimeType,
		&announcement.IsDocument,
		&announcement.CreatedAt,
		&announcement.UpdatedAt,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create announcement: %w", err)
	}

	return &announcement, nil
}

// GetByID gets announcement by ID
func (r *AnnouncementRepository) GetByID(id int) (*models.Announcement, error) {
	query := `
		SELECT id, admin_id, title, announcement_text, image_telegram_file_id, image_file_unique_id, image_file_size, image_mime_type, is_document, created_at, updated_at
		FROM announcements
		WHERE id = $1
	`

	var announcement models.Announcement
	err := r.db.QueryRow(query, id).Scan(
		&announcement.ID,
		&announcement.AdminID,
		&announcement.Title,
		&announcement.AnnouncementText,
		&announcement.ImageTelegramFileID,
		&announcement.ImageFileUniqueID,
		&announcement.ImageFileSize,
		&announcement.ImageMimeType,
		&announcement.IsDocument,
		&announcement.CreatedAt,
		&announcement.UpdatedAt,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get announcement: %w", err)
	}

	return &announcement, nil
}

// GetAll gets all announcements with pagination (for users)
func (r *AnnouncementRepository) GetAll(limit, offset int) ([]*models.Announcement, error) {
	query := `
		SELECT id, admin_id, title, announcement_text, image_telegram_file_id, image_file_unique_id, image_file_size, image_mime_type, is_document, created_at, updated_at
		FROM announcements
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get announcements: %w", err)
	}
	defer rows.Close()

	var announcements []*models.Announcement
	for rows.Next() {
		var announcement models.Announcement
		err := rows.Scan(
			&announcement.ID,
			&announcement.AdminID,
			&announcement.Title,
			&announcement.AnnouncementText,
			&announcement.ImageTelegramFileID,
			&announcement.ImageFileUniqueID,
			&announcement.ImageFileSize,
			&announcement.ImageMimeType,
			&announcement.IsDocument,
			&announcement.CreatedAt,
			&announcement.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan announcement: %w", err)
		}
		announcements = append(announcements, &announcement)
	}

	return announcements, nil
}

// GetByAdminID gets announcements by admin ID (for admin to see their own posts)
func (r *AnnouncementRepository) GetByAdminID(adminID int, limit, offset int) ([]*models.Announcement, error) {
	query := `
		SELECT id, admin_id, title, announcement_text, image_telegram_file_id, image_file_unique_id, image_file_size, image_mime_type, is_document, created_at, updated_at
		FROM announcements
		WHERE admin_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, adminID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get announcements by admin: %w", err)
	}
	defer rows.Close()

	var announcements []*models.Announcement
	for rows.Next() {
		var announcement models.Announcement
		err := rows.Scan(
			&announcement.ID,
			&announcement.AdminID,
			&announcement.Title,
			&announcement.AnnouncementText,
			&announcement.ImageTelegramFileID,
			&announcement.ImageFileUniqueID,
			&announcement.ImageFileSize,
			&announcement.ImageMimeType,
			&announcement.IsDocument,
			&announcement.CreatedAt,
			&announcement.UpdatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan announcement: %w", err)
		}
		announcements = append(announcements, &announcement)
	}

	return announcements, nil
}

// GetAllWithAdmin gets all announcements with admin info using view
func (r *AnnouncementRepository) GetAllWithAdmin(limit, offset int) ([]*models.AnnouncementWithAdmin, error) {
	query := `
		SELECT id, admin_id, title, announcement_text, image_telegram_file_id, image_file_unique_id, image_file_size, image_mime_type, is_document, created_at, updated_at,
		       admin_phone, admin_telegram_id, admin_username
		FROM v_announcements_with_admin
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get announcements with admin: %w", err)
	}
	defer rows.Close()

	var announcements []*models.AnnouncementWithAdmin
	for rows.Next() {
		var announcement models.AnnouncementWithAdmin
		err := rows.Scan(
			&announcement.ID,
			&announcement.AdminID,
			&announcement.Title,
			&announcement.AnnouncementText,
			&announcement.ImageTelegramFileID,
			&announcement.ImageFileUniqueID,
			&announcement.ImageFileSize,
			&announcement.ImageMimeType,
			&announcement.IsDocument,
			&announcement.CreatedAt,
			&announcement.UpdatedAt,
			&announcement.AdminPhone,
			&announcement.AdminTelegramID,
			&announcement.AdminUsername,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan announcement with admin: %w", err)
		}
		announcements = append(announcements, &announcement)
	}

	return announcements, nil
}

// Update updates an announcement
func (r *AnnouncementRepository) Update(req *models.UpdateAnnouncementRequest) error {
	// Build dynamic query based on what fields are provided
	query := `
		UPDATE announcements
		SET title = $1, announcement_text = $2, updated_at = $3
	`
	args := []interface{}{req.Title, req.AnnouncementText, time.Now()}
	argPos := 4

	// Add image fields if provided
	if req.ImageTelegramFileID != "" {
		query += fmt.Sprintf(", image_telegram_file_id = $%d, image_file_unique_id = $%d, image_file_size = $%d, image_mime_type = $%d",
			argPos, argPos+1, argPos+2, argPos+3)
		args = append(args, req.ImageTelegramFileID, req.ImageFileUniqueID, req.ImageFileSize, req.ImageMimeType)
		argPos += 4
	}

	query += fmt.Sprintf(" WHERE id = $%d", argPos)
	args = append(args, req.ID)

	_, err := r.db.Exec(query, args...)
	if err != nil {
		return fmt.Errorf("failed to update announcement: %w", err)
	}

	return nil
}

// Delete deletes an announcement
func (r *AnnouncementRepository) Delete(id int) error {
	query := `DELETE FROM announcements WHERE id = $1`
	result, err := r.db.Exec(query, id)
	if err != nil {
		return fmt.Errorf("failed to delete announcement: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("announcement not found")
	}

	return nil
}

// Count counts total announcements
func (r *AnnouncementRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM announcements").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count announcements: %w", err)
	}
	return count, nil
}

// CountByAdminID counts announcements by admin ID
func (r *AnnouncementRepository) CountByAdminID(adminID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM announcements WHERE admin_id = $1`
	err := r.db.QueryRow(query, adminID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count admin announcements: %w", err)
	}
	return count, nil
}
