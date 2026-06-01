package handlers

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"wallet_bot/app/utils"
)

func Help(bot *tgbotapi.BotAPI, update tgbotapi.Update, lang *utils.LanguageCache) {
	userLang := lang.Get(update.Message.From.ID)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, HelpText[userLang])
	msg.ParseMode = "HTML"

	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending help message:", err)
	}
}
