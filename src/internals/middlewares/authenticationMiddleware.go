package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/timebetov/readerblog/internals/services"
	"github.com/timebetov/readerblog/internals/utils"
)

func AuthenticationMiddleware(authService *services.AuthService) fiber.Handler {
	return func(c *fiber.Ctx) error {
		authHeader := c.Get("Authorization")
		if authHeader == "" {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Missing or invalid token",
			})
		}

		tokenStr := strings.Split(authHeader, " ")[1]
		claims, err := utils.ParseToken(tokenStr)
		if err != nil {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Invalid or expired token",
			})
		}

		// Check if the token is blacklisted
		if authService.IsTokenBlacklisted(tokenStr) {
			return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{
				"status":  "error",
				"message": "Token has been blacklisted",
			})
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}
