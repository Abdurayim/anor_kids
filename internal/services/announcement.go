package services

import (
	"fmt"

	"anor-kids/internal/models"
	"anor-kids/internal/repository"
)

// AnnouncementService handles announcement-related business logic
type AnnouncementService struct {
	repo      *repository.AnnouncementRepository
	adminRepo *repository.AdminRepository
}

// NewAnnouncementService creates a new announcement service
func NewAnnouncementService(repo *repository.AnnouncementRepository, adminRepo *repository.AdminRepository) *AnnouncementService {
	return &AnnouncementService{
		repo:      repo,
		adminRepo: adminRepo,
	}
}

// CreateAnnouncement creates a new announcement
func (s *AnnouncementService) CreateAnnouncement(req *models.CreateAnnouncementRequest) (*models.Announcement, error) {
	// Verify admin exists
	admin, err := s.adminRepo.GetByID(req.AdminID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify admin: %w", err)
	}
	if admin == nil {
		return nil, fmt.Errorf("admin not found")
	}

	announcement, err := s.repo.Create(req)
	if err != nil {
		return nil, fmt.Errorf("failed to create announcement: %w", err)
	}

	return announcement, nil
}

// GetAnnouncementByID gets announcement by ID
func (s *AnnouncementService) GetAnnouncementByID(id int) (*models.Announcement, error) {
	announcement, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get announcement: %w", err)
	}

	return announcement, nil
}

// GetAllAnnouncements gets all announcements with pagination (for users)
func (s *AnnouncementService) GetAllAnnouncements(limit, offset int) ([]*models.Announcement, error) {
	announcements, err := s.repo.GetAll(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get announcements: %w", err)
	}

	return announcements, nil
}

// GetAdminAnnouncements gets announcements by admin ID (for admin to see their own posts)
func (s *AnnouncementService) GetAdminAnnouncements(adminID int, limit, offset int) ([]*models.Announcement, error) {
	announcements, err := s.repo.GetByAdminID(adminID, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get admin announcements: %w", err)
	}

	return announcements, nil
}

// GetAllAnnouncementsWithAdmin gets all announcements with admin info
func (s *AnnouncementService) GetAllAnnouncementsWithAdmin(limit, offset int) ([]*models.AnnouncementWithAdmin, error) {
	announcements, err := s.repo.GetAllWithAdmin(limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to get announcements with admin: %w", err)
	}

	return announcements, nil
}

// UpdateAnnouncement updates an announcement
func (s *AnnouncementService) UpdateAnnouncement(req *models.UpdateAnnouncementRequest) error {
	// Verify announcement exists
	announcement, err := s.repo.GetByID(req.ID)
	if err != nil {
		return fmt.Errorf("failed to verify announcement: %w", err)
	}
	if announcement == nil {
		return fmt.Errorf("announcement not found")
	}

	err = s.repo.Update(req)
	if err != nil {
		return fmt.Errorf("failed to update announcement: %w", err)
	}

	return nil
}

// DeleteAnnouncement deletes an announcement
func (s *AnnouncementService) DeleteAnnouncement(id int, adminID int) error {
	// Verify announcement exists and belongs to admin
	announcement, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to verify announcement: %w", err)
	}
	if announcement == nil {
		return fmt.Errorf("announcement not found")
	}

	// Check if admin owns this announcement
	if announcement.AdminID != adminID {
		return fmt.Errorf("unauthorized: announcement does not belong to this admin")
	}

	err = s.repo.Delete(id)
	if err != nil {
		return fmt.Errorf("failed to delete announcement: %w", err)
	}

	return nil
}

// CountAnnouncements counts total announcements
func (s *AnnouncementService) CountAnnouncements() (int, error) {
	count, err := s.repo.Count()
	if err != nil {
		return 0, fmt.Errorf("failed to count announcements: %w", err)
	}

	return count, nil
}

// CountAdminAnnouncements counts announcements by admin ID
func (s *AnnouncementService) CountAdminAnnouncements(adminID int) (int, error) {
	count, err := s.repo.CountByAdminID(adminID)
	if err != nil {
		return 0, fmt.Errorf("failed to count admin announcements: %w", err)
	}

	return count, nil
}
