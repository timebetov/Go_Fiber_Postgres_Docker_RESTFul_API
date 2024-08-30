package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/timebetov/readerblog/database"
	"github.com/timebetov/readerblog/internals/routes"
)

// Entrypoint of the application
func main() {
	// Connecting to DB
	database.ConnectDB()

	// Initializng fiber app
	app := fiber.New()

	// Setting up API routes
	routes.SetupRoutes(app)

	// Starting server on port: 3000
	app.Listen(":3000")
}
