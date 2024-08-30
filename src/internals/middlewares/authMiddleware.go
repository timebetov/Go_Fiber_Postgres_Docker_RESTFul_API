package middlewares

import (
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/timebetov/readerblog/internals/utils"
)

func AuthMiddleware(requiredRole string) fiber.Handler {
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

		if claims.Role != requiredRole && claims.Role != "admin" {
			return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
				"status":  "error",
				"message": "insufficient permissions",
			})
		}

		c.Locals("claims", claims)
		return c.Next()
	}
}
