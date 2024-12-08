package routes

import (
	"edjr-trk/internal/api/middlewares/auth"
	"edjr-trk/internal/api/middlewares/dto_validator/validate_article"
	"edjr-trk/internal/api/middlewares/param_validator"
	"edjr-trk/internal/ioc"
	"github.com/gofiber/fiber/v2"
)

// RegisterArticleRoutes - регистрирует маршруты для работы со статьями
func RegisterArticleRoutes(app *fiber.App, container *ioc.Container) {
	app.Post("/articles",
		validate_article.ValidateCreateArticleMiddleware(container.Logger),
		container.ArticleHandler.CreateArticle,
	)

	app.Patch("/articles/:id",
		param_validator.ValidateArticleIDMiddleware(container.Logger),
		validate_article.ValidatePatchArticleMiddleware(container.Logger),
		container.ArticleHandler.PatchArticleById,
	)

	app.Delete("/articles/:id",
		param_validator.ValidateArticleIDMiddleware(container.Logger),
		container.ArticleHandler.RemoveArticleById,
	)

	app.Get("/articles/:id",
		param_validator.ValidateArticleIDMiddleware(container.Logger),
		container.ArticleHandler.GetArticleById,
	)

	app.Get("/articles",
		//middlewares.AuthMiddleware,
		auth.BasicAuthMiddleware(),
		validate_article.ValidatePaginationMiddleware(container.Logger),
		container.ArticleHandler.GetAllArticles,
	)
}
