package handlers

import (
	"edjr-trk/internal/api/dto"
	"edjr-trk/internal/service"
	"edjr-trk/pkg/http_error"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type emailHandler struct {
	service service.EmailServiceInterface
	logger  *zap.Logger
}

type EmailHandlerInterface interface {
	SendMsg(c *fiber.Ctx) error
}

func NewEmailHandler(service service.EmailServiceInterface, logger *zap.Logger) EmailHandlerInterface {
	return &emailHandler{
		service: service,
		logger:  logger,
	}
}

func (h *emailHandler) SendMsg(c *fiber.Ctx) error {
	h.logger.Info("Received request to send email")

	reqInterface := c.Locals("validatedBody")
	body, ok := reqInterface.(dto.SendEmailRequest)

	if !ok {
		h.logger.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	err := h.service.SendMessage(&body)
	if err != nil {
		h.logger.Error("Failed to send email", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to send email", nil).Send(c)
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"status": true})
}
