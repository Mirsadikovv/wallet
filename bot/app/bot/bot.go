package bot

import (
	"log"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"

	"wallet_bot/app/config"
	handlers "wallet_bot/app/handler/privates"
	"wallet_bot/app/utils"
)

func Start(cfg *config.Config, db *gorm.DB) error {
	bot, err := tgbotapi.NewBotAPI(cfg.BotToken)
	{
		if err != nil {
			log.Printf("failed to create bot: %v", err)
			return err
		}
	}

	bot.Debug = cfg.Debug

	u := tgbotapi.NewUpdate(0)
	u.Timeout = 100

	handlers.Register(bot, db)
	utils.NotifyAdmins(bot, cfg)

	updates := bot.GetUpdatesChan(u)
	langCache := utils.NewLanguageCache()
	stateCache := utils.NewStateCache()
	otpCache := utils.NewOTPCache()

	for update := range updates {
		if update.Message != nil {
			curLang := handlers.HandleUpdate(bot, update, cfg, db, langCache, stateCache, otpCache)
			langCache.Set(update.Message.From.ID, curLang.Get(update.Message.From.ID))
		}
	}

	return nil
}
