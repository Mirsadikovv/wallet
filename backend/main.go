package main

import (
	"log"

	"github.com/caarlos0/env/v11"
	"github.com/joho/godotenv"
	"wallet/src"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	cfg := &src.Env{}
	{
		if err := env.Parse(cfg); err != nil {
			log.Fatalf("Failed to parse environment variables: %v", err)
		}
	}

	src.Exec(cfg)
}
