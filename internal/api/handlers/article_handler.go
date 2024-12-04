package handlers

import (
	"edjr-trk/internal/api/dto"
	"edjr-trk/internal/service"
	"encoding/base64"
	"fmt"
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

func validateImg(fl validator.FieldLevel) bool {
	img, ok := fl.Field().Interface().(*string)
	if !ok {
		fmt.Println("Img field is not of type *string")
		return false
	}

	if img == nil {
		fmt.Println("Img field is null - validation passed")
		return true
	}

	fmt.Println("Validating img field:", *img)

	decoded, err := base64.StdEncoding.DecodeString(*img)
	if err != nil {
		fmt.Println("Base64 decode error:", err)
		return false
	}

	if len(decoded) == 0 {
		fmt.Println("Decoded Base64 is empty")
		return false
	}

	return true
}
func init() {
	validate = validator.New()

	// Регистрируем кастомное правило валидации
	err := validate.RegisterValidation("img_base64_or_null", validateImg)
	if err != nil {
		// Логируем ошибку
		zap.L().Fatal("Failed to register custom validation rule", zap.Error(err))
	}
}

func ValidateStruct(input interface{}) error {
	fmt.Printf("Validating struct: %+v\n", input)
	return validate.Struct(input)
}

// NewArticleHandler - создаёт новый экземпляр ArticleHandler.
func NewArticleHandler(service *service.ArticleService, logger *zap.Logger) *ArticleHandler {
	return &ArticleHandler{
		service: service,
		logger:  logger,
	}
}

func (h *ArticleHandler) CreateArticle(c *fiber.Ctx) error {
	h.logger.Info("Received request to create a new article")

	// Парсинг входных данных.
	var req dto.CreateArticleRequest
	if err := c.BodyParser(&req); err != nil {
		h.logger.Error("Failed to parse request body", zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid request body",
		})
	}

	// Валидация данных.
	if err := ValidateStruct(&req); err != nil {
		h.logger.Error("Validation failed for request body", zap.Error(err))
		fmt.Printf("Validation error: %+v\n", err)

		// Если ошибка является ошибкой валидации.
		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			errorMessages := make(map[string]string)
			for _, validationErr := range validationErrors {
				field := validationErr.Field() // Поле, которое не прошло проверку
				tag := validationErr.Tag()     // Тег правила, которое не прошло
				switch tag {
				case "required":
					errorMessages[field] = "This field is required"
				case "min":
					errorMessages[field] = "The field does not meet the minimum length requirement"
				case "img_base64_or_null":
					errorMessages[field] = "The field must be null or a valid Base64 string"
				default:
					errorMessages[field] = "Validation failed on tag: " + tag
				}
			}
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error":   "Validation error",
				"details": errorMessages,
			})
		}

		// Если ошибка не связана с валидацией.
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Invalid input",
		})
	}

	// Создание новой статьи через сервис.
	article, err := h.service.CreateArticle(c.Context(), req)
	if err != nil {
		h.logger.Error("Failed to create article", zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Failed to create article",
		})
	}

	return c.Status(fiber.StatusCreated).JSON(article)
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
