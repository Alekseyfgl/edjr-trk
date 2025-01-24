package routes

import (
	"edjr-trk/internal/api/middlewares/validator/dto_validator"
	"edjr-trk/internal/ioc"
	"github.com/gofiber/fiber/v2"
)

func RegisterEmailRoutes(app fiber.Router, container *ioc.Container) {
	app.Post("/email",
		dto_validator.RateLimiterMiddleware(container.Logger, container.RateLimitService),
		dto_validator.ValidateSendEmailMiddleware(container.Logger),
		container.EmailHandler.SendMsg,
	)
}
