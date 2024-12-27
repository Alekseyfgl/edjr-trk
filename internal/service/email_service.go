package service

import (
	"edjr-trk/configs/env"
	"edjr-trk/internal/api/dto"
	"edjr-trk/internal/repository"
	"fmt"
	"go.uber.org/zap"
)

type EmailServiceInterface interface {
	SendMessage(dto *dto.SendEmailRequest) error
}

type emailService struct {
	repo   repository.EmailRepositoryInterface
	logger *zap.Logger
}

func NewEmailService(repo repository.EmailRepositoryInterface, logger *zap.Logger) EmailServiceInterface {
	return &emailService{repo, logger}
}

func (s *emailService) SendMessage(dto *dto.SendEmailRequest) error {
	from := env.GetEnv("GMAIL_FROM", "")
	password := env.GetEnv("GMAIL_PASSWORD", "")
	to := env.GetEnv("GMAIL_TO", "")
	subject := "Message from your website!"
	body := fmt.Sprintf(
		`<html>
		<body>
			<p><strong>Name:</strong> %s</p>
			<p><strong>Phone:</strong> <a href="tel:%s">%s</a></p>
			<p><strong>Message:</strong></p>
			<p>%s</p>
			<p>Best wishes,<br>Your team.</p>
		</body>
		</html>`,
		dto.Name, dto.Phone, dto.Phone, dto.Text,
	)

	err := s.repo.SendEmail(from, password, to, subject, body)
	if err != nil {
		s.logger.Error("Ошибка при отправке письма: %v", zap.Error(err))
	}
	return nil
}
