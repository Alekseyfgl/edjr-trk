package routes

import (
	"edjr-trk/internal/api/middlewares/auth"
	"edjr-trk/internal/api/middlewares/validator/dto_validator"
	"edjr-trk/internal/ioc"
	"github.com/gofiber/fiber/v2"
)

// RegisterArticleRoutes - регистрирует маршруты для работы со статьями
func RegisterArticleRoutes(app *fiber.App, container *ioc.Container) {
	app.Post("/articles",
		auth.JwtAuthMiddleware(container.JwtService),
		dto_validator.ValidateCreateArticleMiddleware(container.Logger),
		container.ArticleHandler.CreateArticle,
	)

	app.Patch("/articles/:id",
		dto_validator.ValidateArticleIdMiddleware(container.Logger),
		dto_validator.ValidatePatchArticleMiddleware(container.Logger),
		container.ArticleHandler.PatchArticleById,
	)

	app.Delete("/articles/:id",
		dto_validator.ValidateArticleIdMiddleware(container.Logger),
		container.ArticleHandler.RemoveArticleById,
	)

	app.Get("/articles/:id",
		dto_validator.ValidateArticleIdMiddleware(container.Logger),
		container.ArticleHandler.GetArticleById,
	)

	app.Get("/articles",
		dto_validator.ValidatePaginationMiddleware(container.Logger),
		container.ArticleHandler.GetAllArticles,
	)
}
