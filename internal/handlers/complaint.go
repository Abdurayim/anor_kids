package handlers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"anor-kids/internal/i18n"
	"anor-kids/internal/models"
	"anor-kids/internal/services"
	"anor-kids/internal/utils"
	"anor-kids/internal/validator"
)

// HandleComplaintCommand initiates complaint submission
func HandleComplaintCommand(botService *services.BotService, message *tgbotapi.Message) error {
	telegramID := message.From.ID
	chatID := message.Chat.ID

	// Check if user is registered
	user, err := botService.UserService.GetUserByTelegramID(telegramID)
	if err != nil {
		return err
	}

	if user == nil {
		lang := i18n.LanguageUzbek
		text := i18n.Get(i18n.ErrNotRegistered, lang)
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	lang := i18n.GetLanguage(user.Language)

	// Set state to awaiting complaint
	stateData := &models.StateData{
		Language: user.Language,
		Images:   []models.ImageData{},
	}
	err = botService.StateManager.Set(telegramID, models.StateAwaitingComplaint, stateData)
	if err != nil {
		return err
	}

	// Send request message
	text := i18n.Get(i18n.MsgRequestComplaint, lang)
	return botService.TelegramService.SendMessage(chatID, text, nil)
}

// HandleComplaintText handles complaint text input
func HandleComplaintText(botService *services.BotService, message *tgbotapi.Message, stateData *models.StateData) error {
	telegramID := message.From.ID
	chatID := message.Chat.ID
	lang := i18n.GetLanguage(stateData.Language)

	// Validate complaint text
	complaintText, err := validator.ValidateComplaintText(message.Text)
	if err != nil {
		text := i18n.Get(i18n.ErrInvalidComplaint, lang) + "\n\n" + err.Error()
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	// Save complaint text in state
	stateData.ComplaintText = complaintText

	// Move to image collection state
	err = botService.StateManager.Set(telegramID, models.StateAwaitingImages, stateData)
	if err != nil {
		return err
	}

	// Ask if user wants to add images
	text := "✅ Shikoyat matni qabul qilindi.\n\n" +
		"Shikoyatingizga rasm qo'shmoqchimisiz?\n\n" +
		"✅ Текст жалобы принят.\n\n" +
		"Хотите добавить фото к жалобе?"

	keyboard := utils.MakeImagePromptKeyboard(lang)
	return botService.TelegramService.SendMessage(chatID, text, keyboard)
}

// HandleImage handles image input during complaint submission
func HandleImage(botService *services.BotService, message *tgbotapi.Message, stateData *models.StateData) error {
	chatID := message.Chat.ID
	telegramID := message.From.ID
	lang := i18n.GetLanguage(stateData.Language)

	// Check if max images reached
	if len(stateData.Images) >= 10 {
		text := "Maksimum rasm soni (10 ta) yuklandi.\n" +
			"'Tugallash' tugmasini bosing.\n\n" +
			"Достигнуто максимальное количество изображений (10).\n" +
			"Нажмите кнопку 'Завершить'."
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	// Get largest photo size
	var photo *tgbotapi.PhotoSize
	if len(message.Photo) > 0 {
		photo = &message.Photo[len(message.Photo)-1]
	} else {
		text := "Iltimos, rasm yuboring.\n\nПожалуйста, отправьте изображение."
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	// Add image to state
	imageData := models.ImageData{
		FileID:       photo.FileID,
		FileUniqueID: photo.FileUniqueID,
		FileSize:     photo.FileSize,
	}
	stateData.Images = append(stateData.Images, imageData)

	// Save state
	err := botService.StateManager.Set(telegramID, models.StateAwaitingImages, stateData)
	if err != nil {
		return err
	}

	// Send confirmation
	text := fmt.Sprintf("✅ Rasm qabul qilindi (%d/10)\n\n✅ Изображение принято (%d/10)",
		len(stateData.Images), len(stateData.Images))

	keyboard := utils.MakeImageCollectionKeyboard(lang)
	return botService.TelegramService.SendMessage(chatID, text, keyboard)
}

// HandleAddImages handles when user chooses to add images
func HandleAddImages(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
	telegramID := callback.From.ID
	chatID := callback.Message.Chat.ID

	// Get state
	stateData, err := botService.StateManager.GetData(telegramID)
	if err != nil {
		return err
	}

	lang := i18n.GetLanguage(stateData.Language)

	// Answer callback
	_ = botService.TelegramService.AnswerCallbackQuery(callback.ID, "✅")

	// Send instructions for adding images using file/document icon
	text := "📎 Rasmlarni yuboring (maksimum 10 ta).\n\n" +
		"💡 <b>Yo'riqnoma:</b> Telegramda rasmni yuborish uchun 📎 (qisqich) tugmasini bosing va rasmlarni tanlang.\n\n" +
		"Rasmlar yuborilgandan keyin 'Tugallash' tugmasini bosing.\n\n" +
		"──────────\n\n" +
		"📎 Отправьте фотографии (максимум 10).\n\n" +
		"💡 <b>Инструкция:</b> Чтобы отправить фото в Telegram, нажмите кнопку 📎 (скрепка) и выберите фотографии.\n\n" +
		"После отправки фото нажмите кнопку 'Завершить'."

	keyboard := utils.MakeImageCollectionKeyboard(lang)
	return botService.TelegramService.SendMessage(chatID, text, keyboard)
}

// HandleSkipImages handles skipping image upload
func HandleSkipImages(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
	telegramID := callback.From.ID
	chatID := callback.Message.Chat.ID

	// Get state
	stateData, err := botService.StateManager.GetData(telegramID)
	if err != nil {
		return err
	}

	lang := i18n.GetLanguage(stateData.Language)

	// Answer callback
	_ = botService.TelegramService.AnswerCallbackQuery(callback.ID, "✅")

	// Move to confirmation
	err = botService.StateManager.Set(telegramID, models.StateConfirmingComplaint, stateData)
	if err != nil {
		return err
	}

	// Show confirmation with complaint text and images
	text := "📋 <b>Shikoyatingizni tekshiring / Проверьте вашу жалобу</b>\n\n"
	text += "<b>Matn / Текст:</b>\n"
	text += stateData.ComplaintText + "\n\n"
	text += fmt.Sprintf("📎 <b>Rasmlar / Фото:</b> %d ta\n\n", len(stateData.Images))

	if len(stateData.Images) == 0 {
		text += "──────────\n\n"
		text += "Yuborilsinmi? / Отправить?"
	} else {
		text += "──────────\n\n"
		text += "Shikoyat va rasmlar yuborilsinmi?\nОтправить жалобу с фотографиями?"
	}

	keyboard := utils.MakeConfirmationKeyboard(lang)
	return botService.TelegramService.SendMessage(chatID, text, keyboard)
}

// HandleFinishImages handles finishing image upload
func HandleFinishImages(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
	return HandleSkipImages(botService, callback) // Same logic
}

// HandleComplaintConfirmation handles complaint confirmation and PDF generation
func HandleComplaintConfirmation(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
	telegramID := callback.From.ID
	chatID := callback.Message.Chat.ID

	// Get user
	user, err := botService.UserService.GetUserByTelegramID(telegramID)
	if err != nil {
		return err
	}

	if user == nil {
		return botService.TelegramService.AnswerCallbackQuery(callback.ID, "User not found")
	}

	lang := i18n.GetLanguage(user.Language)

	// Get complaint text and images from state
	stateData, err := botService.StateManager.GetData(telegramID)
	if err != nil {
		return err
	}

	if stateData.ComplaintText == "" {
		return botService.TelegramService.AnswerCallbackQuery(callback.ID, "Complaint text not found")
	}

	// Answer callback query
	_ = botService.TelegramService.AnswerCallbackQuery(callback.ID, "✅")

	// Generate PDF document with text and images
	pdfPath, filename, err := botService.DocumentService.GenerateComplaintPDF(user, stateData.ComplaintText, stateData.Images)
	if err != nil {
		log.Printf("Failed to generate PDF: %v", err)
		text := i18n.Get(i18n.ErrDatabaseError, lang)
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	// Upload PDF to Telegram and get file_id
	fileID, err := botService.TelegramService.UploadDocument(chatID, pdfPath, filename)
	if err != nil {
		log.Printf("Failed to upload PDF: %v", err)
		// Clean up temp file
		_ = botService.DocumentService.DeleteTempFile(pdfPath)
		text := i18n.Get(i18n.ErrDatabaseError, lang)
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	// Clean up temp file after upload
	_ = botService.DocumentService.DeleteTempFile(pdfPath)

	// Save complaint to database with PDF info
	complaintReq := &models.CreateComplaintRequest{
		UserID:            user.ID,
		ComplaintText:     stateData.ComplaintText,
		PDFTelegramFileID: fileID,
		PDFFilename:       filename,
	}

	complaint, err := botService.ComplaintService.CreateComplaint(complaintReq)
	if err != nil {
		log.Printf("Failed to save complaint: %v", err)
		text := i18n.Get(i18n.ErrDatabaseError, lang)
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	// Save images to database
	for i, img := range stateData.Images {
		err = botService.ComplaintService.CreateComplaintImage(complaint.ID, &img, i)
		if err != nil {
			log.Printf("Failed to save complaint image %d: %v", i, err)
		}
	}

	// Clear state
	_ = botService.StateManager.Clear(telegramID)

	// Send success message
	text := i18n.Get(i18n.MsgComplaintSubmitted, lang)

	// Check if user is admin to show appropriate keyboard
	isAdmin, _ := botService.IsAdmin(user.PhoneNumber, user.TelegramID)
	keyboard := utils.MakeMainMenuKeyboardForUser(lang, isAdmin)
	_ = botService.TelegramService.SendMessage(chatID, text, keyboard)

	// Notify admins with PDF document
	go notifyAdminsWithPDF(botService, user, complaint, fileID, len(stateData.Images))

	return nil
}

// HandleComplaintCancellation handles complaint cancellation
func HandleComplaintCancellation(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
	telegramID := callback.From.ID
	chatID := callback.Message.Chat.ID

	// Get user
	user, err := botService.UserService.GetUserByTelegramID(telegramID)
	if err != nil {
		return err
	}

	if user == nil {
		return botService.TelegramService.AnswerCallbackQuery(callback.ID, "User not found")
	}

	lang := i18n.GetLanguage(user.Language)

	// Clear state
	_ = botService.StateManager.Clear(telegramID)

	// Answer callback query
	_ = botService.TelegramService.AnswerCallbackQuery(callback.ID, i18n.Get(i18n.MsgComplaintCancelled, lang))

	// Send cancellation message
	text := i18n.Get(i18n.MsgComplaintCancelled, lang)

	// Check if user is admin to show appropriate keyboard
	isAdmin, _ := botService.IsAdmin(user.PhoneNumber, user.TelegramID)
	keyboard := utils.MakeMainMenuKeyboardForUser(lang, isAdmin)

	return botService.TelegramService.SendMessage(chatID, text, keyboard)
}

// notifyAdminsWithPDF sends complaint as PDF document to all admins
func notifyAdminsWithPDF(botService *services.BotService, user *models.User, complaint *models.Complaint, fileID string, imageCount int) {
	// Get admin telegram IDs
	adminIDs, err := botService.GetAdminTelegramIDs()
	if err != nil {
		log.Printf("Failed to get admin IDs: %v", err)
		return
	}

	if len(adminIDs) == 0 {
		log.Println("No admins configured")
		return
	}

	// Generate caption for the document
	username := user.TelegramUsername
	if username == "" {
		username = "yo'q / нет"
	}

	caption := fmt.Sprintf(
		"<b>YANGI SHIKOYAT / НОВАЯ ЖАЛОБА</b>\n\n"+
			"ID: #%d\n"+
			"Farzand / Ребенок: <b>%s</b>\n"+
			"Sinf / Класс: <b>%s</b>\n"+
			"Telefon / Телефон: %s\n"+
			"Username: @%s\n"+
			"Sana / Дата: %s\n"+
			"📷 Rasmlar / Изображений: %d\n\n"+
			"Shikoyat PDF hujjat sifatida yuqorida\n"+
			"Жалоба в формате PDF выше",
		complaint.ID,
		user.ChildName,
		user.ChildClass,
		user.PhoneNumber,
		username,
		utils.FormatDateTime(complaint.CreatedAt),
		imageCount,
	)

	// Send document to all admins
	err = botService.TelegramService.SendDocumentToAdmins(adminIDs, fileID, caption)
	if err != nil {
		log.Printf("Failed to send document to admins: %v", err)
	}
}

// HandleMyComplaintsCommand shows user's complaint history
func HandleMyComplaintsCommand(botService *services.BotService, message *tgbotapi.Message) error {
	telegramID := message.From.ID
	chatID := message.Chat.ID

	// Get user
	user, err := botService.UserService.GetUserByTelegramID(telegramID)
	if err != nil {
		return err
	}

	if user == nil {
		lang := i18n.LanguageUzbek
		text := i18n.Get(i18n.ErrNotRegistered, lang)
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	lang := i18n.GetLanguage(user.Language)

	// Get user complaints
	complaints, err := botService.ComplaintService.GetUserComplaints(user.ID, 10, 0)
	if err != nil {
		text := i18n.Get(i18n.ErrDatabaseError, lang)
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	if len(complaints) == 0 {
		text := "Sizda hali shikoyatlar yo'q / У вас пока нет жалоб"
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	// Format complaints list
	text := "📋 Sizning shikoyatlaringiz / Ваши жалобы:\n\n"
	for i, c := range complaints {
		status := "⏳"
		if c.Status == models.StatusReviewed {
			status = "✅"
		}

		preview := utils.TruncateText(c.ComplaintText, 50)
		text += fmt.Sprintf("%d. %s %s\n   📅 %s\n\n",
			i+1,
			status,
			preview,
			utils.FormatDateTime(c.CreatedAt),
		)
	}

	return botService.TelegramService.SendMessage(chatID, text, nil)
}

// HandleSettingsCommand shows settings menu
func HandleSettingsCommand(botService *services.BotService, message *tgbotapi.Message) error {
	telegramID := message.From.ID
	chatID := message.Chat.ID

	// Get user
	user, err := botService.UserService.GetUserByTelegramID(telegramID)
	if err != nil {
		return err
	}

	if user == nil {
		lang := i18n.LanguageUzbek
		text := i18n.Get(i18n.ErrNotRegistered, lang)
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	lang := i18n.GetLanguage(user.Language)
	_ = lang // Will be used in future for localized settings

	// Format user info
	text := "⚙️ Sozlamalar / Настройки\n\n"
	text += fmt.Sprintf("👤 Farzand / Ребенок: %s\n", user.ChildName)
	text += fmt.Sprintf("🎓 Sinf / Класс: %s\n", user.ChildClass)
	text += fmt.Sprintf("📱 Telefon / Телефон: %s\n", utils.FormatPhoneNumber(user.PhoneNumber))
	text += fmt.Sprintf("🌍 Til / Язык: %s\n", user.Language)

	return botService.TelegramService.SendMessage(chatID, text, nil)
}
