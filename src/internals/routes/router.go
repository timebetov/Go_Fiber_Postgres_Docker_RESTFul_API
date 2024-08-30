package routes

import (
	"os"

	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/timebetov/readerblog/database"
	"github.com/timebetov/readerblog/internals/controllers"
	"github.com/timebetov/readerblog/internals/middlewares"
	"github.com/timebetov/readerblog/internals/repositories"
	"github.com/timebetov/readerblog/internals/services"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api", logger.New())

	api.Get("/", func(c *fiber.Ctx) error {
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"status":  "success",
			"message": "API is up and running"})
	})

	// Initializing repository and controller
	var userRepo = repositories.NewUserRepository(database.DB)
	var authService = services.NewAuthService(userRepo)
	var authController = controllers.NewAuthController(authService)

	// Authentication routes
	api.Post("/register", authController.Register)
	api.Post("/login", authController.Login)
	api.Get("/profile", middlewares.AuthMiddleware(os.Getenv("WRITER_ROLE")), authController.Profile)

	// Setting up user routes only 'admins' can access
	SetupUserRoutes(api)
}
