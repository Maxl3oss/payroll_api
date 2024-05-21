package main

import (
	"maxl3oss/pkg/configs"
	"maxl3oss/pkg/middleware"
	"maxl3oss/pkg/utils"
	"maxl3oss/platform/database"
	"maxl3oss/platform/routes"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables from .env file
	err := godotenv.Load()
	if err != nil {
		panic("Error loading .env file")
	}

	// Define Fiber config.
	config := configs.FiberConfig()
	// Limit of 10MB
	config.BodyLimit = 10 * 1024 * 1024
	// Define a new Fiber app with config.
	app := fiber.New(config)

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
