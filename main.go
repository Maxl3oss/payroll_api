package main

import (
	"maxl3oss/pkg/configs"
	"maxl3oss/pkg/middleware"
	"maxl3oss/pkg/utils"
	"maxl3oss/platform/database"
	"maxl3oss/platform/routes"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/limiter"
)

func main() {
	// Load environment variables from .env file
	// err := godotenv.Load()
	// if err != nil {
	// 	panic("Error loading .env file")
	// }

	// Define Fiber config.
	config := configs.FiberConfig()
	// Limit of 20MB
	config.BodyLimit = 20 * 1024 * 1024
	// Define a new Fiber app with config.
	app := fiber.New(config)
	// Set up rate limiter middleware
	app.Use(limiter.New(limiter.Config{
		Max:        100,              // Maximum number of requests per period
		Expiration: 30 * time.Second, // Duration of the rate limit window
		KeyGenerator: func(c *fiber.Ctx) string {
			return c.IP() // Limit by IP address
		},
		LimitReached: func(c *fiber.Ctx) error {
			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"error": "Too many requests",
			})
		},
	}))

	// static file
	app.Static("/uploads", "./uploads")

	// Middlewares.
	middleware.FiberMiddleware(app) // Register Fiber's middleware for app.

	// connect db
	db, err := database.PostgreSQLConnection()
	if err != nil {
		panic(err)
	}

	// Add routes
	router := app.Group("/api/v1")
	routes.PublicRoutes(router, db)
	routes.PrivateRoutes(router, db)

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
