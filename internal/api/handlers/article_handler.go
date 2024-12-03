package handlers

import (
	"edjr-trk/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

type ArticleHandler struct {
	service *service.ArticleService
	logger  *zap.Logger
}

var validate *validator.Validate

func init() {
	validate = validator.New()
}

func ValidateStruct(input interface{}) error {
	return validate.Struct(input)
}

// NewArticleHandler - создаёт новый экземпляр ArticleHandler.
func NewArticleHandler(service *service.ArticleService, logger *zap.Logger) *ArticleHandler {
	return &ArticleHandler{
		service: service,
		logger:  logger,
	}
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
