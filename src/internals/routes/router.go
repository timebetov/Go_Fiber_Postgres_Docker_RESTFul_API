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

	// Initializing repositories and controllers
	userRepo := repositories.NewUserRepository(database.DB)
	authRepo := repositories.NewAuthRepository(database.DB)

	// AUTH
	authService := services.NewAuthService(authRepo)
	authController := controllers.NewAuthController(authService)

	// USER
	userService := services.NewUserService(userRepo)
	userController := controllers.NewUserController(userService)

	// Authentication routes
	api.Post("/register", userController.CreateUser)
	api.Post("/login", authController.Login)
	api.Get("/profile", middlewares.AuthMiddleware(os.Getenv("WRITER_ROLE")), authController.Profile)

	// Setting up user routes only 'admins' can access
	SetupUserRoutes(api)
}
