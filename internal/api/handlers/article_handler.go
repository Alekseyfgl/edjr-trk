package handlers

import (
	"edjr-trk/internal/api/dto"
	"edjr-trk/internal/api/validators" // Подключение кастомного валидатора
	"edjr-trk/internal/service"
	"edjr-trk/pkg/http_error"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"strconv"
)

type ArticleHandler struct {
	service  *service.ArticleService
	logger   *zap.Logger
	validate *validator.Validate
}

// NewArticleHandler creates a new instance of ArticleHandler.
func NewArticleHandler(service *service.ArticleService, logger *zap.Logger) *ArticleHandler {
	validate := validator.New()
	validators.RegArticleValidators(validate) // Регистрируем кастомные валидаторы
	return &ArticleHandler{
		service:  service,
		logger:   logger,
		validate: validate,
	}
}

// ValidateStruct validates the given input structure.
func (h *ArticleHandler) ValidateStruct(input interface{}) error {
	return h.validate.Struct(input)
}

// CreateArticle handles creating a new article.
func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	h.logger.Info("Received request to create a new article")

	// Parse the input data.
	var req dto.CreateArticleRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body", zap.Error(err))
		return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid request body", nil).Send(c)
	}

	// Validate the input data.
	if err := h.ValidateStruct(&req); err != nil {
		h.logger.Error("Validation failed for request body", zap.Error(err))

		// If the error is a validation error.
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorDetails := h.formatValidationErrors(validationErrors)
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Validation error", errorDetails).Send(c)
		}

		// For other validation errors.
		return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid input", nil).Send(c)
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

	// Parse the page number.
	pageStr := c.Query("page", "1")
	pageNumber, err := strconv.Atoi(pageStr)
	if err != nil || pageNumber < 1 {
		pageNumber = 1
	}

	// Parse the page size.
	sizeStr := c.Query("size", "10")
	pageSize, err := strconv.Atoi(sizeStr)
	if err != nil || pageSize < 1 {
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

// formatValidationErrors formats validation errors into an array of ErrorItem.
func (h *ArticleHandler) formatValidationErrors(validationErrors validator.ValidationErrors) []http_error.ErrorItem {
	var errorDetails []http_error.ErrorItem
	for _, validationErr := range validationErrors {
		field := validationErr.Field() // The field that failed validation
		tag := validationErr.Tag()     // The validation tag that failed
		var errorMsg string
		switch tag {
		case "required":
			errorMsg = "This field is required"
		case "min":
			errorMsg = "The field does not meet the minimum length requirement"
		case "img_base64_or_null":
			errorMsg = "The field must be null or a valid Base64 string"
		default:
			errorMsg = "Validation failed on tag: " + tag
		}
		errorDetails = append(errorDetails, http_error.ErrorItem{
			Field: field,
			Error: errorMsg,
		})
	}
	return errorDetails
}
