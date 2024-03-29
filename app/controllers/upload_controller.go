package controllers

import (
	"log"
	"maxl3oss/pkg/response"
	"maxl3oss/pkg/utils"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type UploadController struct {
	DB *gorm.DB
}

func NewUploadController(db *gorm.DB) *UploadController {
	return &UploadController{DB: db}
}

func (u *UploadController) HandleUpload(c *fiber.Ctx) error {
	// get form
	form, err := c.MultipartForm()
	if err != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, err.Error())
	}

	// get form files
	files := form.File["files"]
	for _, file := range files {
		log.Printf("Processing uploaded file %s", filepath.Base(file.Filename))
	}

	// Process each uploaded file
	for _, file := range files {
		err := utils.ProcessFile(file)
		if err != nil {
			return response.Message(c, fiber.StatusInternalServerError, false, err.Error())
		}
	}

	return response.Message(c, fiber.StatusOK, true, "Files uploaded successfully!")
}
