package handlers

import (
	"fmt"
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"wallet_bot/app/config"
	keyboard "wallet_bot/app/keyboards/defaults"
	"wallet_bot/app/utils"
)

func Menu(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg *config.Config, lang *utils.LanguageCache) *utils.LanguageCache {
	userID := update.Message.From.ID
	userLang := lang.Get(userID)

	text := fmt.Sprintf(MenuText[userLang], userID, update.Message.From.FirstName)
	msg := tgbotapi.NewMessage(update.Message.Chat.ID, text)
	msg.ParseMode = "HTML"
	msg.ReplyMarkup = utils.MakeReplyMarkup(keyboard.GetMainMenuKeyboard(userLang, cfg.WebAppURL))

	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending menu:", err)
	}

	return lang
}
