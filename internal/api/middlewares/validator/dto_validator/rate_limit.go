package dto_validator

import (
	"edjr-trk/internal/service"
	"edjr-trk/pkg/http_error"
	"edjr-trk/pkg/utils"
	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

func RateLimiterMiddleware(logger *zap.Logger, rateLimiter *service.RateLimiter) fiber.Handler {
	return func(c *fiber.Ctx) error {
		ip := utils.GetClientIP(c)
		// Проверяем запрос на rate limit.
		if err := rateLimiter.ValidateRequest(c); err != nil {
			logger.Warn("Rate limit exceeded", zap.String("ip", ip), zap.Error(err))
			return http_error.NewHTTPError(fiber.StatusTooManyRequests, "Too many requests", nil).Send(c)
		}
		return c.Next()
	}
}
