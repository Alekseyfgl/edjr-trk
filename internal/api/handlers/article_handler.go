package handlers

import (
	"edjr-trk/internal/service"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
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

// GetAllArticles - обработчик для получения всех статей с пагинацией.
func (h *ArticleHandler) GetAllArticles(c *fiber.Ctx) error {
	h.logger.Info("Received request to fetch paginated articles")

	// Получение номера страницы
	pageStr := c.Query("page", "1")
	pageNumber, err := strconv.Atoi(pageStr)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	// Получение размера страницы
	sizeStr := c.Query("size", "10")
	pageSize, err := strconv.Atoi(sizeStr)
	if err != nil || pageSize < 1 {
		pageSize = 10
	}

	articles, err := h.service.GetAllArticles(c.Context(), pageNumber, pageSize)
	if err != nil {
		h.logger.Error("Failed to fetch paginated articles", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to fetch articles",
		})
	}

	h.logger.Info("Paginated articles fetched successfully",
		zap.Int("pageNumber", pageNumber),
		zap.Int("pageSize", pageSize),
		zap.Int("fetchedItems", len(articles.Items)),
	)

	return c.Status(fiber.StatusOK).JSON(articles)
}
