package utils

import (
	"log"
	"strconv"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"wallet_bot/app/config"
)

func NotifyAdmins(bot *tgbotapi.BotAPI, cfg *config.Config) {
	for _, idStr := range cfg.AdminIds {
		id, err := strconv.ParseInt(idStr, 10, 64)
		{
			if err != nil {
				log.Printf("Invalid admin ID: %v", err)
				continue
			}
		}

		msg := tgbotapi.NewMessage(id, "🤖 Bot started successfully")
		if _, err := bot.Send(msg); err != nil {
			log.Printf("Failed to notify admin %d: %v", id, err)
		}
	}
}
