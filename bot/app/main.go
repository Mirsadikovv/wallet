package app

import (
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"wallet_bot/app/bot"
	"wallet_bot/app/config"
)

func Exec(cfg *config.Config) {
	dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%d sslmode=%s TimeZone=%s",
		cfg.PostgresHost,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		cfg.PostgresDB,
		cfg.PostgresPort,
		cfg.PostgresSSLMode,
		cfg.PostgresTimeZone,
	)

	gormLogLevel := logger.Silent
	if cfg.Debug {
		gormLogLevel = logger.Info
	}

	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(gormLogLevel),
	})
	{
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}
	}

	if err := bot.Start(cfg, db); err != nil {
		log.Fatalf("Bot stopped: %v", err)
	}
}
