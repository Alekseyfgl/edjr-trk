package routes

import (
	"edjr-trk/internal/api/middlewares/validator/dto_validator"
	"edjr-trk/internal/ioc"
	"github.com/gofiber/fiber/v2"
)

func RegisterAuthRoutes(app fiber.Router, container *ioc.Container) {
	app.Post("/auth/login",
		dto_validator.ValidateLoginMiddleware(container.Logger),
		container.AuthHandler.Login,
	)
}
