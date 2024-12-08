package dto_validators

import (
	"edjr-trk/pkg/http_error"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ValidateArticleIdMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		articleID := c.Params("id")
		if articleID == "" {
			logger.Error("Article ID is missing in the request")
			return http_error.NewHTTPError(fiber.StatusBadRequest, "Article ID is required", nil).Send(c)
		}

		// Additional validation for article ID can be added here.

		// Store the article ID in context for use in the handler.
		c.Locals("articleID", articleID)

		return c.Next()
	}
}
