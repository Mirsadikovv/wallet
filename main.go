package main

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"wallet_test/src"
)

func main() {
	// Загружаем .env файл (если существует)
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Парсим переменные окружения в структуру
	cfg := &src.Env{}
	if err := env.Parse(cfg); err != nil {
		log.Fatalf("Failed to parse environment variables: %v", err)
	}

	log.Printf("Starting server on port %d", cfg.HTTP_Port)
	log.Printf("TON Network: %s", cfg.TON_Network)

	src.Exec(cfg)
}
