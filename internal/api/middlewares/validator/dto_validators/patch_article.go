package dto_validators

import (
	"edjr-trk/internal/api/dto"
	"edjr-trk/internal/api/middlewares/validator/format_validation_error"
	"edjr-trk/pkg/http_error"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ValidatePatchArticleMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.PatchArticleRequest
		if err := c.BodyParser(&req); err != nil {
			logger.Error("Failed to parse request body", zap.Error(err))
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid request body", nil).Send(c)
		}

		// validate the input data.
		if err := validate.Struct(&req); err != nil {
			logger.Error("Validation failed for request body", zap.Error(err))

			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errorDetails := format_validation_error.FormatValidationErrors(validationErrors)
				return http_error.NewHTTPError(fiber.StatusBadRequest, "Validation error", errorDetails).Send(c)
			}

			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid input", nil).Send(c)
		}

		// Store the validated data in context for use in the handler.
		c.Locals("validatedBody", req)

		return c.Next()
	}
}