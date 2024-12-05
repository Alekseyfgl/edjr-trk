package handlers

import (
	"edjr-trk/internal/api/dto"
	"edjr-trk/internal/service"
	"edjr-trk/pkg/http_error"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ArticleHandler struct {
	service *service.ArticleService
	logger  *zap.Logger
}

// NewArticleHandler creates a new instance of ArticleHandler.
func NewArticleHandler(service *service.ArticleService, logger *zap.Logger) *ArticleHandler {
	return &ArticleHandler{
		service: service,
		logger:  logger,
	}
}

// CreateArticle handles creating a new article.
func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	h.logger.Info("Received request to create a new article")

	// Retrieve validated data from context.
	reqInterface := c.Locals("validatedBody")
	req, ok := reqInterface.(dto.CreateArticleRequest)
	if !ok {
		h.logger.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	// Create a new article via the service.
	article, err := h.service.CreateArticle(c.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create article", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to create article", nil).Send(c)
	}

	return c.Status(fiber.StatusCreated).JSON(article)
}

// GetAllArticles handles fetching all articles with pagination.
func (h *ArticleHandler) GetAllArticles(c *fiber.Ctx) error {
	h.logger.Info("Received request to fetch paginated articles")

	// Retrieve pagination parameters from context.
	pageNumberInterface := c.Locals("pageNumber")
	pageSizeInterface := c.Locals("pageSize")

	pageNumber, ok := pageNumberInterface.(int)
	if !ok {
		h.logger.Error("Failed to retrieve page number from context")
		pageNumber = 1
	}

	pageSize, ok := pageSizeInterface.(int)
	if !ok {
		h.logger.Error("Failed to retrieve page size from context")
		pageSize = 10
	}

	// Fetch articles via the service.
	articles, err := h.service.GetAllArticles(c.Context(), pageNumber, pageSize)
	if err != nil {
		h.logger.Error("Failed to fetch paginated articles", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to fetch articles", nil).Send(c)
	}

	h.logger.Info("Paginated articles fetched successfully",
		zap.Int("pageNumber", pageNumber),
		zap.Int("pageSize", pageSize),
		zap.Int("fetchedItems", len(articles.Items)),
	)

	return c.Status(fiber.StatusOK).JSON(articles)
}

// GetArticleById handles fetching a single article by its ID.
func (h *ArticleHandler) GetArticleById(c *fiber.Ctx) error {
	h.logger.Info("Received request to fetch an article by ID")

	// Retrieve article ID from context.
	articleIDInterface := c.Locals("articleID")
	articleID, ok := articleIDInterface.(string)
	if !ok || articleID == "" {
		h.logger.Error("Article ID is missing in the context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Article ID is required", nil).Send(c)
	}

	// Fetch the article via the service.
	article, err := h.service.GetArticleById(c.Context(), articleID)
	if err != nil {
		h.logger.Error("Failed to fetch article by ID", zap.String("articleID", articleID), zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to fetch article", nil).Send(c)
	}

	if article == nil {
		return http_error.NewHTTPError(fiber.StatusNotFound, "Article not found", nil).Send(c)
	}

	h.logger.Info("Article fetched successfully", zap.String("articleID", articleID))
	return c.Status(fiber.StatusOK).JSON(article)
}

// RemoveArticleById handles removing an article by its ID.
func (h *ArticleHandler) RemoveArticleById(c *fiber.Ctx) error {
	h.logger.Info("Received request to remove an article")

	// Retrieve article ID from context.
	articleIDInterface := c.Locals("articleID")
	articleID, ok := articleIDInterface.(string)
	if !ok || articleID == "" {
		h.logger.Error("Article ID is missing in the context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Article ID is required", nil).Send(c)
	}

	// Remove the article via the service.
	removedArticle, err := h.service.RemoveArticleById(c.Context(), articleID)
	if err != nil {
		h.logger.Error("Failed to remove article", zap.String("articleID", articleID), zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to remove article", nil).Send(c)
	}

	h.logger.Info("Article removed successfully", zap.String("articleID", articleID))
	return c.Status(fiber.StatusOK).JSON(fiber.Map{"id": removedArticle})
}

// PatchArticleById handles partial updates of an article.
func (h *ArticleHandler) PatchArticleById(c *fiber.Ctx) error {
	h.logger.Info("Received request to patch an article")

	// Retrieve article ID from context.
	articleIDInterface := c.Locals("articleID")
	articleID, ok := articleIDInterface.(string)
	if !ok || articleID == "" {
		h.logger.Error("Article ID is missing in the context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Article ID is required", nil).Send(c)
	}

	// Retrieve validated data from context.
	reqInterface := c.Locals("validatedBody")
	req, ok := reqInterface.(dto.PatchArticleRequest)
	if !ok {
		h.logger.Error("Failed to retrieve validated request from context")
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Internal Server Error", nil).Send(c)
	}

	// Update the article via the service.
	updatedArticle, err := h.service.PatchArticleById(c.Context(), req, articleID)
	if err != nil {
		h.logger.Error("Failed to update article", zap.String("articleID", articleID), zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusInternalServerError, "Failed to update article", nil).Send(c)
	}

	h.logger.Info("Article updated successfully", zap.String("articleID", articleID))
	return c.Status(fiber.StatusOK).JSON(updatedArticle)
}
