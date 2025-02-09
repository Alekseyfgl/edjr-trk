package dto_validator

import (
	"edjr-trk/internal/api/dto"
	"edjr-trk/internal/api/middlewares/validator/format_validation_error"
	"edjr-trk/pkg/http_error"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ValidatePatchProductMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req dto.PatchProductRequest
		if err := c.BodyParser(&req); err != nil {
			logger.Error("Failed to parse request body", zap.Error(err))
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid request body", nil).Send(c)
		}

		if err := validate.Struct(&req); err != nil {
			logger.Error("Validation failed for request body", zap.Error(err))

			if validationErrors, ok := err.(validator.ValidationErrors); ok {
				errorDetails := format_validation_error.FormatValidationErrors(validationErrors)
				return http_error.NewHTTPError(fiber.StatusBadRequest, "Validation error", errorDetails).Send(c)
			}

			return http_error.NewHTTPError(fiber.StatusBadRequest, "Invalid input", nil).Send(c)
		}

		c.Locals("validatedBody", req)

		return c.Next()
	}
}
