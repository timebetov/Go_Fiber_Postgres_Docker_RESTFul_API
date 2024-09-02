package routes

import (
	fiber "github.com/gofiber/fiber/v2"
	"github.com/timebetov/readerblog/internals/controllers"
	"github.com/timebetov/readerblog/internals/middlewares"
	"github.com/timebetov/readerblog/internals/services"
)

func SetupAuthRoutes(api fiber.Router, authService *services.AuthService, authController *controllers.AuthController) {
	api.Post("/register", authController.RegisterUser)
	api.Post("/login", authController.Login)
	api.Post("/logout", middlewares.AuthenticationMiddleware(authService), authController.Logout)
	api.Get("/profile", middlewares.AuthenticationMiddleware(authService), authController.Profile)
}
