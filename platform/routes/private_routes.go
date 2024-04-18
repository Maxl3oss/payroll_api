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
	userRoute.Get("/get-all", userController.GetAll)
	userRoute.Get("/get/:id", userController.GetByUser)
	userRoute.Get("/get-role", userController.GetRole)
	userRoute.Post("/add", userController.Create)
	userRoute.Delete("/del/:id", userController.DeleteById)
	userRoute.Patch("/update-pass/:id", userController.ChangePassByAdmin)

	// Route group upload:
	uploadRoute := route.Group("/uploads")
	uploadRoute.Use(middleware.JWTProtectedAdmin())

	uploadRoute.Post("", uploadController.HandleUpload)

	// Route group salary:
	salaryRoute := route.Group("/salary")
	salaryRoute.Use(middleware.JWTProtectedAdmin())

	salaryRoute.Get("/get-all", salaryController.GetAll)
	salaryRoute.Get("/get-by", salaryController.GetAll)
	salaryRoute.Post("/uploads", salaryController.UploadSalaryV2)
	salaryRoute.Delete("/delete-by-month", salaryController.DeleteManySalary)
}
