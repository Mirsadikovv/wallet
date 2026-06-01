package handlers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"

	"wallet_bot/app/config"
	keyboard "wallet_bot/app/keyboards/defaults"
	"wallet_bot/app/utils"
)

func ChooseLanguage(bot *tgbotapi.BotAPI, update tgbotapi.Update, lang *utils.LanguageCache) *utils.LanguageCache {
	userLang := lang.Get(update.Message.From.ID)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, ChooseLanguageText[userLang])
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = utils.MakeReplyMarkup(keyboard.LanguageKeyboard)

	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending message:", err)
	}

	return lang
}

func ChangeLanguage(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg *config.Config, db *gorm.DB, lang *utils.LanguageCache, state *utils.StateCache) *utils.LanguageCache {
	userID := update.Message.From.ID

	switch update.Message.Text {
	case "🇺🇿 O'zbek":
		lang.Set(userID, "uz")
	case "🇷🇺 Русский":
		lang.Set(userID, "ru")
	default:
		lang.Set(userID, "en")
	}

	var count int64
	db.Table("users").Where("telegram_id = ?", userID).Count(&count)

	if count > 0 {
		return Menu(bot, update, cfg, lang)
	}

	state.Set(userID, "await_email")

	userLang := lang.Get(userID)
	text := fmt.Sprintf(WelcomeText[userLang], userID, update.Message.From.FirstName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending message:", err)
	}

	return lang
}
