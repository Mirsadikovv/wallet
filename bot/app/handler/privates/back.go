package handlers

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

	"wallet_bot/app/config"
	"wallet_bot/app/utils"
)

func Back(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg *config.Config, lang *utils.LanguageCache, state *utils.StateCache) *utils.LanguageCache {
	state.Clear(update.Message.From.ID)
	return Menu(bot, update, cfg, lang)
}
