package middlewares

import (
	"time"

	"github.com/gofiber/fiber/v2"
	"go.uber.org/zap"
)

// RequestLoggerMiddleware logs request processing time
func RequestLoggerMiddleware(logger *zap.Logger) func(c *fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		start := time.Now()
		err := c.Next()
		duration := time.Since(start)

		// Log the request details
		logger.Info("Request processed",
			zap.String("method", c.Method()),
			zap.String("path", c.Path()),
			zap.Duration("duration", duration))
		return err
	}
}
