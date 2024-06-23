package controllers

import (
	"errors"
	"log"
	"maxl3oss/app/models"
	"maxl3oss/pkg/response"
	"maxl3oss/pkg/utils"
	"path/filepath"
	"strconv"
	"time"

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
	inpType := c.Query("type", "0")
	salaryType, err := strconv.Atoi(inpType)
	if err != nil {
		return response.Message(c, fiber.StatusBadRequest, false, "ข้อมูลรูปแบบไม่ถูกต้อง")
	}

	if errDate != nil {
		return response.Message(c, fiber.StatusBadRequest, false, errDate.Error())
	}

	// Start a transaction
	tx := u.DB.Begin()
	if tx.Error != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, "Failed to start transaction")
	}

	// Define the date condition
	dateCondition := "EXTRACT(YEAR FROM DATE(created_at)) = ? AND EXTRACT(MONTH FROM DATE(created_at)) = ?"
	if salaryType != 0 {
		dateCondition += " AND salary_type_id = ?"
	} else {
		dateCondition += " OR salary_type_id = ?"
	}

	// Delete salary
	resultSalary := tx.Delete(&models.Salary{}, dateCondition, date.Year(), date.Month(), salaryType)
	if resultSalary.Error != nil {
		tx.Rollback()
		return response.Message(c, fiber.StatusInternalServerError, false, "Failed to delete salary data")
	}

	// Delete salary other
	resultOther := tx.Delete(&models.SalaryOther{}, dateCondition, date.Year(), date.Month(), salaryType)
	if resultOther.Error != nil {
		tx.Rollback()
		return response.Message(c, fiber.StatusInternalServerError, false, "Failed to delete salary other data")
	}

	// Check if any rows were affected
	if resultSalary.RowsAffected == 0 && resultOther.RowsAffected == 0 {
		tx.Rollback()
		return response.Message(c, fiber.StatusNotFound, false, "ไม่พบข้อมูลในช่วงเวลานี้")
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, "Failed to commit transaction")
	}

	return response.Message(c, fiber.StatusOK, true, "ลบข้อมูลเงินเดือนสำเร็จ")
}

func (u *SalaryController) GetAll(c *fiber.Ctx) error {
	page, _ := strconv.Atoi(c.Query("pageNumber", "1"))
	limit, _ := strconv.Atoi(c.Query("pageSize", "10"))
	search := c.Query("search")
	inpType := c.Query("type", "0")
	date, _ := utils.ToThaiTime(c.Query("month"))

	offset := (page - 1) * limit
	var salaries []models.Salary
	var resCount int64
	var err error

	salaryType, err := strconv.Atoi(inpType)
	if err != nil {
		return response.Message(c, fiber.StatusBadRequest, false, "ข้อมูลรูปแบบไม่ถูกต้อง")
	}

	// In your main function:
	query := buildQueryGetAll(u.DB, search, &date, salaryType)

	// Count
	if err := query.Count(&resCount).Error; err != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, err.Error())
	}

	// Data
	result := query.Preload("SalaryType").Preload("SalaryOther").
		Limit(limit).Offset(offset).Find(&salaries)

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

func buildQueryGetAll(db *gorm.DB, search string, date *time.Time, salaryType int) *gorm.DB {
	query := db.Model(&models.Salary{}).Where("full_name LIKE ?", "%"+search+"%")

	if date != nil {
		query = query.Where("EXTRACT(YEAR FROM DATE(created_at)) = ? AND EXTRACT(MONTH FROM DATE(created_at)) = ?", date.Year(), date.Month())
	}

	if salaryType != 0 {
		query = query.Where("salary_type_id = ?", uint(salaryType))
	}

	return query
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
	queryData := u.DB.Preload("SalaryType").Model(&models.Salary{}).Where("user_id = ?", userID)
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
	others := convertOthers(c)
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
		case ".xlsx":
			err = utils.ProcessFileBack(u.DB, pathFile, dateInfo, dataSalaryType, others)
		default:
			err = errors.New("only (.xlsx) files allowed")
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

func (u *SalaryController) GetSalaryOther(c *fiber.Ctx) error {
	var dataSalaryOther models.SalaryOther
	inpType := c.Query("type", "0")
	salaryType, err := strconv.Atoi(inpType)
	if err != nil {
		return response.Message(c, fiber.StatusBadRequest, false, "ข้อมูลรูปแบบไม่ถูกต้อง")
	}

	var result = u.DB.Where(&models.SalaryOther{SalaryTypeID: uint(salaryType)}).Order("created_at DESC").First(&dataSalaryOther)
	if result.Error != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, result.Error.Error())
	}
	return response.SendData(c, fiber.StatusOK, true, dataSalaryOther, nil)
}

// get others name
func convertOthers(c *fiber.Ctx) models.TypeOthersName {
	result := models.TypeOthersName{
		Other1Name: c.FormValue("other1_name", ""),
		Other2Name: c.FormValue("other2_name", ""),
		Other3Name: c.FormValue("other3_name", ""),
		Other4Name: c.FormValue("other4_name", ""),
		Other5Name: c.FormValue("other5_name", ""),
		Other6Name: c.FormValue("other6_name", ""),
		Other7Name: c.FormValue("other7_name", ""),
		Other8Name: c.FormValue("other8_name", ""),
	}

	return result
}
