package handlers

import (
	"log"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"

	"wallet_bot/app/config"
	"wallet_bot/app/utils"
)

func Register(bot *tgbotapi.BotAPI, db *gorm.DB) {
	cmds := utils.GetBotCommands()
	if _, err := bot.Request(tgbotapi.NewSetMyCommands(cmds...)); err != nil {
		log.Println("Failed to set bot commands:", err)
	}
}

func HandleUpdate(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg *config.Config, db *gorm.DB, langCache *utils.LanguageCache, stateCache *utils.StateCache, otpCache *utils.OTPCache) *utils.LanguageCache {
	if update.Message.From.LanguageCode != "" && langCache.Get(update.Message.From.ID) == "" {
		lang := update.Message.From.LanguageCode
		if lang != "uz" && lang != "ru" && lang != "en" {
			lang = "en"
		}
		langCache.Set(update.Message.From.ID, lang)
	}

	if update.Message.IsCommand() {
		switch update.Message.Command() {
		case "start":
			return Start(bot, update, cfg, db, langCache, stateCache)
		case "help":
			Help(bot, update, langCache)
		default:
			Echo(bot, update, langCache)
		}
		return langCache
	}

	switch update.Message.Text {
	case "🇬🇧 English", "🇷🇺 Русский", "🇺🇿 O'zbek":
		return ChangeLanguage(bot, update, cfg, db, langCache, stateCache)
	case "Назад ⬅️", "Ortga ⬅️", "Back ⬅️":
		return Back(bot, update, cfg, langCache, stateCache)
	case "Change language🌐", "Tilni o'zgartirish🌐", "Поменять язык🌐":
		return ChooseLanguage(bot, update, langCache)
	}

	state := stateCache.Get(update.Message.From.ID)
	switch {
	case state == "await_email":
		return EnterEmail(bot, update, cfg, db, langCache, stateCache, otpCache)
	case strings.HasPrefix(state, "await_otp:"):
		return EnterOTP(bot, update, cfg, db, langCache, stateCache, otpCache)
	}

	Echo(bot, update, langCache)
	return langCache
}
