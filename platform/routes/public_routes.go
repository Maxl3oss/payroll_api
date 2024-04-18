package routes

import (
	"maxl3oss/app/controllers"
	"maxl3oss/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(route fiber.Router, db *gorm.DB) {
	// News controller
	authController := controllers.NewAuthController(db)
	salaryController := controllers.NewSalaryController(db)
	userController := controllers.NewUserController(db)

	// Route group auth:
	authRoute := route.Group("/auth")

	authRoute.Post("/login", authController.Login)
	// authRoute.Post("/register", authController.Register)
	// authRoute.Post("/refreshToken", authController.RefreshToken)

	// New controller use JWT
	salaryRoute := route.Group("/salary")
	salaryRoute.Use(middleware.JWTProtected())

	salaryRoute.Get("/get-by-user/:id", salaryController.GetByUser)
	salaryRoute.Get("/get", salaryController.GetAll)

	// Route group auth:
	userRoute := route.Group("/user")
	userRoute.Use(middleware.JWTProtected())

	userRoute.Get("/profile", userController.GetProfile)
	userRoute.Patch("/update/:id", userController.Update)
	userRoute.Patch("/change-pass/:id", userController.ChangePassByUser)
}
