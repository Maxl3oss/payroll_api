package routes

import (
	"maxl3oss/app/controllers"
	"maxl3oss/pkg/middleware"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

// PublicRoutes func for describe group of public routes.
func PrivateRoutes(route fiber.Router, db *gorm.DB) {
	// Use JWT protected middleware
	// route.Use(middleware.JWTProtectedAdmin())

	// News controller
	userController := controllers.NewUserController(db)
	uploadController := controllers.NewUploadController(db)
	salaryController := controllers.NewSalaryController(db)

	// Route group user:
	userRoute := route.Group("/user")
	userRoute.Use(middleware.JWTProtectedAdmin())

	// Add route
	userRoute.Get("", userController.GetAll)
	userRoute.Post("", userController.Create)
	userRoute.Delete("/:id", userController.DeleteById)

	// Route group upload:
	uploadRoute := route.Group("/uploads")
	uploadRoute.Use(middleware.JWTProtectedAdmin())

	uploadRoute.Post("", uploadController.HandleUpload)

	// Route group salary:
	salaryRoute := route.Group("/salary")
	salaryRoute.Use(middleware.JWTProtectedAdmin())

	salaryRoute.Post("/uploads", salaryController.UploadSalaryV2)
	salaryRoute.Get("/get-all", salaryController.GetAll)
	salaryRoute.Delete("/delete-by-month", salaryController.DeleteManySalary)

	salaryRoute.Get("/get-by", salaryController.GetAll)
}
