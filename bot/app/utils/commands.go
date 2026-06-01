package utils

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"

func GetBotCommands() []tgbotapi.BotCommand {
	return []tgbotapi.BotCommand{
		{Command: "start", Description: "Start the bot"},
		{Command: "help", Description: "Get help"},
	}
}
