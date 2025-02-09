package dto_validator

import (
	"edjr-trk/pkg/http_error"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ValidateProductIdMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		productID := c.Params("id")
		if productID == "" {
			logger.Error("Product ID is missing in the request")
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Product ID is required", nil).Send(c)
		}
		c.Locals("productID", productID)
		return c.Next()
	}
}
