package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"os"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type User struct {
	ID        int64  `json:"id"`
	Username  string `json:"username"`
	FirstName string `json:"first_name"`
}

const (
	fileWithUsername    = "users_username.json"
	fileWithoutUsername = "users_no_username.json"
)

func load(file string) []User {
	data, err := os.ReadFile(file)
	if err != nil {
		return []User{}
	}

	var users []User
	_ = json.Unmarshal(data, &users)
	return users
}

func save(file string, users []User) {
	data, _ := json.MarshalIndent(users, "", "  ")
	_ = os.WriteFile(file, data, 0644)
}

func addUser(u User) {
	if u.Username != "" {
		users := load(fileWithUsername)
		if !exists(users, u.ID) {
			users = append(users, u)
			save(fileWithUsername, users)
			log.Printf("new user: id=%d @%s (%s)", u.ID, u.Username, u.FirstName)
		}
	} else {
		users := load(fileWithoutUsername)
		if !exists(users, u.ID) {
			users = append(users, u)
			save(fileWithoutUsername, users)
			log.Printf("new user: id=%d first_name=%s", u.ID, u.FirstName)
		}
	}
}

func exists(users []User, id int64) bool {
	for _, u := range users {
		if u.ID == id {
			return true
		}
	}
	return false
}

func mention(u User) string {
	if u.Username != "" {
		return "@" + u.Username
	}
	return fmt.Sprintf(
		`<a href="tg://user?id=%d">%s</a>`,
		u.ID,
		html.EscapeString(u.FirstName),
	)
}

func buildAll() string {
	var result []string

	for _, u := range load(fileWithUsername) {
		result = append(result, mention(u))
	}

	for _, u := range load(fileWithoutUsername) {
		result = append(result, mention(u))
	}

	return strings.Join(result, " ")
}

func splitAndSend(bot *tgbotapi.BotAPI, chatID int64, text string) {
	runes := []rune(text)

	for len(runes) > 0 {
		size := min(3500, len(runes))

		part := string(runes[:size])
		runes = runes[size:]

		msg := tgbotapi.NewMessage(chatID, part)
		msg.ParseMode = "HTML"
		_, err := bot.Send(msg)
		if err != nil {
			log.Println(err)
		}
	}
}

func main() {
	bot, err := tgbotapi.NewBotAPI("7912742673:AAHheL_JApcM65YAqy-vUlWc9NE9rpX6bqY")
	if err != nil {
		panic(err)
	}
	log.Printf("Authorized on account %s", bot.Self.UserName)
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	updates := bot.GetUpdatesChan(u)

	for update := range updates {
		if update.Message == nil {
			continue
		}

		if len(update.Message.NewChatMembers) > 0 {
			for _, member := range update.Message.NewChatMembers {
				addUser(User{
					ID:        int64(member.ID),
					Username:  member.UserName,
					FirstName: member.FirstName,
				})
			}
			continue
		}

		if update.Message.From == nil {
			continue
		}

		user := update.Message.From
		addUser(User{
			ID:        int64(user.ID),
			Username:  user.UserName,
			FirstName: user.FirstName,
		})

		// команда /all
		if update.Message.Text == "/all" {
			text := buildAll()

			splitAndSend(bot, update.Message.Chat.ID, text)
		}
	}
}
