package services

import (
	"fmt"
	"net/smtp"
	"os"
	"strconv"
	"strings"
)

type EmailService struct {
	host     string
	port     string
	username string
	password string
	from     string
}

func NewEmailServiceFromEnv() *EmailService {
	return &EmailService{
		host:     envOrDefault("SMTP_HOST", "smtp.gmail.com"),
		port:     envOrDefault("SMTP_PORT", "587"),
		username: os.Getenv("SMTP_USERNAME"),
		password: os.Getenv("SMTP_PASSWORD"),
		from:     envOrDefault("SMTP_FROM", os.Getenv("SMTP_USERNAME")),
	}
}

func (s *EmailService) SendPasswordResetEmail(toEmail, resetLink string) error {
	if s.host == "" || s.port == "" || s.username == "" || s.password == "" || s.from == "" {
		return fmt.Errorf("smtp configuration is incomplete")
	}

	ttl := passwordResetTTLMinutes()
	subject := "Reset your AI Inference Gateway password"
	body := fmt.Sprintf(
		"Use this link to reset your AI Inference Gateway password:\r\n\r\n%s\r\n\r\nThis link expires in %d minutes.\r\n",
		resetLink,
		ttl,
	)

	message := strings.Join([]string{
		"From: " + s.from,
		"To: " + toEmail,
		"Subject: " + subject,
		"MIME-Version: 1.0",
		"Content-Type: text/plain; charset=UTF-8",
		"",
		body,
	}, "\r\n")

	auth := smtp.PlainAuth("", s.username, s.password, s.host)
	return smtp.SendMail(s.host+":"+s.port, auth, s.from, []string{toEmail}, []byte(message))
}

func envOrDefault(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}

func passwordResetTTLMinutes() int {
	value, err := strconv.Atoi(os.Getenv("PASSWORD_RESET_TOKEN_TTL_MINUTES"))
	if err != nil || value <= 0 {
		return 15
	}
	return value
}
