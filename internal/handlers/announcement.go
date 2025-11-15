package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"anor-kids/internal/i18n"
	"anor-kids/internal/models"
	"anor-kids/internal/services"
	"anor-kids/internal/utils"
)

const announcementsPerPage = 5

// sendAnnouncementMessage sends an announcement with proper handling for photos vs documents
func sendAnnouncementMessage(bot *tgbotapi.BotAPI, chatID int64, announcement *models.Announcement, caption string, keyboard *tgbotapi.InlineKeyboardMarkup) error {
	// Check if it's a document or photo
	if announcement.IsDocument {
		// Send as document (for HEIC, SVG, etc.)
		doc := tgbotapi.NewDocument(chatID, tgbotapi.FileID(announcement.ImageTelegramFileID))
		doc.Caption = caption
		doc.ParseMode = "HTML"
		if keyboard != nil {
			doc.ReplyMarkup = keyboard
		}
		_, err := bot.Send(doc)
		return err
	} else {
		// Send as photo (for JPG, PNG, etc.)
		photo := tgbotapi.NewPhoto(chatID, tgbotapi.FileID(announcement.ImageTelegramFileID))
		photo.Caption = caption
		photo.ParseMode = "HTML"
		if keyboard != nil {
			photo.ReplyMarkup = keyboard
		}
		_, err := bot.Send(photo)
		return err
	}
}

// HandleViewAnnouncements handles when users view announcements
func HandleViewAnnouncements(botService *services.BotService, message *tgbotapi.Message) error {
	user, err := botService.UserService.GetUserByTelegramID(message.From.ID)
	if err != nil || user == nil {
		return HandleStart(botService, message)
	}

	lang := i18n.GetLanguage(user.Language)
	chatID := message.Chat.ID

	// Get all announcements (ordered by newest first)
	announcements, err := botService.AnnouncementService.GetAllAnnouncements(100, 0)
	if err != nil {
		text := i18n.Get(i18n.ErrDatabaseError, lang)
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	if len(announcements) == 0 {
		text := i18n.Get(i18n.MsgNoAnnouncements, lang)
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	// Send each announcement with proper media handling (HTML-escaped to prevent injection)
	for _, announcement := range announcements {
		// Escape HTML special characters to prevent injection
		escapedTitle := utils.EscapeHTML(announcement.Title)
		escapedText := utils.EscapeHTML(announcement.AnnouncementText)
		caption := fmt.Sprintf("📢 <b>%s</b>\n\n%s", escapedTitle, escapedText)

		err := sendAnnouncementMessage(botService.Bot, chatID, announcement, caption, nil)
		if err != nil {
			// If failed to send media, try sending text only
			_ = botService.TelegramService.SendMessage(chatID, caption, nil)
		}
	}

	return nil
}

// HandleAdminCreateAnnouncementCallback handles create announcement button click
func HandleAdminCreateAnnouncementCallback(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
	// Verify admin
	admin, err := botService.AdminRepo.GetByTelegramID(callback.From.ID)
	if err != nil || admin == nil {
		return botService.TelegramService.AnswerCallbackQuery(callback.ID, "Unauthorized")
	}

	// Get admin's language preference (check if admin is also a user)
	lang := i18n.LanguageUzbek // Default
	user, err := botService.UserService.GetUserByTelegramID(callback.From.ID)
	if err == nil && user != nil {
		lang = i18n.GetLanguage(user.Language)
	}

	chatID := callback.Message.Chat.ID

	// Set state to awaiting announcement title
	stateData := &models.StateData{
		Language: string(lang),
	}
	err = botService.StateManager.Set(callback.From.ID, models.StateAwaitingAnnouncementTitle, stateData)
	if err != nil {
		return err
	}

	// Ask for announcement title
	text := i18n.Get(i18n.MsgRequestAnnouncementTitle, lang)
	return botService.TelegramService.SendMessage(chatID, text, nil)
}

// HandleAnnouncementTitle handles announcement title input
func HandleAnnouncementTitle(botService *services.BotService, message *tgbotapi.Message, stateData *models.StateData) error {
	lang := i18n.GetLanguage(stateData.Language)
	chatID := message.Chat.ID
	title := strings.TrimSpace(message.Text)

	// Validate title
	if len(title) < 3 || len(title) > 200 {
		text := "❌ Sarlavha 3 dan 200 gacha belgi bo'lishi kerak.\n\n❌ Заголовок должен содержать от 3 до 200 символов."
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	// Save title to state
	stateData.AnnouncementTitle = title

	// Update state to awaiting announcement text
	err := botService.StateManager.Set(message.From.ID, models.StateAwaitingAnnouncementText, stateData)
	if err != nil {
		return err
	}

	// Ask for announcement text
	text := i18n.Get(i18n.MsgRequestAnnouncementText, lang)
	return botService.TelegramService.SendMessage(chatID, text, nil)
}

// HandleAnnouncementText handles announcement text input
func HandleAnnouncementText(botService *services.BotService, message *tgbotapi.Message, stateData *models.StateData) error {
	lang := i18n.GetLanguage(stateData.Language)
	chatID := message.Chat.ID
	text := strings.TrimSpace(message.Text)

	// Validate text
	if len(text) < 10 || len(text) > 5000 {
		errorText := "❌ E'lon matni 10 dan 5000 gacha belgi bo'lishi kerak.\n\n❌ Текст объявления должен содержать от 10 до 5000 символов."
		return botService.TelegramService.SendMessage(chatID, errorText, nil)
	}

	// Save text to state
	stateData.AnnouncementText = text

	// Update state to awaiting announcement image
	err := botService.StateManager.Set(message.From.ID, models.StateAwaitingAnnouncementImage, stateData)
	if err != nil {
		return err
	}

	// Ask for announcement image
	requestText := i18n.Get(i18n.MsgRequestAnnouncementImage, lang)
	return botService.TelegramService.SendMessage(chatID, requestText, nil)
}

// HandleAnnouncementImage handles announcement image upload
func HandleAnnouncementImage(botService *services.BotService, message *tgbotapi.Message, stateData *models.StateData) error {
	lang := i18n.GetLanguage(stateData.Language)
	chatID := message.Chat.ID

	// Get the largest photo size
	var photo *tgbotapi.PhotoSize
	if message.Photo != nil && len(message.Photo) > 0 {
		photo = &message.Photo[len(message.Photo)-1]
	} else if message.Document != nil {
		// Handle document (HEIC, etc.)
		// Telegram converts HEIC to JPG automatically when sent as photo
		// If sent as document, we'll store it as is
		doc := message.Document
		stateData.AnnouncementImage = &models.ImageData{
			FileID:       doc.FileID,
			FileUniqueID: doc.FileUniqueID,
			FileSize:     doc.FileSize,
			MimeType:     doc.MimeType,
			IsDocument:   true, // Mark as document
		}
	}

	if photo != nil {
		stateData.AnnouncementImage = &models.ImageData{
			FileID:       photo.FileID,
			FileUniqueID: photo.FileUniqueID,
			FileSize:     photo.FileSize,
			IsDocument:   false, // Mark as photo
		}
	}

	if stateData.AnnouncementImage == nil {
		errorText := "❌ Iltimos, rasm yuboring.\n\n❌ Пожалуйста, отправьте изображение."
		return botService.TelegramService.SendMessage(chatID, errorText, nil)
	}

	// Get admin info
	admin, err := botService.AdminRepo.GetByTelegramID(message.From.ID)
	if err != nil || admin == nil {
		errorText := i18n.Get(i18n.ErrDatabaseError, lang)
		return botService.TelegramService.SendMessage(chatID, errorText, nil)
	}

	// Create announcement
	req := &models.CreateAnnouncementRequest{
		AdminID:             admin.ID,
		Title:               stateData.AnnouncementTitle,
		AnnouncementText:    stateData.AnnouncementText,
		ImageTelegramFileID: stateData.AnnouncementImage.FileID,
		ImageFileUniqueID:   stateData.AnnouncementImage.FileUniqueID,
		ImageFileSize:       stateData.AnnouncementImage.FileSize,
		ImageMimeType:       stateData.AnnouncementImage.MimeType,
		IsDocument:          stateData.AnnouncementImage.IsDocument,
	}

	announcement, err := botService.AnnouncementService.CreateAnnouncement(req)
	if err != nil {
		errorText := i18n.Get(i18n.ErrDatabaseError, lang)
		return botService.TelegramService.SendMessage(chatID, errorText, nil)
	}

	// Clear state (delete it instead of setting to registered, as admin might not be a user)
	err = botService.StateManager.Delete(message.From.ID)
	if err != nil {
		// Log error but continue - not critical
		fmt.Printf("Warning: failed to delete state for admin %d: %v\n", message.From.ID, err)
	}

	// Send preview of announcement to admin (what users will see)
	escapedTitle := utils.EscapeHTML(announcement.Title)
	escapedText := utils.EscapeHTML(announcement.AnnouncementText)
	caption := fmt.Sprintf("📢 <b>%s</b>\n\n%s", escapedTitle, escapedText)

	err = sendAnnouncementMessage(botService.Bot, chatID, announcement, caption, nil)
	if err != nil {
		// If failed to send media, try sending text only
		_ = botService.TelegramService.SendMessage(chatID, caption, nil)
	}

	// Notify admin of success
	successText := i18n.Get(i18n.MsgAnnouncementCreated, lang)
	_ = botService.TelegramService.SendMessage(chatID, successText, nil)

	// Send announcement to all users in background
	go BroadcastAnnouncement(botService, announcement, chatID)

	return nil
}

// BroadcastAnnouncement sends announcement to all registered users
// Runs in background goroutine with rate limiting
func BroadcastAnnouncement(botService *services.BotService, announcement *models.Announcement, adminChatID int64) {
	// Get all users - for large user bases, implement pagination
	const batchSize = 10000
	users, err := botService.UserRepo.GetAll(batchSize, 0)
	if err != nil {
		// Notify admin of broadcast failure
		errorMsg := fmt.Sprintf("⚠️ E'lonni yuborishda xatolik: %v\n\n⚠️ Ошибка при рассылке объявления: %v", err, err)
		_ = botService.TelegramService.SendMessage(adminChatID, errorMsg, nil)
		return
	}

	if len(users) == 0 {
		return
	}

	// Escape HTML to prevent injection
	escapedTitle := utils.EscapeHTML(announcement.Title)
	escapedText := utils.EscapeHTML(announcement.AnnouncementText)
	caption := fmt.Sprintf("📢 <b>%s</b>\n\n%s", escapedTitle, escapedText)

	successCount := 0
	failCount := 0
	totalUsers := len(users)

	// Rate limiting: Telegram allows ~30 messages per second
	// We'll send 20 per second to be safe
	ticker := time.NewTicker(50 * time.Millisecond) // 20 messages per second
	defer ticker.Stop()

	for _, user := range users {
		<-ticker.C // Wait for rate limit

		err := sendAnnouncementMessage(botService.Bot, user.TelegramID, announcement, caption, nil)
		if err != nil {
			failCount++
		} else {
			successCount++
		}
	}

	// Notify admin of broadcast completion
	summaryMsg := fmt.Sprintf(
		"✅ E'lon yuborish yakunlandi!\n"+
			"Jami foydalanuvchilar: %d\n"+
			"Muvaffaqiyatli: %d\n"+
			"Xato: %d\n\n"+
			"✅ Рассылка завершена!\n"+
			"Всего пользователей: %d\n"+
			"Успешно: %d\n"+
			"Ошибок: %d",
		totalUsers, successCount, failCount,
		totalUsers, successCount, failCount,
	)
	_ = botService.TelegramService.SendMessage(adminChatID, summaryMsg, nil)
}

// HandleAdminManageAnnouncementsCallback handles manage announcements button click
func HandleAdminManageAnnouncementsCallback(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
	// Verify admin
	admin, err := botService.AdminRepo.GetByTelegramID(callback.From.ID)
	if err != nil || admin == nil {
		return botService.TelegramService.AnswerCallbackQuery(callback.ID, "Unauthorized")
	}

	lang := i18n.LanguageUzbek // Default
	chatID := callback.Message.Chat.ID

	// Get admin's announcements
	announcements, err := botService.AnnouncementService.GetAdminAnnouncements(admin.ID, announcementsPerPage, 0)
	if err != nil {
		text := i18n.Get(i18n.ErrDatabaseError, lang)
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	if len(announcements) == 0 {
		text := i18n.Get(i18n.MsgNoAnnouncements, lang)
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	// Send each announcement with delete button
	for _, announcement := range announcements {
		// Escape HTML to prevent injection
		escapedTitle := utils.EscapeHTML(announcement.Title)
		escapedText := utils.EscapeHTML(announcement.AnnouncementText)

		// Telegram caption limit is 1024 characters
		// Create caption with title and date, check if full text fits
		dateStr := announcement.CreatedAt.Format("02.01.2006 15:04")
		caption := fmt.Sprintf("📢 <b>%s</b>\n\n%s\n\n📅 %s",
			escapedTitle,
			escapedText,
			dateStr,
		)

		keyboard := utils.MakeAnnouncementDeleteKeyboard(announcement.ID, lang)

		// If caption is too long, truncate and send full text separately
		if len(caption) > 1024 {
			// Send media with truncated caption
			shortCaption := fmt.Sprintf("📢 <b>%s</b>\n\n%s...\n\n📅 %s",
				escapedTitle,
				utils.TruncateText(escapedText, 800),
				dateStr,
			)

			err := sendAnnouncementMessage(botService.Bot, chatID, announcement, shortCaption, nil)
			if err != nil {
				// If failed to send media, send text only
				_ = botService.TelegramService.SendMessage(chatID, shortCaption, nil)
			}

			// Send full text as separate message with keyboard
			fullText := fmt.Sprintf("<b>Kengaytirilgan matn / Полный текст:</b>\n\n%s", escapedText)
			msg := tgbotapi.NewMessage(chatID, fullText)
			msg.ParseMode = "HTML"
			msg.ReplyMarkup = keyboard
			_, _ = botService.Bot.Send(msg)
		} else {
			// Send media with full caption
			err := sendAnnouncementMessage(botService.Bot, chatID, announcement, caption, &keyboard)
			if err != nil {
				// If failed to send media, try sending text only
				_ = botService.TelegramService.SendMessage(chatID, caption, &keyboard)
			}
		}
	}

	// Get total count for pagination info
	totalCount, _ := botService.AnnouncementService.CountAdminAnnouncements(admin.ID)
	totalPages := (totalCount + announcementsPerPage - 1) / announcementsPerPage

	if totalPages > 1 {
		paginationText := fmt.Sprintf("📄 Sahifa 1/%d", totalPages)
		keyboard := utils.MakeAnnouncementListKeyboard(0, totalPages, lang)
		_ = botService.TelegramService.SendMessage(chatID, paginationText, &keyboard)
	}

	return botService.TelegramService.AnswerCallbackQuery(callback.ID, "")
}

// HandleDeleteAnnouncementCallback handles delete announcement button click
func HandleDeleteAnnouncementCallback(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
	// Verify admin
	admin, err := botService.AdminRepo.GetByTelegramID(callback.From.ID)
	if err != nil || admin == nil {
		return botService.TelegramService.AnswerCallbackQuery(callback.ID, "Unauthorized")
	}

	// Extract announcement ID from callback data
	// Format: delete_announcement_123
	parts := strings.Split(callback.Data, "_")
	if len(parts) != 3 {
		return botService.TelegramService.AnswerCallbackQuery(callback.ID, "Invalid data")
	}

	announcementID, err := strconv.Atoi(parts[2])
	if err != nil {
		return botService.TelegramService.AnswerCallbackQuery(callback.ID, "Invalid ID")
	}

	// Get admin's language preference
	lang := i18n.LanguageUzbek
	user, err := botService.UserService.GetUserByTelegramID(callback.From.ID)
	if err == nil && user != nil {
		lang = i18n.GetLanguage(user.Language)
	}

	// Delete announcement
	err = botService.AnnouncementService.DeleteAnnouncement(announcementID, admin.ID)
	if err != nil {
		errorText := "❌ E'lonni o'chirishda xatolik\n\n❌ Ошибка при удалении объявления"
		if err.Error() == "unauthorized: announcement does not belong to this admin" {
			errorText = "❌ Bu e'lon sizga tegishli emas\n\n❌ Это объявление вам не принадлежит"
		}
		return botService.TelegramService.AnswerCallbackQuery(callback.ID, errorText)
	}

	// Delete the message
	deleteMsg := tgbotapi.NewDeleteMessage(callback.Message.Chat.ID, callback.Message.MessageID)
	_, _ = botService.Bot.Request(deleteMsg)

	// Notify admin
	successText := i18n.Get(i18n.MsgAnnouncementDeleted, lang)
	_ = botService.TelegramService.AnswerCallbackQuery(callback.ID, successText)

	return nil
}

// HandleAnnouncementPageCallback handles announcement pagination
func HandleAnnouncementPageCallback(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
	// Verify admin
	admin, err := botService.AdminRepo.GetByTelegramID(callback.From.ID)
	if err != nil || admin == nil {
		return botService.TelegramService.AnswerCallbackQuery(callback.ID, "Unauthorized")
	}

	// Extract page number from callback data
	// Format: announcements_page_1
	parts := strings.Split(callback.Data, "_")
	if len(parts) != 3 {
		return botService.TelegramService.AnswerCallbackQuery(callback.ID, "Invalid data")
	}

	page, err := strconv.Atoi(parts[2])
	if err != nil || page < 0 {
		return botService.TelegramService.AnswerCallbackQuery(callback.ID, "Invalid page")
	}

	lang := i18n.LanguageUzbek
	chatID := callback.Message.Chat.ID

	// Get announcements for this page
	offset := page * announcementsPerPage
	announcements, err := botService.AnnouncementService.GetAdminAnnouncements(admin.ID, announcementsPerPage, offset)
	if err != nil {
		text := i18n.Get(i18n.ErrDatabaseError, lang)
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	// Send announcements
	for _, announcement := range announcements {
		// Escape HTML to prevent injection
		escapedTitle := utils.EscapeHTML(announcement.Title)
		escapedText := utils.EscapeHTML(announcement.AnnouncementText)

		// Telegram caption limit is 1024 characters
		dateStr := announcement.CreatedAt.Format("02.01.2006 15:04")
		caption := fmt.Sprintf("📢 <b>%s</b>\n\n%s\n\n📅 %s",
			escapedTitle,
			escapedText,
			dateStr,
		)

		keyboard := utils.MakeAnnouncementDeleteKeyboard(announcement.ID, lang)

		// If caption is too long, truncate and send full text separately
		if len(caption) > 1024 {
			// Send media with truncated caption
			shortCaption := fmt.Sprintf("📢 <b>%s</b>\n\n%s...\n\n📅 %s",
				escapedTitle,
				utils.TruncateText(escapedText, 800),
				dateStr,
			)

			_ = sendAnnouncementMessage(botService.Bot, chatID, announcement, shortCaption, nil)

			// Send full text as separate message with keyboard
			fullText := fmt.Sprintf("<b>Kengaytirilgan matn / Полный текст:</b>\n\n%s", escapedText)
			msg := tgbotapi.NewMessage(chatID, fullText)
			msg.ParseMode = "HTML"
			msg.ReplyMarkup = keyboard
			_, _ = botService.Bot.Send(msg)
		} else {
			// Send media with full caption
			_ = sendAnnouncementMessage(botService.Bot, chatID, announcement, caption, &keyboard)
		}
	}

	// Update pagination message
	totalCount, _ := botService.AnnouncementService.CountAdminAnnouncements(admin.ID)
	totalPages := (totalCount + announcementsPerPage - 1) / announcementsPerPage

	paginationText := fmt.Sprintf("📄 Sahifa %d/%d", page+1, totalPages)
	keyboard := utils.MakeAnnouncementListKeyboard(page, totalPages, lang)

	editMsg := tgbotapi.NewEditMessageText(chatID, callback.Message.MessageID, paginationText)
	editMsg.ReplyMarkup = &keyboard
	_, _ = botService.Bot.Send(editMsg)

	return botService.TelegramService.AnswerCallbackQuery(callback.ID, "")
}
