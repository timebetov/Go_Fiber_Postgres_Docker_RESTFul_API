package routes

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/timebetov/readerblog/database"
	"github.com/timebetov/readerblog/internals/controllers"
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

	// Initializing redis client
	redisClient := database.NewRedisClient()

	// Initializing repositories
	authRepo := repositories.NewAuthRepository(database.DB)
	userRepo := repositories.NewUserRepository(database.DB)

	// Initializing services
	userService := services.NewUserService(userRepo)
	authService := services.NewAuthService(authRepo, userService, redisClient)

	// Initializing controllers
	authController := controllers.NewAuthController(authService)
	userController := controllers.NewUserController(userService)

	// Authentication routes
	SetupAuthRoutes(api, authService, authController)
	// Setting up user routes only 'admins' can access
	SetupUserRoutes(api, authService, userController)
}
