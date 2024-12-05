// internal/api/routes/routes.go
package routes

import (
	"edjr-trk/internal/api/validators"
	"edjr-trk/internal/ioc"
	"github.com/gofiber/fiber/v2"
)

// RegisterArticleRoutes - регистрирует маршруты для работы со статьями
func RegisterArticleRoutes(app *fiber.App, container *ioc.Container) {
	app.Post("/articles",
		validators.ValidateCreateArticleMiddleware(container.Logger),
		container.ArticleHandler.CreateArticle,
	)

	app.Patch("/articles/:id",
		validators.ValidateArticleIDMiddleware(container.Logger),
		validators.ValidatePatchArticleMiddleware(container.Logger),
		container.ArticleHandler.PatchArticleById,
	)

	app.Delete("/articles/:id",
		validators.ValidateArticleIDMiddleware(container.Logger),
		container.ArticleHandler.RemoveArticleById,
	)

	app.Get("/articles/:id",
		validators.ValidateArticleIDMiddleware(container.Logger),
		container.ArticleHandler.GetArticleById,
	)

	app.Get("/articles",
		//middlewares.AuthMiddleware,
		validators.ValidatePaginationMiddleware(container.Logger),
		container.ArticleHandler.GetAllArticles,
	)
}
