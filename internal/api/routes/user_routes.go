package routes

import (
	"edjr-trk/internal/api/middlewares/auth"
	"edjr-trk/internal/api/middlewares/dto_validator/validate_user"
	"edjr-trk/internal/ioc"
	"github.com/gofiber/fiber/v2"
)

// RegisterUserRoutes - регистрирует маршруты для работы с user
func RegisterUserRoutes(app *fiber.App, container *ioc.Container) {
	app.Post("/users",
		auth.BasicAuthMiddleware(),
		validate_user.ValidateCreateUserMiddleware(container.Logger),
		container.UserHandler.CreateUser,
	)

	app.Delete("/users/:id",
		auth.BasicAuthMiddleware(),
		container.UserHandler.RemoveUserById,
	)

	app.Get("/users",
		auth.BasicAuthMiddleware(),
		container.UserHandler.GetAllUsers,
	)
}
