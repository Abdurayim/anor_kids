package repository

import (
	"database/sql"
	"fmt"

	"anor-kids/internal/models"
)

type ComplaintRepository struct {
	db *sql.DB
}

func NewComplaintRepository(db *sql.DB) *ComplaintRepository {
	return &ComplaintRepository{db: db}
}

// Create creates a new complaint with PDF file
func (r *ComplaintRepository) Create(req *models.CreateComplaintRequest) (*models.Complaint, error) {
	query := `
		INSERT INTO complaints (user_id, complaint_text, pdf_telegram_file_id, pdf_filename)
		VALUES ($1, $2, $3, $4)
		RETURNING id, user_id, complaint_text, pdf_telegram_file_id, pdf_filename, created_at, status
	`

	var complaint models.Complaint
	err := r.db.QueryRow(
		query,
		req.UserID,
		req.ComplaintText,
		req.PDFTelegramFileID,
		req.PDFFilename,
	).Scan(
		&complaint.ID,
		&complaint.UserID,
		&complaint.ComplaintText,
		&complaint.PDFTelegramFileID,
		&complaint.PDFFilename,
		&complaint.CreatedAt,
		&complaint.Status,
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create complaint: %w", err)
	}

	return &complaint, nil
}

// GetByID gets complaint by ID
func (r *ComplaintRepository) GetByID(id int) (*models.Complaint, error) {
	query := `
		SELECT id, user_id, complaint_text, pdf_telegram_file_id, pdf_filename, created_at, status
		FROM complaints
		WHERE id = $1
	`

	var complaint models.Complaint
	err := r.db.QueryRow(query, id).Scan(
		&complaint.ID,
		&complaint.UserID,
		&complaint.ComplaintText,
		&complaint.PDFTelegramFileID,
		&complaint.PDFFilename,
		&complaint.CreatedAt,
		&complaint.Status,
	)

	if err == sql.ErrNoRows {
		return nil, nil
	}

	if err != nil {
		return nil, fmt.Errorf("failed to get complaint: %w", err)
	}

	return &complaint, nil
}

// GetByUserID gets complaints by user ID (indexed, fast query)
func (r *ComplaintRepository) GetByUserID(userID int, limit, offset int) ([]*models.Complaint, error) {
	query := `
		SELECT id, user_id, complaint_text, pdf_telegram_file_id, pdf_filename, created_at, status
		FROM complaints
		WHERE user_id = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, userID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get complaints: %w", err)
	}
	defer rows.Close()

	var complaints []*models.Complaint
	for rows.Next() {
		var complaint models.Complaint
		err := rows.Scan(
			&complaint.ID,
			&complaint.UserID,
			&complaint.ComplaintText,
			&complaint.PDFTelegramFileID,
			&complaint.PDFFilename,
			&complaint.CreatedAt,
			&complaint.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan complaint: %w", err)
		}
		complaints = append(complaints, &complaint)
	}

	return complaints, nil
}

// GetAll gets all complaints with pagination (for admin)
func (r *ComplaintRepository) GetAll(limit, offset int) ([]*models.Complaint, error) {
	query := `
		SELECT id, user_id, complaint_text, pdf_telegram_file_id, pdf_filename, created_at, status
		FROM complaints
		ORDER BY created_at DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get complaints: %w", err)
	}
	defer rows.Close()

	var complaints []*models.Complaint
	for rows.Next() {
		var complaint models.Complaint
		err := rows.Scan(
			&complaint.ID,
			&complaint.UserID,
			&complaint.ComplaintText,
			&complaint.PDFTelegramFileID,
			&complaint.PDFFilename,
			&complaint.CreatedAt,
			&complaint.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan complaint: %w", err)
		}
		complaints = append(complaints, &complaint)
	}

	return complaints, nil
}

// GetAllWithUser gets all complaints with user info using view (optimized for admin)
func (r *ComplaintRepository) GetAllWithUser(limit, offset int) ([]*models.ComplaintWithUser, error) {
	query := `
		SELECT id, user_id, complaint_text, pdf_telegram_file_id, pdf_filename, created_at, status,
		       user_telegram_id, telegram_username, phone_number, child_name, child_class
		FROM v_complaints_with_user
		LIMIT $1 OFFSET $2
	`

	rows, err := r.db.Query(query, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get complaints with user: %w", err)
	}
	defer rows.Close()

	var complaints []*models.ComplaintWithUser
	for rows.Next() {
		var complaint models.ComplaintWithUser
		err := rows.Scan(
			&complaint.ID,
			&complaint.UserID,
			&complaint.ComplaintText,
			&complaint.PDFTelegramFileID,
			&complaint.PDFFilename,
			&complaint.CreatedAt,
			&complaint.Status,
			&complaint.UserTelegramID,
			&complaint.TelegramUsername,
			&complaint.PhoneNumber,
			&complaint.ChildName,
			&complaint.ChildClass,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan complaint with user: %w", err)
		}
		complaints = append(complaints, &complaint)
	}

	return complaints, nil
}

// GetByStatus gets complaints by status (indexed, fast query)
func (r *ComplaintRepository) GetByStatus(status string, limit, offset int) ([]*models.Complaint, error) {
	query := `
		SELECT id, user_id, complaint_text, pdf_telegram_file_id, pdf_filename, created_at, status
		FROM complaints
		WHERE status = $1
		ORDER BY created_at DESC
		LIMIT $2 OFFSET $3
	`

	rows, err := r.db.Query(query, status, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get complaints by status: %w", err)
	}
	defer rows.Close()

	var complaints []*models.Complaint
	for rows.Next() {
		var complaint models.Complaint
		err := rows.Scan(
			&complaint.ID,
			&complaint.UserID,
			&complaint.ComplaintText,
			&complaint.PDFTelegramFileID,
			&complaint.PDFFilename,
			&complaint.CreatedAt,
			&complaint.Status,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan complaint: %w", err)
		}
		complaints = append(complaints, &complaint)
	}

	return complaints, nil
}

// UpdateStatus updates complaint status
func (r *ComplaintRepository) UpdateStatus(id int, status string) error {
	query := `UPDATE complaints SET status = $1 WHERE id = $2`
	_, err := r.db.Exec(query, status, id)
	if err != nil {
		return fmt.Errorf("failed to update complaint status: %w", err)
	}
	return nil
}

// Count counts total complaints
func (r *ComplaintRepository) Count() (int, error) {
	var count int
	err := r.db.QueryRow("SELECT COUNT(*) FROM complaints").Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count complaints: %w", err)
	}
	return count, nil
}

// CountByStatus counts complaints by status
func (r *ComplaintRepository) CountByStatus(status string) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM complaints WHERE status = $1`
	err := r.db.QueryRow(query, status).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count complaints by status: %w", err)
	}
	return count, nil
}

// CountByUserID counts complaints by user ID
func (r *ComplaintRepository) CountByUserID(userID int) (int, error) {
	var count int
	query := `SELECT COUNT(*) FROM complaints WHERE user_id = $1`
	err := r.db.QueryRow(query, userID).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("failed to count user complaints: %w", err)
	}
	return count, nil
}

// CreateComplaintImage adds an image to a complaint
func (r *ComplaintRepository) CreateComplaintImage(complaintID int, image *models.ImageData, orderIndex int) error {
	query := `
		INSERT INTO complaint_images (complaint_id, telegram_file_id, file_unique_id, file_size, mime_type, order_index)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	_, err := r.db.Exec(query, complaintID, image.FileID, image.FileUniqueID, image.FileSize, image.MimeType, orderIndex)
	if err != nil {
		return fmt.Errorf("failed to create complaint image: %w", err)
	}
	return nil
}

// GetComplaintImages gets all images for a complaint
func (r *ComplaintRepository) GetComplaintImages(complaintID int) ([]models.ComplaintImage, error) {
	query := `
		SELECT id, complaint_id, telegram_file_id, file_unique_id, file_size, mime_type, order_index, created_at
		FROM complaint_images
		WHERE complaint_id = $1
		ORDER BY order_index ASC
	`

	rows, err := r.db.Query(query, complaintID)
	if err != nil {
		return nil, fmt.Errorf("failed to get complaint images: %w", err)
	}
	defer rows.Close()

	var images []models.ComplaintImage
	for rows.Next() {
		var img models.ComplaintImage
		err := rows.Scan(
			&img.ID,
			&img.ComplaintID,
			&img.TelegramFileID,
			&img.FileUniqueID,
			&img.FileSize,
			&img.MimeType,
			&img.OrderIndex,
			&img.CreatedAt,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan complaint image: %w", err)
		}
		images = append(images, img)
	}

	return images, nil
}
