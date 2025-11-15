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

// HandleProposalCommand initiates proposal submission
func HandleProposalCommand(botService *services.BotService, message *tgbotapi.Message) error {
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

	// Set state to awaiting proposal
	stateData := &models.StateData{
		Language: user.Language,
		Images:   []models.ImageData{},
	}
	err = botService.StateManager.Set(telegramID, models.StateAwaitingProposal, stateData)
	if err != nil {
		return err
	}

	// Send request message
	text := i18n.Get(i18n.MsgRequestProposal, lang)
	return botService.TelegramService.SendMessage(chatID, text, nil)
}

// HandleProposalText handles proposal text input
func HandleProposalText(botService *services.BotService, message *tgbotapi.Message, stateData *models.StateData) error {
	telegramID := message.From.ID
	chatID := message.Chat.ID
	lang := i18n.GetLanguage(stateData.Language)

	// Validate proposal text
	proposalText, err := validator.ValidateProposalText(message.Text)
	if err != nil {
		text := i18n.Get(i18n.ErrInvalidProposal, lang) + "\n\n" + err.Error()
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	// Save proposal text in state
	stateData.ProposalText = proposalText

	// Move to image collection state
	err = botService.StateManager.Set(telegramID, models.StateAwaitingProposalImages, stateData)
	if err != nil {
		return err
	}

	// Ask if user wants to add images
	text := "✅ Taklif matni qabul qilindi.\n\n" +
		"Taklifingizga rasm qo'shmoqchimisiz?\n\n" +
		"✅ Текст предложения принят.\n\n" +
		"Хотите добавить фото к предложению?"

	keyboard := utils.MakeProposalImagePromptKeyboard(lang)
	return botService.TelegramService.SendMessage(chatID, text, keyboard)
}

// HandleProposalImage handles image input during proposal submission
func HandleProposalImage(botService *services.BotService, message *tgbotapi.Message, stateData *models.StateData) error {
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
	err := botService.StateManager.Set(telegramID, models.StateAwaitingProposalImages, stateData)
	if err != nil {
		return err
	}

	// Send confirmation
	text := fmt.Sprintf("✅ Rasm qabul qilindi (%d/10)\n\n✅ Изображение принято (%d/10)",
		len(stateData.Images), len(stateData.Images))

	keyboard := utils.MakeProposalImageCollectionKeyboard(lang)
	return botService.TelegramService.SendMessage(chatID, text, keyboard)
}

// HandleAddProposalImages handles when user chooses to add images
func HandleAddProposalImages(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
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

	keyboard := utils.MakeProposalImageCollectionKeyboard(lang)
	return botService.TelegramService.SendMessage(chatID, text, keyboard)
}

// HandleSkipProposalImages handles skipping image upload
func HandleSkipProposalImages(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
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
	err = botService.StateManager.Set(telegramID, models.StateConfirmingProposal, stateData)
	if err != nil {
		return err
	}

	// Show confirmation with proposal text and images
	text := "📋 <b>Taklifingizni tekshiring / Проверьте ваше предложение</b>\n\n"
	text += "<b>Matn / Текст:</b>\n"
	text += stateData.ProposalText + "\n\n"
	text += fmt.Sprintf("📎 <b>Rasmlar / Фото:</b> %d ta\n\n", len(stateData.Images))

	if len(stateData.Images) == 0 {
		text += "──────────\n\n"
		text += "Yuborilsinmi? / Отправить?"
	} else {
		text += "──────────\n\n"
		text += "Taklif va rasmlar yuborilsinmi?\nОтправить предложение с фотографиями?"
	}

	keyboard := utils.MakeProposalConfirmationKeyboard(lang)
	return botService.TelegramService.SendMessage(chatID, text, keyboard)
}

// HandleFinishProposalImages handles finishing image upload
func HandleFinishProposalImages(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
	return HandleSkipProposalImages(botService, callback) // Same logic
}

// HandleProposalConfirmation handles proposal confirmation and PDF generation
func HandleProposalConfirmation(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
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

	// Get proposal text and images from state
	stateData, err := botService.StateManager.GetData(telegramID)
	if err != nil {
		return err
	}

	if stateData.ProposalText == "" {
		return botService.TelegramService.AnswerCallbackQuery(callback.ID, "Proposal text not found")
	}

	// Answer callback query
	_ = botService.TelegramService.AnswerCallbackQuery(callback.ID, "✅")

	// Generate PDF document with text and images
	pdfPath, filename, err := botService.DocumentService.GenerateProposalPDF(user, stateData.ProposalText, stateData.Images)
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

	// Save proposal to database with PDF info
	proposalReq := &models.CreateProposalRequest{
		UserID:            user.ID,
		ProposalText:      stateData.ProposalText,
		PDFTelegramFileID: fileID,
		PDFFilename:       filename,
	}

	proposal, err := botService.ProposalService.CreateProposal(proposalReq)
	if err != nil {
		log.Printf("Failed to save proposal: %v", err)
		text := i18n.Get(i18n.ErrDatabaseError, lang)
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	// Save images to database
	for i, img := range stateData.Images {
		err = botService.ProposalService.CreateProposalImage(proposal.ID, &img, i)
		if err != nil {
			log.Printf("Failed to save proposal image %d: %v", i, err)
		}
	}

	// Clear state
	_ = botService.StateManager.Clear(telegramID)

	// Send success message
	text := i18n.Get(i18n.MsgProposalSubmitted, lang)

	// Check if user is admin to show appropriate keyboard
	isAdmin, _ := botService.IsAdmin(user.PhoneNumber, user.TelegramID)
	keyboard := utils.MakeMainMenuKeyboardForUser(lang, isAdmin)
	_ = botService.TelegramService.SendMessage(chatID, text, keyboard)

	// Notify admins with PDF document
	go notifyAdminsWithProposalPDF(botService, user, proposal, fileID, len(stateData.Images))

	return nil
}

// HandleProposalCancellation handles proposal cancellation
func HandleProposalCancellation(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
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
	_ = botService.TelegramService.AnswerCallbackQuery(callback.ID, i18n.Get(i18n.MsgProposalCancelled, lang))

	// Send cancellation message
	text := i18n.Get(i18n.MsgProposalCancelled, lang)

	// Check if user is admin to show appropriate keyboard
	isAdmin, _ := botService.IsAdmin(user.PhoneNumber, user.TelegramID)
	keyboard := utils.MakeMainMenuKeyboardForUser(lang, isAdmin)

	return botService.TelegramService.SendMessage(chatID, text, keyboard)
}

// notifyAdminsWithProposalPDF sends proposal as PDF document to all admins
func notifyAdminsWithProposalPDF(botService *services.BotService, user *models.User, proposal *models.Proposal, fileID string, imageCount int) {
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
		"<b>YANGI TAKLIF / НОВОЕ ПРЕДЛОЖЕНИЕ</b>\n\n"+
			"ID: #%d\n"+
			"Farzand / Ребенок: <b>%s</b>\n"+
			"Sinf / Класс: <b>%s</b>\n"+
			"Telefon / Телефон: %s\n"+
			"Username: @%s\n"+
			"Sana / Дата: %s\n"+
			"📷 Rasmlar / Изображений: %d\n\n"+
			"Taklif PDF hujjat sifatida yuqorida\n"+
			"Предложение в формате PDF выше",
		proposal.ID,
		user.ChildName,
		user.ChildClass,
		user.PhoneNumber,
		username,
		utils.FormatDateTime(proposal.CreatedAt),
		imageCount,
	)

	// Send document to all admins
	err = botService.TelegramService.SendDocumentToAdmins(adminIDs, fileID, caption)
	if err != nil {
		log.Printf("Failed to send document to admins: %v", err)
	}
}

// HandleMyProposalsCommand shows user's proposal history
func HandleMyProposalsCommand(botService *services.BotService, message *tgbotapi.Message) error {
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

	// Get user proposals
	proposals, err := botService.ProposalService.GetUserProposals(user.ID, 10, 0)
	if err != nil {
		text := i18n.Get(i18n.ErrDatabaseError, lang)
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	if len(proposals) == 0 {
		text := "Sizda hali takliflar yo'q / У вас пока нет предложений"
		return botService.TelegramService.SendMessage(chatID, text, nil)
	}

	// Format proposals list
	text := "📋 Sizning takliflaringiz / Ваши предложения:\n\n"
	for i, p := range proposals {
		status := "⏳"
		if p.Status == models.ProposalStatusReviewed {
			status = "✅"
		}

		preview := utils.TruncateText(p.ProposalText, 50)
		text += fmt.Sprintf("%d. %s %s\n   📅 %s\n\n",
			i+1,
			status,
			preview,
			utils.FormatDateTime(p.CreatedAt),
		)
	}

	return botService.TelegramService.SendMessage(chatID, text, nil)
}
