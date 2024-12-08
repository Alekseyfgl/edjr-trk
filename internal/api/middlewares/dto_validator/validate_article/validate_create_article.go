package validate_article

import (
	"edjr-trk/internal/api/dto"
	"edjr-trk/internal/api/middlewares/dto_validator/validate_common_format"
	"edjr-trk/pkg/http_error"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// ValidateCreateArticleMiddleware validates the create article request.
func ValidateCreateArticleMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse the input data.
		var req dto.CreateArticleRequest
		if err := c.BodyParser(&req); err != nil {
			logger.Error("Failed to parse request body", zap.Error(err))
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid request body", nil).Send(c)
		}

		// Create a validator and register custom validators.
		validate := validator.New()
		RegArticleValidators(validate)

		// Validate the input data.
		if err := validate.Struct(&req); err != nil {
			logger.Error("Validation failed for request body", zap.Error(err))

			// If the error is a validation error.
			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errorDetails := validate_common_format.FormatValidationErrors(validationErrors)
				return http_error.NewHTTPError(fiber.StatusBadRequest, "Validation error", errorDetails).Send(c)
			}

			// For other validation errors.
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid input", nil).Send(c)
		}

		// Store the validated data in context for use in the handler.
		c.Locals("validatedBody", req)

		// Proceed to the next handler.
		return c.Next()
	}
}

//// formatValidationErrors formats validation errors into an array of ErrorItem.
//func formatValidationErrors(validationErrors validator.ValidationErrors) []http_error.ErrorItem {
//	var errorDetails []http_error.ErrorItem
//	for _, validationErr := range validationErrors {
//		field := validationErr.Field()
//		tag := validationErr.Tag()
//		var errorMsg string
//		switch tag {
//		case "required":
//			errorMsg = "This field is required"
//		case "min":
//			errorMsg = "The field does not meet the minimum length requirement"
//		case "img_base64_or_null":
//			errorMsg = "The field must be null or a valid Base64 string"
//		default:
//			errorMsg = "Validation failed on tag: " + tag
//		}
//		errorDetails = append(errorDetails, http_error.ErrorItem{
//			Field: field,
//			Error: errorMsg,
//		})
//	}
//	return errorDetails
//}
