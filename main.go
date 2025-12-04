package main

import (
	"log"
	"os"

	"github.com/AllanKariuki/fiber-go-rest-api/config"
	"github.com/AllanKariuki/fiber-go-rest-api/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/recover"
)

func main() {
	// load environment variables
	config.LoadEnv()

	// Initialize database
	config.ConnectDB()

	// Create Fiber app
	app := fiber.New(fiber.Config{
		AppName:      "Go fiver MVC API v1.0.0",
		ErrorHandler: config.ErrorHandler,
	})

	// Middleware
	app.Use(logger.New())
	app.Use(recover.New())
	app.Use(cors.New())

	routes.SetupRoutes(app)

	// Health check
	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status": "ok",
		})
	})

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Fatal(app.Listen(":" + port))
}
