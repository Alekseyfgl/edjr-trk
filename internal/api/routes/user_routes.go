package routes

import (
	"edjr-trk/internal/api/middlewares/auth"
	"edjr-trk/internal/api/middlewares/validator/dto_validator"
	"edjr-trk/internal/ioc"
	"github.com/gofiber/fiber/v2"
)

// RegisterUserRoutes - регистрирует маршруты для работы с user
func RegisterUserRoutes(app fiber.Router, container *ioc.Container) {
	app.Post("/users",
		auth.BasicAuthMiddleware(),
		dto_validator.ValidateCreateUserMiddleware(container.Logger),
		container.UserHandler.CreateUser,
	)

	app.Delete("/users/:id",
		auth.BasicAuthMiddleware(),
		dto_validator.ValidateUserIdMiddleware(container.Logger),
		container.UserHandler.RemoveUserById,
	)

	app.Get("/users",
		auth.BasicAuthMiddleware(),
		dto_validator.ValidatePaginationMiddleware(container.Logger),
		container.UserHandler.GetAllUsers,
	)
}
