package routes

import (
	"edjr-trk/internal/api/middlewares"
	"edjr-trk/internal/ioc"
	"github.com/gofiber/fiber/v2"
)

// RegisterArticleRoutes - регистрирует маршруты для работы со статьями
func RegisterArticleRoutes(app *fiber.App, container *ioc.Container) {
	app.Post("/articles", middlewares.ValidateCreateArticleMiddleware(container.Logger), container.ArticleHandler.CreateArticle)
	app.Patch("/articles/:id", container.ArticleHandler.PatchArticleById)
	app.Delete("/articles/:id", container.ArticleHandler.RemoveArticleById)
	app.Get("/articles/:id", container.ArticleHandler.GetArticleById)
	app.Get("/articles", middlewares.AuthMiddleware, container.ArticleHandler.GetAllArticles)
}
