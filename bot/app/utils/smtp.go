package utils

import (
	"fmt"
	"net/smtp"

	"wallet_bot/app/config"
)

func SendOTP(cfg *config.Config, toEmail, code string) error {
	from := cfg.SmtpFrom
	to := []string{toEmail}

	subject := "Your verification code"
	body := fmt.Sprintf("Your verification code: <b>%s</b><br>Valid for 5 minutes.", code)

	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nMIME-Version: 1.0\r\nContent-Type: text/html; charset=UTF-8\r\n\r\n%s",
		from, toEmail, subject, body,
	)

	auth := smtp.PlainAuth("", cfg.SmtpUsername, cfg.SmtpPassword, cfg.SmtpHost)
	return smtp.SendMail(cfg.SmtpHost+":"+cfg.SmtpPort, auth, from, to, []byte(msg))
}
