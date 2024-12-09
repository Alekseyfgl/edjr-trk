package handlers

import (
	"edjr-trk/internal/api/dto"
	"edjr-trk/internal/service"
	"edjr-trk/pkg/http_error"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type authHandler struct {
	service service.AuthServiceInterface
	logger  *zap.Logger
}

type AuthHandlerInterface interface {
	Login(c *fiber.Ctx) error
}

// NewAuthHandler creates a new instance of UserHandler.
func NewAuthHandler(service service.AuthServiceInterface, logger *zap.Logger) AuthHandlerInterface {
	return &authHandler{
		service: service,
		logger:  logger,
	}
}

func (h *authHandler) Login(c *fiber.Ctx) error {
	h.logger.Info("Received request to login the user")

	//Retrieve validated data from context.
	reqInterface := c.Locals("validatedBody")
	body, ok := reqInterface.(dto.LoginRequest)

	if !ok {
		h.logger.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	user, err := h.service.Login(c.Context(), body.Email, body.Password)
	if err != nil {
		h.logger.Error("Failed to create article", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to login user", nil).Send(c)
	}

	return c.Status(fiber.StatusCreated).JSON(user)
}
