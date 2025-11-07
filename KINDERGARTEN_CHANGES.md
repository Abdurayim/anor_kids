# Kindergarten Bot - Implementation Summary

## Overview
This bot has been transformed from the parent-bot to a kindergarten complaint system with support for image uploads and PDF generation.

## Key Features

### 1. **Multi-Image Support**
- Users can attach up to 10 images to each complaint
- Images are stored in Telegram cloud (no local storage needed)
- Each image is tracked with metadata (file ID, size, MIME type)

### 2. **PDF Generation**
- Complaints are automatically converted to PDF documents
- PDFs include:
  - User information (child name, class, phone, date)
  - Complaint text
  - All attached images embedded in the document
- PDFs are stored in Telegram cloud for easy access

### 3. **Database Structure**

#### complaints table
- `pdf_telegram_file_id` - File ID of PDF in Telegram cloud
- `pdf_filename` - Name of the PDF file
- `complaint_text` - Text of the complaint
- `user_id` - Reference to user who submitted

#### complaint_images table (NEW)
- `complaint_id` - Reference to complaint
- `telegram_file_id` - Image file ID from Telegram
- `file_unique_id` - Unique identifier
- `order_index` - Order of images in complaint
- `file_size`, `mime_type` - Metadata

### 4. **User Flow**

1. User starts complaint: `/shikoyat` or button press
2. User enters complaint text
3. User uploads images (optional, max 10)
   - Can skip image upload
   - Can finish after uploading desired images
4. User confirms complaint preview
5. Bot generates PDF with text + images
6. PDF is uploaded to Telegram and stored
7. Admin receives notification with PDF

## Technical Changes

### Dependencies Added
```go
github.com/jung-kurt/gofpdf v1.16.2  // PDF generation library
```

### New States
- `StateAwaitingImages` - Collecting images from user

### New Handlers
- `HandleImage()` - Process photo uploads
- `HandleSkipImages()` - Skip image collection
- `HandleFinishImages()` - Complete image collection

### Modified Services

#### DocumentService
- **OLD**: `GenerateComplaintDocument()` - Generated DOCX files
- **NEW**: `GenerateComplaintPDF()` - Generates PDF with images
- `downloadTelegramImage()` - Downloads images from Telegram servers

#### ComplaintService
- `CreateComplaintImage()` - Save image metadata
- `GetComplaintImages()` - Retrieve complaint images

### Font Support
- DejaVu Sans fonts added for Cyrillic character support in PDFs
- Located in `/fonts` directory

## Configuration

### Environment Variables (`.env`)
```bash
# Bot Configuration
BOT_TOKEN=your_telegram_bot_token_here
BOT_MODE=webhook  # or polling

# Webhook Configuration (if using webhook mode)
WEBHOOK_URL=https://your-domain.com/webhook
WEBHOOK_PORT=8080

# Database Configuration
DB_TYPE=sqlite
DB_PATH=./anor_kids_bot.db

# Admin Configuration
ADMIN_PHONES=+998901234567,+998907654321
```

## File Structure
```
anor-kids/
├── fonts/
│   ├── DejaVuSans.ttf
│   └── DejaVuSans-Bold.ttf
├── internal/
│   ├── database/
│   │   └── migrations/
│   │       └── 001_initial_sqlite.sql  # Updated schema
│   ├── handlers/
│   │   ├── complaint.go              # Image handling
│   │   └── router.go                 # New state routing
│   ├── models/
│   │   ├── complaint.go              # Updated models
│   │   └── state.go                  # ImageData type
│   ├── repository/
│   │   └── complaint_repo.go         # Image CRUD
│   ├── services/
│   │   ├── document.go               # PDF generation
│   │   ├── complaint.go              # Image service methods
│   │   └── bot.go                    # Updated initialization
│   └── utils/
│       ├── helpers.go                # GeneratePDFFilename()
│       └── keyboard.go               # Image collection keyboard
├── temp_docs/                        # Temporary PDF storage
├── go.mod                            # Updated dependencies
└── anor-kids                         # Compiled binary
```

## Language Support
- **Uzbek**: Primary language
- **Russian**: Secondary language
- Both languages supported throughout the entire flow

## Storage Strategy
- **PDFs**: Stored in Telegram cloud (via file_id)
- **Images**: Stored in Telegram cloud (via file_id)
- **Database**: Only stores metadata and file IDs
- **Benefit**: Minimal database size, infinite file storage via Telegram

## Admin Notifications
Admins receive:
- Complaint ID
- Child name and class
- Parent phone number
- Submission date
- Number of images attached
- PDF document with full complaint

## Build & Run

### Build
```bash
cd /Users/abdurayim/Desktop/PROJECTS/anor-kids
go build -o anor-kids ./cmd/bot/main.go
```

### Run
```bash
./anor-kids
```

## Database Migration
The bot will automatically create tables on first run. If migrating from parent-bot:
1. Backup existing database
2. Run migration SQL to add `complaint_images` table
3. Update existing `complaints` table schema

## Testing Checklist
- [ ] User can submit complaint with text only
- [ ] User can upload 1-10 images
- [ ] User can skip image upload
- [ ] PDF generation works with Cyrillic text
- [ ] PDF includes all uploaded images
- [ ] Admin receives notification with PDF
- [ ] Database correctly stores complaint and image metadata
- [ ] Telegram cloud storage works for files

## Known Limitations
- Maximum 10 images per complaint (configurable in code)
- PDF size depends on image count and quality
- Requires internet connection to download images from Telegram

## Future Enhancements
- [ ] Image compression before PDF generation
- [ ] Support for video attachments
- [ ] PDF password protection
- [ ] Export complaints to Excel/CSV
- [ ] Analytics dashboard for admins

---
Generated for **Anor Kids Kindergarten Bot** - November 2025
