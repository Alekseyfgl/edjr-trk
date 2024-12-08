package routes

import (
	"edjr-trk/internal/api/middlewares/validator/dto_validators"
	"edjr-trk/internal/ioc"
	"github.com/gofiber/fiber/v2"
)

// RegisterArticleRoutes - регистрирует маршруты для работы со статьями
func RegisterArticleRoutes(app *fiber.App, container *ioc.Container) {
	app.Post("/articles",
		dto_validators.ValidateCreateArticleMiddleware(container.Logger),
		container.ArticleHandler.CreateArticle,
	)

	app.Patch("/articles/:id",
		dto_validators.ValidateArticleIdMiddleware(container.Logger),
		dto_validators.ValidatePatchArticleMiddleware(container.Logger),
		container.ArticleHandler.PatchArticleById,
	)

	app.Delete("/articles/:id",
		dto_validators.ValidateArticleIdMiddleware(container.Logger),
		container.ArticleHandler.RemoveArticleById,
	)

	app.Get("/articles/:id",
		dto_validators.ValidateArticleIdMiddleware(container.Logger),
		container.ArticleHandler.GetArticleById,
	)

	app.Get("/articles",
		dto_validators.ValidatePaginationMiddleware(container.Logger),
		container.ArticleHandler.GetAllArticles,
	)
}
