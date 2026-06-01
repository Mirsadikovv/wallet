package config

import (
	"os"
	"strconv"
	"strings"
)

type Config struct {
	BotToken         string
	AdminIds         []string
	Debug            bool
	WebAppURL        string
	BackendURL       string
	TONNetwork       string
	SmtpHost         string
	SmtpPort         string
	SmtpUsername     string
	SmtpPassword     string
	SmtpFrom         string
	PostgresHost     string
	PostgresUser     string
	PostgresPassword string
	PostgresDB       string
	PostgresPort     int
	PostgresSSLMode  string
	PostgresTimeZone string
}

func Load() (*Config, error) {
	postgresPort, _ := strconv.Atoi(os.Getenv("POSTGRES_PORT"))
	if postgresPort == 0 {
		postgresPort = 5432
	}

	var admins []string
	if adminEnv := os.Getenv("ADMINS"); adminEnv != "" {
		admins = strings.Split(adminEnv, ",")
	}

	return &Config{
		BotToken:         os.Getenv("BOT_TOKEN"),
		AdminIds:         admins,
		Debug:            os.Getenv("DEBUG") == "true",
		WebAppURL:        os.Getenv("WEB_APP_URL"),
		BackendURL:       os.Getenv("BACKEND_URL"),
		TONNetwork:       getEnvOrDefault("TON_NETWORK", "testnet"),
		SmtpHost:         os.Getenv("SMTP_HOST"),
		SmtpPort:         getEnvOrDefault("SMTP_PORT", "587"),
		SmtpUsername:     os.Getenv("SMTP_USERNAME"),
		SmtpPassword:     os.Getenv("SMTP_PASSWORD"),
		SmtpFrom:         os.Getenv("SMTP_FROM"),
		PostgresHost:     os.Getenv("POSTGRES_HOST"),
		PostgresUser:     os.Getenv("POSTGRES_USER"),
		PostgresPassword: os.Getenv("POSTGRES_PASSWORD"),
		PostgresDB:       os.Getenv("POSTGRES_DB"),
		PostgresPort:     postgresPort,
		PostgresSSLMode:  getEnvOrDefault("POSTGRES_SSL_MODE", "disable"),
		PostgresTimeZone: getEnvOrDefault("POSTGRES_TIME_ZONE", "UTC"),
	}, nil
}

func getEnvOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
