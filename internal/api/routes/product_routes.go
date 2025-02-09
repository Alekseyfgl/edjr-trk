package routes

import (
	"edjr-trk/internal/api/middlewares/auth"
	"edjr-trk/internal/api/middlewares/validator/dto_validator"
	"edjr-trk/internal/ioc"
	"github.com/gofiber/fiber/v2"
)

func RegisterProductRoutes(app fiber.Router, container *ioc.Container) {
	app.Post("/products",
		auth.JwtAuthMiddleware(container.JwtService),
		dto_validator.ValidateCreateProductMiddleware(container.Logger),
		container.ProductHandler.CreateProduct,
	)

	app.Patch("/products/:id",
		auth.JwtAuthMiddleware(container.JwtService),
		dto_validator.ValidateProductIdMiddleware(container.Logger),
		dto_validator.ValidatePatchProductMiddleware(container.Logger),
		container.ProductHandler.PatchProductById,
	)

	app.Delete("/products/:id",
		auth.JwtAuthMiddleware(container.JwtService),
		dto_validator.ValidateProductIdMiddleware(container.Logger),
		container.ProductHandler.RemoveProductById,
	)

	app.Get("/products/:id",
		dto_validator.ValidateProductIdMiddleware(container.Logger),
		container.ProductHandler.GetProductById,
	)

	app.Get("/products",
		dto_validator.ValidatePaginationMiddleware(container.Logger),
		container.ProductHandler.GetAllProducts,
	)
}
