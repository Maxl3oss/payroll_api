package controllers

import (
	"errors"
	"log"
	"maxl3oss/app/models"
	"maxl3oss/pkg/response"
	"maxl3oss/pkg/utils"
	"path/filepath"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SalaryController struct {
	DB *gorm.DB
}

func NewSalaryController(db *gorm.DB) *SalaryController {
	return &SalaryController{DB: db}
}

func (u *SalaryController) DeleteManySalary(c *fiber.Ctx) error {
	date, errDate := utils.ToThaiTime(c.Query("month"))
	if errDate != nil {
		return response.Message(c, fiber.ErrBadRequest.Code, false, errDate.Error())
	}

	// check if in month nil
	check := u.DB.Where("EXTRACT(YEAR FROM DATE(created_at)) = ? AND EXTRACT(MONTH FROM DATE(created_at)) = ?", date.Year(), date.Month()).First(&models.Salary{})
	if check.Error != nil {
		return response.Message(c, fiber.ErrBadRequest.Code, false, "ไม่มีข้อมูลในช่วงเวลานี้")
	}

	result := u.DB.Delete(&models.Salary{}, "EXTRACT(YEAR FROM DATE(created_at)) = ? AND EXTRACT(MONTH FROM DATE(created_at)) = ?", date.Year(), date.Month())
	if result.Error != nil {
		return response.Message(c, fiber.ErrBadRequest.Code, false, result.Error.Error())
	}

	return response.Message(c, fiber.StatusOK, true, "Delete data salary successfully!")
}

func (u *SalaryController) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("pageNumber", "1"))
	limit, _ := strconv.Atoi(c.Query("pageSize", "10"))
	search := c.Query("search")
	inpType := c.Query("type", "0")

	// date, errDate := time.Parse(time.RFC3339, c.Query("month"))
	date, errDate := utils.ToThaiTime(c.Query("month"))

	offset := (page - 1) * limit
	var salaries []models.Salary
	var resCount int64
	var err error

	salaryType, err := strconv.Atoi(inpType)
	if err != nil {
		return response.Message(c, fiber.StatusBadRequest, false, err.Error())
	}

	// count
	queryCount := u.DB.Model(&models.Salary{}).Where("full_name LIKE ?", "%"+search+"%")
	if errDate == nil {
		queryCount = queryCount.Where("EXTRACT(YEAR FROM DATE(created_at)) = ? AND EXTRACT(MONTH FROM DATE(created_at)) = ?", date.Year(), date.Month())
	}
	if salaryType != 0 {
		queryCount = queryCount.Where("salary_type_id = ?", uint(salaryType))
	}
	err = queryCount.Count(&resCount).Error
	if err != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, err.Error())
	}

	// data
	queryData := u.DB.Model(&models.Salary{}).Where("full_name LIKE ?", "%"+search+"%")
	if errDate == nil {
		queryData = queryData.Where("EXTRACT(YEAR FROM DATE(created_at)) = ? AND EXTRACT(MONTH FROM DATE(created_at)) = ?", date.Year(), date.Month())
	}
	if salaryType != 0 {
		queryData = queryData.Where("salary_type_id = ?", uint(salaryType))
	}
	result := queryData.Limit(limit).Offset(offset).Find(&salaries)
	if result.Error != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, result.Error.Error())
	}

	pagin := response.Pagination{
		PageNumber:  page,
		PageSize:    limit,
		TotalRecord: int(resCount),
	}

	return response.SendData(c, fiber.StatusOK, true, salaries, &pagin)
}

func (u *SalaryController) GetByUser(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("pageNumber", "1"))
	limit, _ := strconv.Atoi(c.Query("pageSize", "10"))
	userID := c.Params("id")
	date, errDate := utils.ToThaiTime(c.Query("month"))

	offset := (page - 1) * limit
	var salaries []models.Salary
	var resCount int64
	var err error

	// check user
	check := u.DB.Where("id = ?", userID).First(&models.User{})
	if check.Error != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, check.Error.Error())
	}

	// count
	queryCount := u.DB.Model(&models.Salary{}).Where("user_id = ?", userID)
	if errDate == nil {
		queryCount = queryCount.Where("EXTRACT(YEAR FROM DATE(created_at)) = ? AND EXTRACT(MONTH FROM DATE(created_at)) = ?", date.Year(), date.Month())
	}
	err = queryCount.Count(&resCount).Error
	if err != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, err.Error())
	}

	// data
	queryData := u.DB.Model(&models.Salary{}).Where("user_id = ?", userID)
	if errDate == nil {
		queryData = queryData.Where("EXTRACT(YEAR FROM DATE(created_at)) = ? AND EXTRACT(MONTH FROM DATE(created_at)) = ?", date.Year(), date.Month())
	}
	result := queryData.Limit(limit).Offset(offset).Find(&salaries)
	if result.Error != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, result.Error.Error())
	}

	pagin := response.Pagination{
		PageNumber:  page,
		PageSize:    limit,
		TotalRecord: int(resCount),
	}

	return response.SendData(c, fiber.StatusOK, true, salaries, &pagin)
}

func (u *SalaryController) GetSalaryType(c *fiber.Ctx) error {
	var dataType []models.SalaryType

	result := u.DB.Model(&models.SalaryType{}).Find(&dataType)
	if result.Error != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, result.Error.Error())
	}

	return response.SendData(c, fiber.StatusOK, true, dataType, nil)
}

func (u *SalaryController) UploadSalary(c *fiber.Ctx) error {
	// get form
	form, err := c.MultipartForm()
	if err != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, err.Error())
	}

	// get form files
	files := form.File["files[]"]
	dateInfo := c.FormValue("month")

	// get salary type
	inpSalaryType := c.FormValue("type", "0")
	var salaryType int
	var dataSalaryType models.SalaryType

	// validation
	if dateInfo == "" {
		return response.Message(c, fiber.StatusBadRequest, false, "ไม่พบข้อมูลเดือน")
	}
	if files == nil {
		return response.Message(c, fiber.StatusBadRequest, false, "ไม่พบไฟล์")
	}
	if inpSalaryType == "0" {
		return response.Message(c, fiber.StatusBadRequest, false, "ไม่มีข้อมูลรูปแบบ")
	} else {
		salaryType, err = strconv.Atoi(inpSalaryType)
		if err != nil {
			return response.Message(c, fiber.StatusBadRequest, false, "ข้อมูลรูปแบบไม่ถูกต้อง")
		}
	}

	// get name salary type
	if result := u.DB.Where(&models.SalaryType{ID: uint(salaryType)}).First(&dataSalaryType); result.Error != nil {
		return response.Message(c, fiber.StatusBadRequest, false, "ไม่พบข้อมูลรูปแบบ")
	}

	// Begin database transaction
	tx := u.DB.Begin()

	// Process each uploaded file concurrently
	for _, file := range files {
		// Log file processing
		log.Printf("Processing uploaded file %s", filepath.Base(file.Filename))

		// Save file
		_, pathFile, err := utils.SaveFile(file)
		if err != nil {
			return response.Message(c, fiber.StatusBadRequest, false, err.Error())
		}

		// Determine file extension
		ext := filepath.Ext(file.Filename)
		switch ext {
		case ".xlsx", ".xls":
			err = utils.ProcessFileBack(u.DB, pathFile, dateInfo, dataSalaryType, ext)
		default:
			err = errors.New("only (.xlsx, .xls) files allowed")
		}

		// Handle errors
		if err != nil {
			return response.Message(c, fiber.StatusBadRequest, false, err.Error())
		}
	}

	// Commit the transaction if all operations succeed
	tx.Commit()
	return response.Message(c, fiber.StatusOK, true, "Upload successfully!")
}
