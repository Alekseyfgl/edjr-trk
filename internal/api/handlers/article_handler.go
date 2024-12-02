package handlers

import (
	"edjr-trk/internal/repository"
	"edjr-trk/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ArticleHandler struct {
	service *service.ArticleService
	logger  *zap.Logger
}

// NewArticleHandler - создаёт новый экземпляр ArticleHandler.
func NewArticleHandler(service *service.ArticleService, logger *zap.Logger) *ArticleHandler {
	return &ArticleHandler{
		service: service,
		logger:  logger,
	}
}

// CreateArticle - обработчик для создания новой статьи.
func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	h.logger.Info("Received request to create an article")

	var articleInput repository.Article
	if err := c.BodyParser(&articleInput); err != nil {
		h.logger.Error("Invalid request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	id, err := h.service.CreateArticle(c.Context(), articleInput)
	if err != nil {
		h.logger.Error("Failed to create article", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create article",
		})
	}

	h.logger.Info("Article created successfully", zap.String("id", id.Hex()))
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"id": id.Hex(),
	})
}

// GetArticle - обработчик для получения статьи по ID.
func (h *ArticleHandler) GetArticle(c *fiber.Ctx) error {
	id := c.Params("id")
	h.logger.Info("Received request to get an article", zap.String("id", id))

	article, err := h.service.GetArticleByID(c.Context(), id)
	if err != nil {
		h.logger.Error("Failed to fetch article by ID", zap.String("id", id), zap.Error(err))
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"error": "Article not found",
		})
	}

	h.logger.Info("Article fetched successfully", zap.String("id", id))
	return c.Status(fiber.StatusOK).JSON(article)
}

// GetAllArticles - обработчик для получения всех статей.
func (h *ArticleHandler) GetAllArticles(c *fiber.Ctx) error {
	h.logger.Info("Received request to fetch all articles")

	articles, err := h.service.GetAllArticles(c.Context())
	if err != nil {
		h.logger.Error("Failed to fetch all articles", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch articles",
		})
	}

	h.logger.Info("All articles fetched successfully", zap.Int("count", len(articles)))
	return c.Status(fiber.StatusOK).JSON(articles)
}
