package routes

import (
	"edjr-trk/internal/api/middlewares/dto_validator"
	"edjr-trk/internal/api/middlewares/param_validator"
	"edjr-trk/internal/ioc"
	"github.com/gofiber/fiber/v2"
)

// RegisterArticleRoutes - регистрирует маршруты для работы со статьями
func RegisterArticleRoutes(app *fiber.App, container *ioc.Container) {
	app.Post("/articles",
		dto_validator.ValidateCreateArticleMiddleware(container.Logger),
		container.ArticleHandler.CreateArticle,
	)

	app.Patch("/articles/:id",
		param_validator.ValidateArticleIDMiddleware(container.Logger),
		dto_validator.ValidatePatchArticleMiddleware(container.Logger),
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
		dto_validator.ValidatePaginationMiddleware(container.Logger),
		container.ArticleHandler.GetAllArticles,
	)
}
