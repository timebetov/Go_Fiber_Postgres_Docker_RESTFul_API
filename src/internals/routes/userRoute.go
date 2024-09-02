package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timebetov/readerblog/internals/controllers"
	"github.com/timebetov/readerblog/internals/middlewares"
	"github.com/timebetov/readerblog/internals/services"
)

// All routes related to user
func SetupUserRoutes(api fiber.Router, authService *services.AuthService, userController *controllers.UserController) {
	users := api.Group("/users")
	users.Use(middlewares.AuthenticationMiddleware(authService))
	users.Use(middlewares.AuthorizationMiddleware())

	users.Post("/", userController.CreateUser)
	users.Get("/", userController.GetUsers)
	users.Get("/:userId", userController.GetUser)
	users.Patch("/:userId", userController.UpdateUser)
	users.Delete("/:userId", userController.DeleteUser)
	users.Put("/:userId/restore", userController.RestoreUser)
}
