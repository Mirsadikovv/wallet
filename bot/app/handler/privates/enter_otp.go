package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"

	"wallet_bot/app/config"
	"wallet_bot/app/utils"
)

type botUser struct {
	ID               int64  `gorm:"primaryKey;autoIncrement"`
	TelegramID       *int64 `gorm:"column:telegram_id"`
	TelegramUsername string `gorm:"column:telegram_username"`
	Username         string `gorm:"column:username"`
	Email            string `gorm:"column:email"`
}

func (botUser) TableName() string { return "users" }

func clauseOnConflict() clause.OnConflict {
	return clause.OnConflict{
		Columns: []clause.Column{{Name: "telegram_id"}},
		DoUpdates: clause.Assignments(map[string]any{
			"telegram_username": gorm.Expr("excluded.telegram_username"),
		}),
	}
}

func EnterOTP(bot *tgbotapi.BotAPI, update tgbotapi.Update, cfg *config.Config, db *gorm.DB, lang *utils.LanguageCache, state *utils.StateCache, otp *utils.OTPCache) *utils.LanguageCache {
	userID := update.Message.From.ID
	userLang := lang.Get(userID)
	code := strings.TrimSpace(update.Message.Text)

	email := strings.TrimPrefix(state.Get(userID), "await_otp:")

	if !otp.Verify(email, code) {
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, OTPInvalidText[userLang])
		msg.ParseMode = "HTML"
		if _, err := bot.Send(msg); err != nil {
			log.Println("Error sending message:", err)
		}
		return lang
	}

	tgID := userID
	username := update.Message.From.UserName
	if username == "" {
		username = fmt.Sprintf("user_%d", userID)
	}

	newUser := botUser{
		TelegramID:       &tgID,
		TelegramUsername: update.Message.From.UserName,
		Username:         username,
		Email:            email,
	}

	if err := db.Clauses(clauseOnConflict()).Create(&newUser).Error; err != nil {
		log.Printf("Failed to create user %d: %v", userID, err)
		msg := tgbotapi.NewMessage(update.Message.Chat.ID, RegistrationErrorText[userLang])
		msg.ParseMode = "HTML"
		if _, err := bot.Send(msg); err != nil {
			log.Println("Error sending message:", err)
		}
		return lang
	}

	if err := createWallet(cfg, newUser.ID); err != nil {
		log.Printf("Failed to create wallet for user %d: %v", newUser.ID, err)
	}

	state.Clear(userID)

	msg := tgbotapi.NewMessage(update.Message.Chat.ID, OTPSuccessText[userLang])
	msg.ParseMode = "HTML"
	if _, err := bot.Send(msg); err != nil {
		log.Println("Error sending message:", err)
	}

	return Menu(bot, update, cfg, lang)
}

func createWallet(cfg *config.Config, userID int64) error {
	body, err := json.Marshal(map[string]any{
		"user_id":     userID,
		"wallet_type": "v4r2",
		"network":     cfg.TONNetwork,
	})
	{
		if err != nil {
			return fmt.Errorf("failed to marshal wallet request: %w", err)
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, cfg.BackendURL+"/api/v1/wallet", bytes.NewReader(body))
	{
		if err != nil {
			return fmt.Errorf("failed to build wallet request: %w", err)
		}
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	{
		if err != nil {
			return fmt.Errorf("wallet API call failed: %w", err)
		}
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return fmt.Errorf("wallet API returned status %d", resp.StatusCode)
	}

	return nil
}
