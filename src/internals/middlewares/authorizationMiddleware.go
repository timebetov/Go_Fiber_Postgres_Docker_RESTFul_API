package middlewares

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/timebetov/readerblog/internals/utils"
)

func AuthorizationMiddleware(requiredRole ...string) fiber.Handler {
	return func(c *fiber.Ctx) error {
		claims := c.Locals("claims").(*utils.Claims)

		if len(requiredRole) > 0 && requiredRole[0] != "" {
			if claims.Role != requiredRole[0] && claims.Role != os.Getenv("ADMIN_ROLE") {
				return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
					"status":  "error",
					"message": "insufficient permissions",
				})
			}
		}

		return c.Next()
	}
}
