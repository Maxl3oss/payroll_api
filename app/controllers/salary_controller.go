package controllers

import (
	"errors"
	"log"
	"maxl3oss/app/models"
	"maxl3oss/pkg/response"
	"maxl3oss/pkg/utils"
	"mime/multipart"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

type SalaryController struct {
	DB *gorm.DB
}

func NewSalaryController(db *gorm.DB) *SalaryController {
	return &SalaryController{DB: db}
}

func (u *SalaryController) CreateManySalary(dataSalary []models.Salary, dateInfo string) error {
	for _, salary := range dataSalary {
		date, err := utils.ToThaiTime(dateInfo)
		if err != nil {
			return err
		}
		salary.CreatedAt = date

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

func (u *SalaryController) DeleteManySalary(c *fiber.Ctx) error {
	date, errDate := utils.ToThaiTime(c.Query("month"))
	if errDate != nil {
		return response.Message(c, fiber.ErrBadRequest.Code, false, errDate.Error())
	}

	// check if in month nil
	check := u.DB.Where("EXTRACT(YEAR FROM DATE(created_at)) = ? AND EXTRACT(MONTH FROM DATE(created_at)) = ?", date.Year(), date.Month()).First(&models.User{})
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
	// date, errDate := time.Parse(time.RFC3339, c.Query("month"))
	date, errDate := utils.ToThaiTime(c.Query("month"))

	offset := (page - 1) * limit
	var salaries []models.Salary
	var resCount int64
	var err error

	// count
	queryCount := u.DB.Model(&models.Salary{}).Where("full_name LIKE ?", "%"+search+"%")
	if errDate == nil {
		queryCount = queryCount.Where("EXTRACT(YEAR FROM DATE(created_at)) = ? AND EXTRACT(MONTH FROM DATE(created_at)) = ?", date.Year(), date.Month())
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

func (u *SalaryController) UploadSalaryV2(c *fiber.Ctx) error {
	// get form
	form, err := c.MultipartForm()
	if err != nil {
		return response.Message(c, fiber.StatusInternalServerError, false, err.Error())
	}

	// get form files
	files := form.File["files[]"]
	dateInfo := c.FormValue("month")

	// Begin database transaction
	tx := u.DB.Begin()

	// Process each uploaded file concurrently
	for _, file := range files {
		// Log file processing
		log.Printf("Processing uploaded file %s", filepath.Base(file.Filename))

		// Validate file extension
		if filepath.Ext(file.Filename) != ".xlsx" {
			return response.Message(c, fiber.StatusBadRequest, false, "Only .xlsx files allowed")
		}

		// Save file
		savedFile, pathFile, err := utils.SaveFile(file)
		if err != nil {
			return response.Message(c, fiber.StatusBadRequest, false, err.Error())
		}

		defer savedFile.Close()

		// check format file
		err = utils.CheckFormatFileSalary(pathFile)
		if err != nil {
			return response.Message(c, fiber.StatusBadRequest, false, err.Error())
		}
		// add db
		err = u.ProcessFileBackV2(pathFile, dateInfo)
		if err != nil {
			return response.Message(c, fiber.StatusBadRequest, false, err.Error())
		}

		// break loop if error
		if err != nil {
			break
		}
	}

	// Commit the transaction if all operations succeed
	tx.Commit()
	return response.Message(c, fiber.StatusOK, true, "Upload successfully!")
}

func (u *SalaryController) ProcessFileBackV2(path string, dateInfo string) error {
	// Extract data sheet เงินเดือน
	dataSalary, err := utils.ExtractSheetSalary(path)
	if err != nil {
		return err
	}

	// Extract data from sheet Detail
	dataTransfer, err := utils.ExtractSheetDetail(path)
	if err != nil {
		return err
	}

	// Loop through each salary data
	for idx, salary := range dataSalary {
		// Loop through each transfer data
		for _, transfer := range dataTransfer {
			// Check if the full names match
			if salary.FullName == transfer.ReceiverName || salary.BankAccountNumber == transfer.ReceivingACNo {
				//  Check user have?
				var user models.User
				check := u.DB.Where(&models.User{Email: transfer.Email, FullName: salary.FullName}).First(&user)
				if check.Error == nil {
					dataSalary[idx].UserID = &user.ID
					// log.Printf("old user 1 -> %+v", dataSalary[idx].UserID)
					break
				}
				// log.Printf("%+v", transfer.Email)

				// Create user
				transfer.MobileNo = strings.ReplaceAll(transfer.MobileNo, "-", "")
				transfer.MobileNo = strings.ReplaceAll(transfer.MobileNo, " ", "")
				newUser := models.User{
					Email:          transfer.Email,
					Password:       utils.GeneratePassword(transfer.MobileNo),
					FullName:       salary.FullName,
					CitizenIDTaxID: transfer.CitizenIDTaxID,
					MobileNo:       transfer.MobileNo,
					RoleID:         2,
				}

				// Perform the operation to create the user
				if err := u.DB.Model(&models.User{}).Create(&newUser).Error; err != nil {
					return err
				}

				dataSalary[idx].UserID = &newUser.ID
				// log.Printf("new user 1 -> %+v", dataSalary[idx].UserID)
				break
			}
		}
	}

	// for idx, item := range dataSalary {
	// 	log.Printf("user %+v -> %+v %+v", idx, item.UserID, item.FullName)
	// }
	// log.Printf("data detail 1 -> %+v, %+v", dataTransfer[0].ReceiverName, dataSalary[0].FullName == dataTransfer[0].ReceiverName)
	// Create many salaries
	err = u.CreateManySalary(dataSalary, dateInfo)
	if err != nil {
		return err
	}

	return nil
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

/*
TODO is upload and process in background
* faster upload but cannot handle error
*/
func (u *SalaryController) UploadSalary(c *fiber.Ctx) error {
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
			go u.ProcessFileBack(path, dateInfo)
		}
	}()

	return response.Message(c, fiber.StatusOK, true, "Files are being processed in the background")
}

func (u *SalaryController) ProcessFileBack(path string, dateInfo string) {
	// Extract data from file
	dataSalary, err := utils.ExtractSheetSalary(path)
	if err != nil {
		log.Printf("Error extracting data from file: %s", err)
		return
	}

	// Create many salaries
	err = u.CreateManySalary(dataSalary, dateInfo)
	if err != nil {
		log.Printf("Error creating salaries: %s", err)
		return
	}

	log.Println("Salaries created successfully")
}
