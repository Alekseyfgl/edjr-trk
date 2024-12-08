package validate_user

import (
	"edjr-trk/internal/api/dto"
	"edjr-trk/internal/api/middlewares/dto_validator/validate_common_format"
	"edjr-trk/pkg/http_error"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
	"regexp"
)

var validate = validator.New()

// init registers custom validation rules.
func init() {
	// Register custom email validation.
	validate.RegisterValidation("custom_email", func(fl validator.FieldLevel) bool {
		emailRegex := `^\w+([-+.']\w+)*@\w+([-.]\w+)*\.\w+([-.]\w+)*$`
		re := regexp.MustCompile(emailRegex)
		return re.MatchString(fl.Field().String())
	})
}

// ValidateCreateUserMiddleware validates the create user request.
func ValidateCreateUserMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		// Parse the input data.
		var req dto.CreateUserRequest
		if err := c.BodyParser(&req); err != nil {
			logger.Error("Failed to parse request body", zap.Error(err))
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid request body", nil).Send(c)
		}

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
