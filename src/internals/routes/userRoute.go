package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/timebetov/readerblog/database"
	"github.com/timebetov/readerblog/internals/controllers"
	"github.com/timebetov/readerblog/internals/middlewares"
	"github.com/timebetov/readerblog/internals/repositories"
	"github.com/timebetov/readerblog/internals/services"
)

func SetupUserRoutes(router fiber.Router) {
	user := router.Group("/users")
	user.Use(middlewares.AuthMiddleware(os.Getenv("ADMIN_ROLE")))

	// Initializing the user repository and controller
	userRepository := repositories.NewUserRepository(database.DB)
	userService := services.NewUserService(userRepository)
	userController := controllers.NewUserController(userService)

	// Setting up user routes
	user.Post("/", userController.CreateUser)
	user.Get("/", userController.GetUsers)
	user.Get("/:userId", userController.GetUser)
	user.Patch("/:userId", userController.UpdateUser)
	user.Delete("/:userId", userController.DeleteUser)
	user.Put("/:userId/restore", userController.RestoreUser)
}
