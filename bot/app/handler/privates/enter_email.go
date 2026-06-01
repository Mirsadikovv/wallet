package handlers

import (
	"fmt"
	"log"
	"regexp"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"

	"wallet_bot/app/config"
	keyboard "wallet_bot/app/keyboards/defaults"
	"wallet_bot/app/utils"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func EnterEmail(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg *config.Config, db *gorm.DB, lang *utils.LanguageCache, state *utils.StateCache, otp *utils.OTPCache) *utils.LanguageCache {
	userID := update.Message.From.ID
	userLang := lang.Get(userID)
	email := update.Message.Text

	if !emailRegex.MatchString(email) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, InvalidEmailText[userLang])
		msg.ParseMode = "HTML"
		if _, err := bot.Send(msg); err != nil {
			log.Println("Error sending message:", err)
		}
		return lang
	}

	var count int64
	db.Table("users").Where("email = ?", email).Count(&count)

	if count > 0 {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, EmailTakenText[userLang])
		msg.ParseMode = "HTML"
		if _, err := bot.Send(msg); err != nil {
			log.Println("Error sending message:", err)
		}
		return lang
	}

	code := otp.Generate(email)

	if err := utils.SendOTP(cfg, email, code); err != nil {
		log.Printf("Failed to send OTP to %s: %v", email, err)
	}

	state.Set(userID, "await_otp:"+email)

	text := fmt.Sprintf(OTPSentText[userLang], email)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = utils.MakeReplyMarkup(keyboard.GetBackKeyboard(userLang))

	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending message:", err)
	}

	return lang
}
