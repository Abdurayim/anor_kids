package utils

import (
	"fmt"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"anor-kids/internal/i18n"
	"anor-kids/internal/models"
)

// MakeLanguageKeyboard creates language selection keyboard
func MakeLanguageKeyboard() tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnUzbek, i18n.LanguageUzbek),
				"lang_uz",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnRussian, i18n.LanguageRussian),
				"lang_ru",
			),
		),
	)
}

// MakePhoneKeyboard creates phone number request keyboard
func MakePhoneKeyboard(lang i18n.Language) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButtonContact(i18n.Get(i18n.BtnSharePhone, lang)),
		),
	)
}

// MakeMainMenuKeyboard creates main menu keyboard
func MakeMainMenuKeyboard(lang i18n.Language) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.Get(i18n.BtnSubmitComplaint, lang)),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.Get(i18n.BtnSubmitProposal, lang)),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.Get(i18n.BtnMyComplaints, lang)),
			tgbotapi.NewKeyboardButton(i18n.Get(i18n.BtnMyProposals, lang)),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.Get(i18n.BtnViewAnnouncements, lang)),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.Get(i18n.BtnSettings, lang)),
		),
	)
}

// MakeMainMenuKeyboardWithAdmin creates main menu keyboard with admin button
func MakeMainMenuKeyboardWithAdmin(lang i18n.Language) tgbotapi.ReplyKeyboardMarkup {
	return tgbotapi.NewReplyKeyboard(
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.Get(i18n.BtnSubmitComplaint, lang)),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.Get(i18n.BtnSubmitProposal, lang)),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.Get(i18n.BtnMyComplaints, lang)),
			tgbotapi.NewKeyboardButton(i18n.Get(i18n.BtnMyProposals, lang)),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.Get(i18n.BtnViewAnnouncements, lang)),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.Get(i18n.BtnSettings, lang)),
		),
		tgbotapi.NewKeyboardButtonRow(
			tgbotapi.NewKeyboardButton(i18n.Get(i18n.BtnAdminPanel, lang)),
		),
	)
}

// MakeMainMenuKeyboardForUser creates main menu keyboard based on user's admin status
func MakeMainMenuKeyboardForUser(lang i18n.Language, isAdmin bool) tgbotapi.ReplyKeyboardMarkup {
	if isAdmin {
		return MakeMainMenuKeyboardWithAdmin(lang)
	}
	return MakeMainMenuKeyboard(lang)
}

// MakeConfirmationKeyboard creates confirmation keyboard
func MakeConfirmationKeyboard(lang i18n.Language) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnConfirm, lang),
				"confirm_complaint",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnCancel, lang),
				"cancel_complaint",
			),
		),
	)
}

// MakeImagePromptKeyboard creates keyboard to ask if user wants to add images
func MakeImagePromptKeyboard(lang i18n.Language) tgbotapi.InlineKeyboardMarkup {
	yesText := "✅ Ha, rasm qo'shaman"
	noText := "📤 Yo'q, rasmiz davom etish"

	if lang == i18n.LanguageRussian {
		yesText = "✅ Да, добавить фото"
		noText = "📤 Нет, продолжить без фото"
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				yesText,
				"add_images",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				noText,
				"skip_images",
			),
		),
	)
}

// MakeImageCollectionKeyboard creates keyboard for image collection flow (after user chooses to add images)
func MakeImageCollectionKeyboard(lang i18n.Language) tgbotapi.InlineKeyboardMarkup {
	finishText := "✅ Tugallash"

	if lang == i18n.LanguageRussian {
		finishText = "✅ Завершить"
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				finishText,
				"finish_images",
			),
		),
	)
}

// MakeAdminKeyboard creates admin panel keyboard
func MakeAdminKeyboard(lang i18n.Language) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnManageClasses, lang),
				"admin_manage_classes",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnViewUsers, lang),
				"admin_users",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnViewComplaints, lang),
				"admin_complaints",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnViewProposals, lang),
				"admin_proposals",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnViewStats, lang),
				"admin_stats",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnCreateAnnouncement, lang),
				"admin_create_announcement",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnManageAnnouncements, lang),
				"admin_manage_announcements",
			),
		),
	)
}

// RemoveKeyboard creates a keyboard removal markup
func RemoveKeyboard() tgbotapi.ReplyKeyboardRemove {
	return tgbotapi.NewRemoveKeyboard(true)
}

// MakeClassSelectionKeyboard creates class selection inline keyboard
func MakeClassSelectionKeyboard(classes []*models.Class, lang i18n.Language) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Create buttons in rows of 3
	var row []tgbotapi.InlineKeyboardButton
	for i, class := range classes {
		button := tgbotapi.NewInlineKeyboardButtonData(
			class.ClassName,
			"class_"+class.ClassName,
		)
		row = append(row, button)

		// Add row every 3 buttons or at the end
		if (i+1)%3 == 0 || i == len(classes)-1 {
			rows = append(rows, row)
			row = []tgbotapi.InlineKeyboardButton{}
		}
	}

	_ = lang // Will be used in future for localized buttons

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}

// MakeProposalImagePromptKeyboard creates keyboard to ask if user wants to add images for proposal
func MakeProposalImagePromptKeyboard(lang i18n.Language) tgbotapi.InlineKeyboardMarkup {
	yesText := "✅ Ha, rasm qo'shaman"
	noText := "📤 Yo'q, rasmiz davom etish"

	if lang == i18n.LanguageRussian {
		yesText = "✅ Да, добавить фото"
		noText = "📤 Нет, продолжить без фото"
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				yesText,
				"add_proposal_images",
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				noText,
				"skip_proposal_images",
			),
		),
	)
}

// MakeProposalImageCollectionKeyboard creates keyboard for proposal image collection flow
func MakeProposalImageCollectionKeyboard(lang i18n.Language) tgbotapi.InlineKeyboardMarkup {
	finishText := "✅ Tugallash"

	if lang == i18n.LanguageRussian {
		finishText = "✅ Завершить"
	}

	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				finishText,
				"finish_proposal_images",
			),
		),
	)
}

// MakeProposalConfirmationKeyboard creates proposal confirmation keyboard
func MakeProposalConfirmationKeyboard(lang i18n.Language) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnConfirm, lang),
				"confirm_proposal",
			),
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnCancel, lang),
				"cancel_proposal",
			),
		),
	)
}

// MakeAnnouncementDeleteKeyboard creates keyboard for deleting announcement
func MakeAnnouncementDeleteKeyboard(announcementID int, lang i18n.Language) tgbotapi.InlineKeyboardMarkup {
	return tgbotapi.NewInlineKeyboardMarkup(
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnDelete, lang),
				fmt.Sprintf("delete_announcement_%d", announcementID),
			),
		),
		tgbotapi.NewInlineKeyboardRow(
			tgbotapi.NewInlineKeyboardButtonData(
				i18n.Get(i18n.BtnBack, lang),
				"admin_back",
			),
		),
	)
}

// MakeAnnouncementListKeyboard creates keyboard for announcement pagination (admin view)
func MakeAnnouncementListKeyboard(currentPage, totalPages int, lang i18n.Language) tgbotapi.InlineKeyboardMarkup {
	var rows [][]tgbotapi.InlineKeyboardButton

	// Navigation row
	var navRow []tgbotapi.InlineKeyboardButton
	if currentPage > 0 {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData(
			"◀️",
			fmt.Sprintf("announcements_page_%d", currentPage-1),
		))
	}

	if currentPage < totalPages-1 {
		navRow = append(navRow, tgbotapi.NewInlineKeyboardButtonData(
			"▶️",
			fmt.Sprintf("announcements_page_%d", currentPage+1),
		))
	}

	if len(navRow) > 0 {
		rows = append(rows, navRow)
	}

	// Back button
	rows = append(rows, tgbotapi.NewInlineKeyboardRow(
		tgbotapi.NewInlineKeyboardButtonData(
			i18n.Get(i18n.BtnBack, lang),
			"admin_back",
		),
	))

	return tgbotapi.NewInlineKeyboardMarkup(rows...)
}
