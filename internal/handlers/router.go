package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"anor-kids/internal/i18n"
	"anor-kids/internal/models"
	"anor-kids/internal/services"
	"anor-kids/internal/utils"
)

// RouteByState routes messages based on user's current state
func RouteByState(botService *services.BotService, message *tgbotapi.Message, state string, stateData *models.StateData) error {
	switch state {
	case models.StateAwaitingLanguage:
		// Waiting for language selection (handled by callback)
		return nil

	case models.StateAwaitingPhone:
		return HandlePhoneNumber(botService, message, stateData)

	case models.StateAwaitingChildName:
		return HandleChildName(botService, message, stateData)

	case models.StateAwaitingChildClass:
		return HandleChildClass(botService, message, stateData)

	case models.StateAwaitingComplaint:
		return HandleComplaintText(botService, message, stateData)

	case models.StateAwaitingImages:
		// Handle photo uploads
		if message.Photo != nil && len(message.Photo) > 0 {
			return HandleImage(botService, message, stateData)
		}
		// If not a photo, remind user
		lang := i18n.GetLanguage(stateData.Language)
		text := "📎 Iltimos, rasm yuboring (📎 tugmasidan) yoki 'Tugallash' tugmasini bosing.\n\n" +
			"📎 Пожалуйста, отправьте фото (через кнопку 📎) или нажмите 'Завершить'."
		keyboard := utils.MakeImageCollectionKeyboard(lang)
		return botService.TelegramService.SendMessage(message.Chat.ID, text, keyboard)

	case models.StateConfirmingComplaint:
		// Waiting for confirmation (handled by callback)
		return nil

	case models.StateAwaitingProposal:
		return HandleProposalText(botService, message, stateData)

	case models.StateAwaitingProposalImages:
		// Handle photo uploads for proposal
		if message.Photo != nil && len(message.Photo) > 0 {
			return HandleProposalImage(botService, message, stateData)
		}
		// If not a photo, remind user
		lang := i18n.GetLanguage(stateData.Language)
		text := "📎 Iltimos, rasm yuboring (📎 tugmasidan) yoki 'Tugallash' tugmasini bosing.\n\n" +
			"📎 Пожалуйста, отправьте фото (через кнопку 📎) или нажмите 'Завершить'."
		keyboard := utils.MakeProposalImageCollectionKeyboard(lang)
		return botService.TelegramService.SendMessage(message.Chat.ID, text, keyboard)

	case models.StateConfirmingProposal:
		// Waiting for confirmation (handled by callback)
		return nil

	case models.StateAwaitingAdminPhone:
		return HandleAdminLinkPhone(botService, message)

	case models.StateAwaitingClassName:
		return HandleClassNameInput(botService, message)

	case models.StateAwaitingAnnouncementTitle:
		return HandleAnnouncementTitle(botService, message, stateData)

	case models.StateAwaitingAnnouncementText:
		return HandleAnnouncementText(botService, message, stateData)

	case models.StateAwaitingAnnouncementImage:
		// Handle image uploads for announcement
		if message.Photo != nil && len(message.Photo) > 0 {
			return HandleAnnouncementImage(botService, message, stateData)
		}
		if message.Document != nil {
			return HandleAnnouncementImage(botService, message, stateData)
		}
		// If not an image, remind user
		text := "📸 Iltimos, rasm yuboring.\n\n📸 Пожалуйста, отправьте изображение."
		return botService.TelegramService.SendMessage(message.Chat.ID, text, nil)

	case models.StateRegistered:
		// User is registered, get user data
		user, err := botService.UserService.GetUserByTelegramID(message.From.ID)
		if err != nil {
			// Show error to user instead of silent failure
			lang := i18n.LanguageUzbek
			text := i18n.Get(i18n.ErrDatabaseError, lang)
			_ = botService.TelegramService.SendMessage(message.Chat.ID, text, nil)
			return err
		}
		if user == nil {
			// User not found, restart registration
			return HandleStart(botService, message)
		}
		return HandleRegisteredUserMessage(botService, message, user)

	default:
		// Unknown state, restart
		return HandleStart(botService, message)
	}
}

// HandleRegisteredUserMessage handles messages from registered users
func HandleRegisteredUserMessage(botService *services.BotService, message *tgbotapi.Message, user *models.User) error {
	if user == nil {
		return HandleStart(botService, message)
	}

	lang := i18n.GetLanguage(user.Language)
	chatID := message.Chat.ID

	// Check if message is a button press
	buttonText := message.Text

	// Submit complaint button (check both languages)
	submitBtnUz := i18n.Get(i18n.BtnSubmitComplaint, i18n.LanguageUzbek)
	submitBtnRu := i18n.Get(i18n.BtnSubmitComplaint, i18n.LanguageRussian)
	if buttonText == submitBtnUz || buttonText == submitBtnRu {
		return HandleComplaintCommand(botService, message)
	}

	// Submit proposal button (check both languages)
	submitProposalBtnUz := i18n.Get(i18n.BtnSubmitProposal, i18n.LanguageUzbek)
	submitProposalBtnRu := i18n.Get(i18n.BtnSubmitProposal, i18n.LanguageRussian)
	if buttonText == submitProposalBtnUz || buttonText == submitProposalBtnRu {
		return HandleProposalCommand(botService, message)
	}

	// My complaints button (check both languages)
	myComplaintsBtnUz := i18n.Get(i18n.BtnMyComplaints, i18n.LanguageUzbek)
	myComplaintsBtnRu := i18n.Get(i18n.BtnMyComplaints, i18n.LanguageRussian)
	if buttonText == myComplaintsBtnUz || buttonText == myComplaintsBtnRu {
		return HandleMyComplaintsCommand(botService, message)
	}

	// My proposals button (check both languages)
	myProposalsBtnUz := i18n.Get(i18n.BtnMyProposals, i18n.LanguageUzbek)
	myProposalsBtnRu := i18n.Get(i18n.BtnMyProposals, i18n.LanguageRussian)
	if buttonText == myProposalsBtnUz || buttonText == myProposalsBtnRu {
		return HandleMyProposalsCommand(botService, message)
	}

	// Settings button (check both languages)
	settingsBtnUz := i18n.Get(i18n.BtnSettings, i18n.LanguageUzbek)
	settingsBtnRu := i18n.Get(i18n.BtnSettings, i18n.LanguageRussian)
	if buttonText == settingsBtnUz || buttonText == settingsBtnRu {
		return HandleSettingsCommand(botService, message)
	}

	// View announcements button (check both languages)
	announcementsBtnUz := i18n.Get(i18n.BtnViewAnnouncements, i18n.LanguageUzbek)
	announcementsBtnRu := i18n.Get(i18n.BtnViewAnnouncements, i18n.LanguageRussian)
	if buttonText == announcementsBtnUz || buttonText == announcementsBtnRu {
		return HandleViewAnnouncements(botService, message)
	}

	// Default: show main menu
	text := i18n.Get(i18n.MsgMainMenu, lang)

	// Check if user is admin to show appropriate keyboard
	isAdmin, _ := botService.IsAdmin(user.PhoneNumber, user.TelegramID)
	keyboard := utils.MakeMainMenuKeyboardForUser(lang, isAdmin)

	return botService.TelegramService.SendMessage(chatID, text, keyboard)
}

// HandleCallbackQuery handles inline button clicks
func HandleCallbackQuery(botService *services.BotService, callback *tgbotapi.CallbackQuery) error {
	data := callback.Data

	// Language selection
	if data == "lang_uz" || data == "lang_ru" {
		return HandleLanguageSelection(botService, callback)
	}

	// Class action callbacks (activate, deactivate, delete) - MUST CHECK BEFORE generic "class_"
	if len(data) > 13 && data[:13] == "class_delete_" {
		return HandleClassDeleteCallback(botService, callback)
	}

	if len(data) > 13 && data[:13] == "class_toggle_" {
		return HandleClassToggleCallback(botService, callback)
	}

	// Class selection (starts with "class_") - for user registration
	if len(data) > 6 && data[:6] == "class_" {
		return HandleClassSelection(botService, callback)
	}

	// Complaint confirmation
	if data == "confirm_complaint" {
		return HandleComplaintConfirmation(botService, callback)
	}

	// Complaint cancellation
	if data == "cancel_complaint" {
		return HandleComplaintCancellation(botService, callback)
	}

	// Image collection callbacks
	if data == "add_images" {
		return HandleAddImages(botService, callback)
	}

	if data == "skip_images" {
		return HandleSkipImages(botService, callback)
	}

	if data == "finish_images" {
		return HandleFinishImages(botService, callback)
	}

	// Proposal confirmation
	if data == "confirm_proposal" {
		return HandleProposalConfirmation(botService, callback)
	}

	// Proposal cancellation
	if data == "cancel_proposal" {
		return HandleProposalCancellation(botService, callback)
	}

	// Proposal image collection callbacks
	if data == "add_proposal_images" {
		return HandleAddProposalImages(botService, callback)
	}

	if data == "skip_proposal_images" {
		return HandleSkipProposalImages(botService, callback)
	}

	if data == "finish_proposal_images" {
		return HandleFinishProposalImages(botService, callback)
	}

	// Admin callbacks
	if data == "admin_users" {
		return HandleAdminUsersCallback(botService, callback)
	}

	if data == "admin_complaints" {
		return HandleAdminComplaintsCallback(botService, callback)
	}

	if data == "admin_proposals" {
		return HandleAdminProposalsCallback(botService, callback)
	}

	if data == "admin_stats" {
		return HandleAdminStatsCallback(botService, callback)
	}

	// Admin manage classes callback
	if data == "admin_manage_classes" {
		return HandleAdminManageClassesCallback(botService, callback)
	}

	// Admin create class callback
	if data == "admin_create_class" {
		return HandleAdminCreateClassCallback(botService, callback)
	}

	// Admin back button
	if data == "admin_back" {
		return HandleAdminBackCallback(botService, callback)
	}

	// Admin create announcement callback
	if data == "admin_create_announcement" {
		return HandleAdminCreateAnnouncementCallback(botService, callback)
	}

	// Admin manage announcements callback
	if data == "admin_manage_announcements" {
		return HandleAdminManageAnnouncementsCallback(botService, callback)
	}

	// Delete announcement callbacks (starts with "delete_announcement_")
	if len(data) > 20 && data[:20] == "delete_announcement_" {
		return HandleDeleteAnnouncementCallback(botService, callback)
	}

	// Announcement pagination callbacks (starts with "announcements_page_")
	if len(data) > 18 && data[:18] == "announcements_page" {
		return HandleAnnouncementPageCallback(botService, callback)
	}

	// Unknown callback
	return botService.TelegramService.AnswerCallbackQuery(callback.ID, "Unknown action")
}
