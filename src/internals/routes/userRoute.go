package routes

import (
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/timebetov/readerblog/database"
	"github.com/timebetov/readerblog/internals/controllers"
	"github.com/timebetov/readerblog/internals/middlewares"
	"github.com/timebetov/readerblog/internals/repositories"
)

func SetupUserRoutes(router fiber.Router) {
	user := router.Group("/users")
	user.Use(middlewares.AuthMiddleware(os.Getenv("ADMIN_ROLE")))

	// Initializing repository and controller
	var userRepo = repositories.NewUserRepository(database.DB)
	var userController = controllers.NewUserController(userRepo)

	// Setting up user routes
	user.Post("/", userController.CreateUser)
	user.Get("/", userController.GetUsers)
	user.Get("/deleted", userController.GetDeletedUsers)
	user.Get("/:userId", userController.GetUser)
	user.Patch("/:userId", userController.UpdateUser)
	user.Delete("/:userId", userController.SoftDeleteUser)
	user.Delete("/:userId/force", userController.ForceDeleteUser)
	user.Put("/:userId/restore", userController.RestoreUser)
}
