-- Migration 003: Add is_document field to announcements
-- This field tracks whether the image was uploaded as a document (HEIC, etc.) or as a photo

-- Add is_document field (default to false for existing records)
-- Check if column exists first to make migration idempotent
-- SQLite doesn't support IF NOT EXISTS for ALTER TABLE, so we check in application layer

-- Drop and recreate the view with the new field
DROP VIEW IF EXISTS v_announcements_with_admin;

CREATE VIEW v_announcements_with_admin AS
SELECT
    a.id,
    a.admin_id,
    a.title,
    a.announcement_text,
    a.image_telegram_file_id,
    a.image_file_unique_id,
    a.image_file_size,
    a.image_mime_type,
    a.is_document,
    a.created_at,
    a.updated_at,
    ad.phone_number AS admin_phone,
    ad.telegram_id AS admin_telegram_id,
    ad.name AS admin_username
FROM announcements a
INNER JOIN admins ad ON a.admin_id = ad.id
ORDER BY a.created_at DESC;
