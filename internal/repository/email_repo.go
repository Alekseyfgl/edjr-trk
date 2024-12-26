package repository

import (
	"fmt"
	"go.uber.org/zap"
	"net/smtp"
)

type EmailRepositoryInterface interface {
	SendEmail(from, password, to, subject, body string) error
}

type smtpEmailRepository struct {
	host   string
	port   string
	logger *zap.Logger
}

func NewSMTPEmailRepository(host, port string, logger *zap.Logger) EmailRepositoryInterface {
	return &smtpEmailRepository{host: host, port: port, logger: logger}
}

func (r *smtpEmailRepository) SendEmail(from, password, to, subject, body string) error {
	msg := []byte(fmt.Sprintf(
		"Subject: %s\r\n"+
			"Content-Type: text/html; charset=\"UTF-8\"\r\n"+
			"\r\n"+
			"%s",
		subject, body,
	))

	auth := smtp.PlainAuth("", from, password, r.host)
	addr := fmt.Sprintf("%s:%s", r.host, r.port)

	if err := smtp.SendMail(addr, auth, from, []string{to}, msg); err != nil {
		return err
	}
	return nil
}
