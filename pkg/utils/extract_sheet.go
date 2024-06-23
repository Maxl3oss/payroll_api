package utils

import (
	"errors"
	"maxl3oss/app/models"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// create other name
func createSalaryOther(DB *gorm.DB, others models.TypeOthersName, typeID uint, date time.Time) (uint, error) {
	newData := &models.SalaryOther{
		TypeOthersName: others,
		SalaryTypeID:   typeID,
	}
	newData.CreatedAt = date

	result := DB.Model(&models.SalaryOther{}).Create(newData)
	if result.Error != nil {
		return 0, result.Error
	}
	return newData.ID, nil
}

// create salary
func createManySalary(DB *gorm.DB, dataSalary []models.Salary, dateInfo string, typeID uint, others models.TypeOthersName) error {
	date, err := ToThaiTime(dateInfo)
	if err != nil {
		return err
	}

	// create other
	SOId, resultSalaryOther := createSalaryOther(DB, others, typeID, date)
	if resultSalaryOther != nil {
		return resultSalaryOther
	}

	// loop create salary
	for _, salary := range dataSalary {
		salary.CreatedAt = date
		salary.SalaryTypeID = typeID
		salary.SalaryOtherId = SOId

		// check if fullName empty connect
		if salary.FullName == "" {
			continue
		}

		// check data in month
		resultCheckSalary := DB.Where(&models.Salary{FullName: salary.FullName}).Where("EXTRACT(YEAR FROM DATE(created_at)) = ? AND EXTRACT(MONTH FROM DATE(created_at)) = ?", date.Year(), date.Month()).Find(&salary)
		if resultCheckSalary.Error != nil {
			return resultCheckSalary.Error
		}
		if resultCheckSalary.RowsAffected > 0 {
			return errors.New("รายการในไฟล์ในเดือนนี้ มีข้อมูลแล้ว")
		}

		// Save data
		result := DB.Model(&models.Salary{}).Create(&salary)
		if result.Error != nil {
			return result.Error
		}
	}
	return nil
}

// create user
func createUser(DB *gorm.DB, transfer models.TransferInfo, salary models.Salary) (newUser models.User, err error) {

	// Create user
	transfer.MobileNo = strings.ReplaceAll(transfer.MobileNo, "-", "")
	transfer.MobileNo = strings.ReplaceAll(transfer.MobileNo, " ", "")
	// check email
	var email string

	trimmedEmail := strings.TrimSpace(transfer.Email)
	trimmedCitizenIDTaxID := strings.TrimSpace(transfer.CitizenIDTaxID)

	// Check if Email is not empty or not equal to "-"
	if isValidEmail(trimmedEmail) {
		email = trimmedEmail
	} else if trimmedCitizenIDTaxID != "" && trimmedCitizenIDTaxID != "-" {
		email = trimmedCitizenIDTaxID
	} else {
		// Fallback to random email template
		email = randomEmailTemplate()
	}

	// check password
	password := transfer.MobileNo
	if password == "" {
		password = randomNumericString(10)
	}

	makeNewUser := models.User{
		Email:    email,
		Password: GeneratePassword(password),
		FullName: trimAllSpace(salary.FullName),
		TaxID:    transfer.CitizenIDTaxID,
		Mobile:   transfer.MobileNo,
		RoleID:   2,
	}

	// Perform the operation to create the user
	if err := DB.Model(&models.User{}).Create(&makeNewUser).Error; err != nil {
		return models.User{}, err
	}

	return makeNewUser, nil
}

func updateUser(DB *gorm.DB, transfer models.TransferInfo) error {
	transfer.MobileNo = strings.ReplaceAll(transfer.MobileNo, "-", "")
	transfer.MobileNo = strings.ReplaceAll(transfer.MobileNo, " ", "")
	// check email
	var email string

	trimmedEmail := strings.TrimSpace(transfer.Email)
	trimmedCitizenIDTaxID := strings.TrimSpace(transfer.CitizenIDTaxID)

	// Check if Email is not empty or not equal to "-"
	if isValidEmail(trimmedEmail) {
		email = trimmedEmail
	} else if trimmedCitizenIDTaxID != "" && trimmedCitizenIDTaxID != "-" {
		email = trimmedCitizenIDTaxID
	} else {
		// Fallback to random email template
		email = randomEmailTemplate()
	}

	makeNewUser := models.User{
		Email:    email,
		FullName: trimAllSpace(transfer.ReceiverName),
		TaxID:    transfer.CitizenIDTaxID,
		Mobile:   transfer.MobileNo,
		RoleID:   2,
	}

	if err := DB.Model(&models.User{}).Create(&makeNewUser).Error; err != nil {
		return err
	}

	return nil
}

// for process
func ProcessFileBack(DB *gorm.DB, path string, dateInfo string, salaryType models.SalaryType, others models.TypeOthersName) error {
	var err error
	var xlsxFile *excelize.File
	var dataSalary []models.Salary
	var dataTransfer []models.TransferInfo
	// for sheet detail
	var targetSheet = &TypeTargetSheet{
		Name:  "Detail",
		Cols:  8,
		isUse: true,
	}

	// read files xlsx
	xlsxFile, err = excelize.OpenFile(path)
	if err != nil {
		return err
	}
	defer xlsxFile.Close()

	// get salary
	switch salaryType.Name {
	case "รพสต.":
		if dataSalary, err = extractSheetSalaryHospital(xlsxFile); err != nil {
			return err
		}
	case "สจ.":
		if dataSalary, err = extractSheetSalaryConsultant(xlsxFile); err != nil {
			return err
		}
	case "ฝ่ายประจำ":
		if dataSalary, err = extractSheetSalaryDepartment(xlsxFile); err != nil {
			return err
		}
	case "บำเหน็จรายเดือน":
		targetSheet.Name = "KTB Corporate (2)"
		if dataSalary, err = extractSheetSalaryMonthlyPension(xlsxFile); err != nil {
			return err
		}
	case "เงินเดือนครู":
		targetSheet.Name = "KTB Corporate Online (3)"
		if dataSalary, err = extractSheetSalaryTeacher(xlsxFile); err != nil {
			return err
		}
	case "บำนาญครู":
		targetSheet.Name = "KTB Corporate 3"
		if dataSalary, err = extractSheetTeacherPension(xlsxFile); err != nil {
			return err
		}
	case "บำนาญข้าราชการ":
		targetSheet.Name = "KTB Corporate (2)"
		if dataSalary, err = extractSheetCivilServantPension(xlsxFile); err != nil {
			return err
		}
	}

	// get details
	if targetSheet.isUse {
		if dataTransfer, err = extractSheetDetail(xlsxFile, targetSheet); err != nil {
			return err
		}
	}

	if salaryType.Name == "บำนาญข้าราชการ-old" {
		// Loop through each salary data
		for idx, salary := range dataSalary {
			//  Check user have?
			var user models.User
			check := DB.Where(&models.User{FullName: trimAllSpace(salary.FullName)}).First(&user)
			if check.Error == nil {
				dataSalary[idx].UserID = &user.ID
				// log.Printf("old user 1 -> %+v", dataSalary[idx].UserID)
				break
			}

			// Perform the operation to create the user
			newUser, err := createUser(DB, models.TransferInfo{}, salary)
			if err != nil {
				return err
			}

			dataSalary[idx].UserID = &newUser.ID
		}
	} else {
		// Loop through each salary data
		for idx, salary := range dataSalary {
			// Loop through each transfer data
			for _, transfer := range dataTransfer {
				// Check if the full names match
				if salary.FullName == transfer.ReceiverName || salary.BankAccountNumber == transfer.ReceivingACNo {
					//  Check user have?
					var user models.User
					check := DB.Where(&models.User{FullName: trimAllSpace(transfer.ReceiverName)}).Or(&models.User{TaxID: transfer.CitizenIDTaxID}).First(&user)
					if check.Error == nil {
						dataSalary[idx].UserID = &user.ID
						errUpdate := updateUser(DB, transfer)
						if errUpdate != nil {
							return errUpdate
						}
						break
					}

					// Perform the operation to create the user
					newUser, err := createUser(DB, transfer, salary)
					if err != nil {
						return err
					}

					dataSalary[idx].UserID = &newUser.ID
					break
				}
			}
		}
	}

	// Create many salaries
	err = createManySalary(DB, dataSalary, dateInfo, salaryType.ID, others)
	if err != nil {
		return err
	}

	return nil
}
