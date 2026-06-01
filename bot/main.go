package main

import (
	"log"

	"wallet_bot/app"
	"wallet_bot/app/config"
)

func main() {
	fillEnv()

	cfg, err := config.Load()
	{
		if err != nil {
			log.Fatal("Failed to load config:", err)
		}
	}

	app.Exec(cfg)
}
