package routes

import (
	"edjr-trk/internal/ioc"
	"github.com/gofiber/fiber/v2"
)

// RegisterArticleRoutes - регистрирует маршруты для работы со статьями
func RegisterArticleRoutes(app *fiber.App, container *ioc.Container) {
	app.Post("/articles", container.ArticleHandler.CreateArticle)
	//app.Patch("/articles/:id", container.ArticleHandler.PatchArticle)
	app.Get("/articles/:id", container.ArticleHandler.GetArticleById)
	app.Get("/articles", container.ArticleHandler.GetAllArticles)
}
