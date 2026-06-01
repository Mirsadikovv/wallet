package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"wallet_bot/app/utils"
)

func Echo(bot *tgbotapi.BotAPI, update tgbotapi.Update, lang *utils.LanguageCache) {
	userLang := lang.Get(update.Message.From.ID)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, EchoText[userLang])
	msg.ParseMode = "HTML"

	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending echo:", err)
	}
}
