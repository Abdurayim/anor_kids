package services

import (
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"anor-kids/internal/models"
	"anor-kids/internal/utils"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"github.com/jung-kurt/gofpdf"
)

// DocumentService handles document generation and management
type DocumentService struct {
	tempDir string
	bot     *tgbotapi.BotAPI
}

// NewDocumentService creates a new document service
func NewDocumentService(tempDir string, bot *tgbotapi.BotAPI) *DocumentService {
	return &DocumentService{
		tempDir: tempDir,
		bot:     bot,
	}
}

// GenerateComplaintPDF generates a PDF document for a complaint with text and images
// Returns the file path and filename
func (s *DocumentService) GenerateComplaintPDF(user *models.User, complaintText string, images []models.ImageData) (filePath, filename string, err error) {
	// Generate filename with .pdf extension
	filename = utils.GeneratePDFFilename(user.ChildName, user.ChildClass)

	// Create full path
	filePath = filepath.Join(s.tempDir, filename)

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add UTF-8 font support for Cyrillic characters
	pdf.AddUTF8Font("DejaVu", "", filepath.Join("fonts", "DejaVuSans.ttf"))
	pdf.AddUTF8Font("DejaVu", "B", filepath.Join("fonts", "DejaVuSans-Bold.ttf"))

	// Set font
	pdf.SetFont("DejaVu", "B", 16)

	// Title
	pdf.CellFormat(0, 10, "SHIKOYAT / ЖАЛОБА", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// User information
	pdf.SetFont("DejaVu", "B", 12)
	pdf.Cell(0, 8, "Ma'lumotlar / Информация:")
	pdf.Ln(8)

	pdf.SetFont("DejaVu", "", 11)
	// Strip emojis from user data for PDF compatibility
	pdf.Cell(0, 6, fmt.Sprintf("Farzand / Ребенок: %s", utils.StripEmojis(user.ChildName)))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Guruh / Группа: %s", utils.StripEmojis(user.ChildClass)))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Telefon / Телефон: %s", user.PhoneNumber))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Sana / Дата: %s", time.Now().Format("02.01.2006 15:04")))
	pdf.Ln(10)

	// Complaint text
	pdf.SetFont("DejaVu", "B", 12)
	pdf.Cell(0, 8, "Shikoyat matni / Текст жалобы:")
	pdf.Ln(8)

	pdf.SetFont("DejaVu", "", 11)
	// Strip emojis from complaint text for PDF compatibility
	pdf.MultiCell(0, 6, utils.StripEmojis(complaintText), "", "", false)
	pdf.Ln(5)

	// Add images if any
	if len(images) > 0 {
		pdf.Ln(5)
		pdf.SetFont("DejaVu", "B", 12)
		pdf.Cell(0, 8, fmt.Sprintf("Qo'shimcha rasmlar / Приложенные изображения: %d ta", len(images)))
		pdf.Ln(8)

		for i, img := range images {
			// Download image from Telegram
			imgPath, err := s.downloadTelegramImage(img.FileID, i)
			if err != nil {
				fmt.Printf("[WARN] Failed to download image %d: %v\n", i, err)
				continue
			}
			defer os.Remove(imgPath) // Clean up after adding to PDF

			// Add image number
			pdf.SetFont("DejaVu", "", 10)
			pdf.Cell(0, 5, fmt.Sprintf("Rasm %d / Изображение %d:", i+1, i+1))
			pdf.Ln(6)

			// Get current Y position
			currentY := pdf.GetY()

			// Calculate max image dimensions for A4 page
			// A4 width = 210mm, with 15mm margins on each side = 180mm usable
			maxWidth := 180.0
			maxHeight := 200.0 // Leave space for margins and other content

			// Check if we need a new page (if less than 60mm space remaining)
			if currentY > 240 {
				pdf.AddPage()
				currentY = pdf.GetY()
			}

			// Add image to PDF with proper sizing
			// Using ImageOptions with flow mode to automatically calculate height
			options := gofpdf.ImageOptions{
				ImageType: "",
				ReadDpi:   false,
			}

			// Get image info to calculate aspect ratio
			info := pdf.RegisterImageOptions(imgPath, options)
			if info != nil {
				imgWidth := float64(info.Width())
				imgHeight := float64(info.Height())

				// Calculate scaled dimensions maintaining aspect ratio
				scaledWidth := maxWidth
				scaledHeight := (imgHeight / imgWidth) * maxWidth

				// If height is too large, scale based on height instead
				if scaledHeight > maxHeight {
					scaledHeight = maxHeight
					scaledWidth = (imgWidth / imgHeight) * maxHeight
				}

				// Center the image horizontally
				xPos := (210.0 - scaledWidth) / 2.0

				// Add the image with calculated dimensions
				pdf.Image(imgPath, xPos, currentY, scaledWidth, scaledHeight, false, "", 0, "")

				// Move Y position down by image height plus spacing
				pdf.SetY(currentY + scaledHeight + 10)
			} else {
				// Fallback if image info not available
				pdf.Image(imgPath, 15, currentY, maxWidth, 0, false, "", 0, "")
				pdf.Ln(80) // Approximate spacing
			}
		}
	}

	// Save PDF
	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Verify file was created and has content
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to verify generated file: %w", err)
	}

	if fileInfo.Size() == 0 {
		return "", "", fmt.Errorf("generated file is empty")
	}

	fmt.Printf("[DEBUG] PDF verified: %s, size: %d bytes\n", filename, fileInfo.Size())

	return filePath, filename, nil
}

// downloadTelegramImage downloads an image from Telegram servers
func (s *DocumentService) downloadTelegramImage(fileID string, index int) (string, error) {
	// Get file from Telegram
	file, err := s.bot.GetFile(tgbotapi.FileConfig{FileID: fileID})
	if err != nil {
		return "", fmt.Errorf("failed to get file: %w", err)
	}

	// Download file
	fileURL := file.Link(s.bot.Token)
	resp, err := http.Get(fileURL)
	if err != nil {
		return "", fmt.Errorf("failed to download file: %w", err)
	}
	defer resp.Body.Close()

	// Create temp file
	tempFile := filepath.Join(s.tempDir, fmt.Sprintf("temp_img_%d_%d.jpg", time.Now().Unix(), index))
	out, err := os.Create(tempFile)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	// Save file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return tempFile, nil
}

// DeleteTempFile deletes a temporary file
func (s *DocumentService) DeleteTempFile(filePath string) error {
	if err := os.Remove(filePath); err != nil {
		return fmt.Errorf("failed to delete temp file: %w", err)
	}

	return nil
}

// CleanTempDirectory cleans old temporary files
func (s *DocumentService) CleanTempDirectory(maxAge time.Duration) error {
	files, err := os.ReadDir(s.tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %w", err)
	}

	now := time.Now()
	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		// Delete files older than maxAge
		if now.Sub(info.ModTime()) > maxAge {
			filePath := filepath.Join(s.tempDir, file.Name())
			_ = os.Remove(filePath) // Ignore errors
		}
	}

	return nil
}

// GetTempDir returns the temporary directory path
func (s *DocumentService) GetTempDir() string {
	return s.tempDir
}

// GenerateProposalPDF generates a PDF document for a proposal with text and images
// Returns the file path and filename
func (s *DocumentService) GenerateProposalPDF(user *models.User, proposalText string, images []models.ImageData) (filePath, filename string, err error) {
	// Generate filename with .pdf extension
	filename = utils.GenerateProposalPDFFilename(user.ChildName, user.ChildClass)

	// Create full path
	filePath = filepath.Join(s.tempDir, filename)

	// Create PDF
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	// Add UTF-8 font support for Cyrillic characters
	pdf.AddUTF8Font("DejaVu", "", filepath.Join("fonts", "DejaVuSans.ttf"))
	pdf.AddUTF8Font("DejaVu", "B", filepath.Join("fonts", "DejaVuSans-Bold.ttf"))

	// Set font
	pdf.SetFont("DejaVu", "B", 16)

	// Title
	pdf.CellFormat(0, 10, "TAKLIF / ПРЕДЛОЖЕНИЕ", "", 1, "C", false, 0, "")
	pdf.Ln(5)

	// User information
	pdf.SetFont("DejaVu", "B", 12)
	pdf.Cell(0, 8, "Ma'lumotlar / Информация:")
	pdf.Ln(8)

	pdf.SetFont("DejaVu", "", 11)
	// Strip emojis from user data for PDF compatibility
	pdf.Cell(0, 6, fmt.Sprintf("Farzand / Ребенок: %s", utils.StripEmojis(user.ChildName)))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Guruh / Группа: %s", utils.StripEmojis(user.ChildClass)))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Telefon / Телефон: %s", user.PhoneNumber))
	pdf.Ln(6)
	pdf.Cell(0, 6, fmt.Sprintf("Sana / Дата: %s", time.Now().Format("02.01.2006 15:04")))
	pdf.Ln(10)

	// Proposal text
	pdf.SetFont("DejaVu", "B", 12)
	pdf.Cell(0, 8, "Taklif matni / Текст предложения:")
	pdf.Ln(8)

	pdf.SetFont("DejaVu", "", 11)
	// Strip emojis from proposal text for PDF compatibility
	pdf.MultiCell(0, 6, utils.StripEmojis(proposalText), "", "", false)
	pdf.Ln(5)

	// Add images if any
	if len(images) > 0 {
		pdf.Ln(5)
		pdf.SetFont("DejaVu", "B", 12)
		pdf.Cell(0, 8, fmt.Sprintf("Qo'shimcha rasmlar / Приложенные изображения: %d ta", len(images)))
		pdf.Ln(8)

		for i, img := range images {
			// Download image from Telegram
			imgPath, err := s.downloadTelegramImage(img.FileID, i)
			if err != nil {
				fmt.Printf("[WARN] Failed to download image %d: %v\n", i, err)
				continue
			}
			defer os.Remove(imgPath) // Clean up after adding to PDF

			// Add image number
			pdf.SetFont("DejaVu", "", 10)
			pdf.Cell(0, 5, fmt.Sprintf("Rasm %d / Изображение %d:", i+1, i+1))
			pdf.Ln(6)

			// Get current Y position
			currentY := pdf.GetY()

			// Calculate max image dimensions for A4 page
			// A4 width = 210mm, with 15mm margins on each side = 180mm usable
			maxWidth := 180.0
			maxHeight := 200.0 // Leave space for margins and other content

			// Check if we need a new page (if less than 60mm space remaining)
			if currentY > 240 {
				pdf.AddPage()
				currentY = pdf.GetY()
			}

			// Add image to PDF with proper sizing
			// Using ImageOptions with flow mode to automatically calculate height
			options := gofpdf.ImageOptions{
				ImageType: "",
				ReadDpi:   false,
			}

			// Get image info to calculate aspect ratio
			info := pdf.RegisterImageOptions(imgPath, options)
			if info != nil {
				imgWidth := float64(info.Width())
				imgHeight := float64(info.Height())

				// Calculate scaled dimensions maintaining aspect ratio
				scaledWidth := maxWidth
				scaledHeight := (imgHeight / imgWidth) * maxWidth

				// If height is too large, scale based on height instead
				if scaledHeight > maxHeight {
					scaledHeight = maxHeight
					scaledWidth = (imgWidth / imgHeight) * maxHeight
				}

				// Center the image horizontally
				xPos := (210.0 - scaledWidth) / 2.0

				// Add the image with calculated dimensions
				pdf.Image(imgPath, xPos, currentY, scaledWidth, scaledHeight, false, "", 0, "")

				// Move Y position down by image height plus spacing
				pdf.SetY(currentY + scaledHeight + 10)
			} else {
				// Fallback if image info not available
				pdf.Image(imgPath, 15, currentY, maxWidth, 0, false, "", 0, "")
				pdf.Ln(80) // Approximate spacing
			}
		}
	}

	// Save PDF
	err = pdf.OutputFileAndClose(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to generate PDF: %w", err)
	}

	// Verify file was created and has content
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		return "", "", fmt.Errorf("failed to verify generated file: %w", err)
	}

	if fileInfo.Size() == 0 {
		return "", "", fmt.Errorf("generated file is empty")
	}

	fmt.Printf("[DEBUG] PDF verified: %s, size: %d bytes\n", filename, fileInfo.Size())

	return filePath, filename, nil
}
