//go:build dev

package main

import "github.com/joho/godotenv"

func fillEnv() {
	_ = godotenv.Load(".env")
}
