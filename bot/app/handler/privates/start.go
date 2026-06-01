package handlers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"

	"wallet_bot/app/config"
	"wallet_bot/app/utils"
)

func Start(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg *config.Config, db *gorm.DB, lang *utils.LanguageCache, state *utils.StateCache) *utils.LanguageCache {
	userID := update.Message.From.ID
	userLang := lang.Get(userID)

	var count int64
	db.Table("users").Where("telegram_id = ?", userID).Count(&count)

	if count > 0 {
		state.Clear(userID)
		return Menu(bot, update, cfg, lang)
	}

	state.Set(userID, "await_email")

	text := fmt.Sprintf(WelcomeText[userLang], userID, update.Message.From.FirstName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"

	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending start message:", err)
	}

	return lang
}
