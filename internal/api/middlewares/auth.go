package middlewares

import (
	"edjr-trk/pkg/http_error"
	"github.com/gofiber/fiber/v2"
)

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	if token != "Bearer my-secret-token" {
		return http_error.NewHTTPError(fiber.StatusUnauthorized, "Unauthorized", nil).Send(c)
	}
	return c.Next() // Переход к следующему обработчику
}
