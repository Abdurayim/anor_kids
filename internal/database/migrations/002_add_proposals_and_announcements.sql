-- Migration 002: Add proposals and announcements functionality
-- This migration adds support for user proposals and admin announcements

-- Proposals table (similar to complaints but for suggestions/proposals)
CREATE TABLE IF NOT EXISTS proposals (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    proposal_text TEXT NOT NULL,
    pdf_telegram_file_id TEXT NOT NULL, -- PDF file stored in Telegram cloud
    pdf_filename TEXT NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    status TEXT DEFAULT 'pending' CHECK (status IN ('pending', 'reviewed', 'archived')),
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Indexes for proposals table
CREATE INDEX IF NOT EXISTS idx_proposals_user_id ON proposals(user_id);
CREATE INDEX IF NOT EXISTS idx_proposals_created_at ON proposals(created_at DESC);
CREATE INDEX IF NOT EXISTS idx_proposals_status ON proposals(status);
CREATE INDEX IF NOT EXISTS idx_proposals_combined ON proposals(user_id, created_at DESC);

-- Proposal images table (supports multiple images per proposal)
CREATE TABLE IF NOT EXISTS proposal_images (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    proposal_id INTEGER NOT NULL,
    telegram_file_id TEXT NOT NULL, -- Image file ID from Telegram
    file_unique_id TEXT NOT NULL, -- Unique file identifier from Telegram
    file_size INTEGER,
    mime_type TEXT,
    order_index INTEGER DEFAULT 0, -- Order of images in the proposal
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (proposal_id) REFERENCES proposals(id) ON DELETE CASCADE
);

-- Indexes for proposal_images table
CREATE INDEX IF NOT EXISTS idx_proposal_images_proposal_id ON proposal_images(proposal_id);
CREATE INDEX IF NOT EXISTS idx_proposal_images_order ON proposal_images(proposal_id, order_index);

-- Announcements table (admin posts with image and text)
CREATE TABLE IF NOT EXISTS announcements (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    admin_id INTEGER NOT NULL,
    title TEXT NOT NULL,
    announcement_text TEXT NOT NULL,
    image_telegram_file_id TEXT NOT NULL, -- Image file ID from Telegram
    image_file_unique_id TEXT, -- Unique file identifier from Telegram
    image_file_size INTEGER,
    image_mime_type TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (admin_id) REFERENCES admins(id) ON DELETE CASCADE
);

-- Indexes for announcements table
CREATE INDEX IF NOT EXISTS idx_announcements_admin_id ON announcements(admin_id);
CREATE INDEX IF NOT EXISTS idx_announcements_created_at ON announcements(created_at DESC);

-- Trigger to automatically update updated_at timestamp for announcements
CREATE TRIGGER IF NOT EXISTS update_announcements_updated_at
AFTER UPDATE ON announcements
FOR EACH ROW
WHEN OLD.updated_at = NEW.updated_at
BEGIN
    UPDATE announcements SET updated_at = CURRENT_TIMESTAMP WHERE id = NEW.id;
END;

-- Create view for proposals with user info
CREATE VIEW IF NOT EXISTS v_proposals_with_user AS
SELECT
    p.id,
    p.user_id,
    p.proposal_text,
    p.pdf_telegram_file_id,
    p.pdf_filename,
    p.created_at,
    p.status,
    u.telegram_id AS user_telegram_id,
    u.telegram_username,
    u.phone_number,
    u.child_name,
    u.child_class
FROM proposals p
INNER JOIN users u ON p.user_id = u.id
ORDER BY p.created_at DESC;

-- Create view for announcements with admin info
CREATE VIEW IF NOT EXISTS v_announcements_with_admin AS
SELECT
    a.id,
    a.admin_id,
    a.title,
    a.announcement_text,
    a.image_telegram_file_id,
    a.image_file_unique_id,
    a.image_file_size,
    a.image_mime_type,
    a.created_at,
    a.updated_at,
    ad.phone_number AS admin_phone,
    ad.telegram_id AS admin_telegram_id,
    ad.name AS admin_username
FROM announcements a
INNER JOIN admins ad ON a.admin_id = ad.id
ORDER BY a.created_at DESC;
