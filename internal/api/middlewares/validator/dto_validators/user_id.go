package dto_validators

import (
	"edjr-trk/pkg/http_error"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func ValidateUserIdMiddleware(logger *zap.Logger) fiber.Handler {
	return func(c *fiber.Ctx) error {
		userId := c.Params("id")
		if userId == "" {
			logger.Error("User Id is missing in the request")
			return http_error.NewHTTPError(fiber.StatusBadRequest, "User id is required", nil).Send(c)
		}

		// Store the user ID in context for use in the handler.
		c.Locals("userId", userId)

		return c.Next()
	}
}
