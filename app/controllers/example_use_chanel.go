package controllers

import (
	"errors"
	"log"
	"maxl3oss/app/models"
	"maxl3oss/pkg/response"
	"maxl3oss/pkg/utils"
	"mime/multipart"
	"path/filepath"

	"github.com/gofiber/fiber/v2"
	"github.com/xuri/excelize/v2"
)

/*
TODO is upload and process in background
* faster upload but cannot handle error
*/
func (u *SalaryController) uploadSalary(c *fiber.Ctx) error {
	// get form
	form, err := c.MultipartForm()
	if err != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, err.Error())
	}

	// get form files
	files := form.File["files[]"]
	dateInfo := c.FormValue("month")
	// Create a channel to communicate successful file processing
	successCh := make(chan string)

	// Process each uploaded file concurrently
	for _, file := range files {
		go func(file *multipart.FileHeader) {
			// Log file processing
			log.Printf("Processing uploaded file %s", filepath.Base(file.Filename))

			// Validate file extension
			if filepath.Ext(file.Filename) != ".xlsx" {
				response.Message(c, fiber.StatusBadRequest, false, "Only .xlsx files allowed")
				return
			}

			// Save file
			savedFile, pathFile, err := utils.SaveFile(file)
			if err != nil {
				response.Message(c, fiber.StatusInternalServerError, false, err.Error())
				return
			}

			defer savedFile.Close()

			// // check format file
			// err = utils.CheckFormatFileSalary(pathFile)
			// if err != nil {
			// 	response.Message(c, fiber.StatusInternalServerError, false, err.Error())
			// 	return
			// }

			// Send success signal
			successCh <- pathFile
		}(file)
	}

	// Listen for success signal
	go func() {
		for path := range successCh {
			// Extract data from file and create salaries in the background
			go u.processFileBack(path, dateInfo)
		}
	}()

	return response.Message(c, fiber.StatusOK, true, "Files are being processed in the background")
}

func (u *SalaryController) processFileBack(path string, dateInfo string) {
	// Extract data from file
	f, err := excelize.OpenFile(path)

	// dataSalary, err := utils.ExtractSheetSalary(f)
	if err != nil {
		log.Printf("Error extracting data from file: %s", err)
		return
	}
	f.Close()
	// Create many salaries
	err = u.createManySalary([]models.Salary{}, dateInfo)
	if err != nil {
		log.Printf("Error creating salaries: %s", err)
		return
	}

	log.Println("Salaries created successfully")
}

func (u *SalaryController) createManySalary(dataSalary []models.Salary, dateInfo string) error {
	for _, salary := range dataSalary {
		date, err := utils.ToThaiTime(dateInfo)
		if err != nil {
			return err
		}
		salary.CreatedAt = date
		salary.SalaryTypeID = 1

		// check data in month
		resultCheckSalary := u.DB.Where(&models.Salary{FullName: salary.FullName}).Where("EXTRACT(YEAR FROM DATE(created_at)) = ? AND EXTRACT(MONTH FROM DATE(created_at)) = ?", date.Year(), date.Month()).Find(&salary)
		if resultCheckSalary.Error != nil {
			return resultCheckSalary.Error
		}
		if resultCheckSalary.RowsAffected > 0 {
			return errors.New("รายการในไฟล์ในเดือนนี้ มีข้อมูลแล้ว")
		}
		log.Printf("%+v", resultCheckSalary.RowsAffected)

		// Save data
		result := u.DB.Model(&models.Salary{}).Create(&salary)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}
