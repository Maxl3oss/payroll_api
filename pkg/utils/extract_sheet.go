package utils

import (
	"errors"
	"log"
	"maxl3oss/app/models"
	"strings"

	"github.com/shakinm/xlsReader/xls"
	"github.com/xuri/excelize/v2"
	"gorm.io/gorm"
)

// create salary
func createManySalary(DB *gorm.DB, dataSalary []models.Salary, dateInfo string, typeID uint) error {
	for _, salary := range dataSalary {
		date, err := ToThaiTime(dateInfo)
		if err != nil {
			return err
		}
		salary.CreatedAt = date
		salary.SalaryTypeID = typeID

		// check data in month
		resultCheckSalary := DB.Where(&models.Salary{FullName: salary.FullName}).Where("EXTRACT(YEAR FROM DATE(created_at)) = ? AND EXTRACT(MONTH FROM DATE(created_at)) = ?", date.Year(), date.Month()).Find(&salary)
		if resultCheckSalary.Error != nil {
			return resultCheckSalary.Error
		}
		if resultCheckSalary.RowsAffected > 0 {
			return errors.New("รายการในไฟล์ในเดือนนี้ มีข้อมูลแล้ว")
		}
		log.Printf("%+v", resultCheckSalary.RowsAffected)

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
	makeNewUser := models.User{
		Email:    transfer.Email,
		Password: GeneratePassword(transfer.MobileNo),
		FullName: salary.FullName,
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

// for process
func ProcessFileBack(DB *gorm.DB, path string, dateInfo string, salaryType models.SalaryType, ext string) error {
	var err error
	var xlsFile xls.Workbook
	var xlsxFile *excelize.File
	var dataSalary []models.Salary
	var dataTransfer []models.TransferInfo

	if ext == ".xls" {
		// read file xls
		xlsFile, err = xls.OpenFile(path)
		if err != nil {
			return err
		}

		// get sheet detail
		dataTransfer, err = extractSheetDetail(xlsFile)
		if err != nil {
			return err
		}
	} else {
		// read files xlsx
		xlsxFile, err = excelize.OpenFile(path)
		if err != nil {
			return err
		}
		defer xlsxFile.Close()
	}

	switch salaryType.Name {
	case "รพสต.":
		if dataSalary, err = extractSheetSalaryHospital(xlsxFile); err != nil {
			return err
		}
		if dataTransfer, err = extractSheetDetailHospital(xlsxFile); err != nil {
			return err
		}
	case "สจ.":
		if dataSalary, err = extractSheetSalaryConsultant(xlsFile); err != nil {
			return err
		}
	case "ฝ่ายประจำ":
		if dataSalary, err = extractSheetSalaryDepartment(xlsFile); err != nil {
			return err
		}
	}

	// Loop through each salary data
	for idx, salary := range dataSalary {
		// Loop through each transfer data
		for _, transfer := range dataTransfer {
			// Check if the full names match
			if salary.FullName == transfer.ReceiverName || salary.BankAccountNumber == transfer.ReceivingACNo {
				//  Check user have?
				var user models.User
				check := DB.Where(&models.User{Email: transfer.Email, FullName: salary.FullName}).First(&user)
				if check.Error == nil {
					dataSalary[idx].UserID = &user.ID
					// log.Printf("old user 1 -> %+v", dataSalary[idx].UserID)
					break
				}

				// Perform the operation to create the user
				newUser, err := createUser(DB, transfer, salary)
				if err != nil {
					return err
				}

				dataSalary[idx].UserID = &newUser.ID
				// log.Printf("new user 1 -> %+v", dataSalary[idx].UserID)
				break
			}
		}
	}

	// Create many salaries
	err = createManySalary(DB, dataSalary, dateInfo, salaryType.ID)
	if err != nil {
		return err
	}

	return nil
}

// ==================================== sheet detail xls
func extractSheetDetail(xlsFile xls.Workbook) ([]models.TransferInfo, error) {
	// targetSheetName := "Detail"

	// Loop through all sheets to find the one with the target name
	// var targetSheet *xls.WorkSheet
	// for i := 0; i < xlsFile.NumSheets(); i++ {
	// 	if xlsFile.GetSheet(i).Name == targetSheetName {
	// 		targetSheet = xlsFile.GetSheet(i)
	// 		break
	// 	}
	// }

	// if targetSheet == nil {
	// 	return nil, errors.New("sheet not found")
	// }

	// add data form xlsx to models
	var transferInfos []models.TransferInfo
	// fmt.Println(isInt(targetSheet.Row(35).Col(0)))
	// fmt.Println(isInt(targetSheet.Row(36).Col(0)))

	// Iterate through rows and columns to get the data
	// for row := 3; row <= 30; row++ {
	// 	rowData := targetSheet.Row(row)
	// 	if rowData == nil {
	// 		continue
	// 	}
	// 	if check := isInt.MatchString(rowData.Col(0)); !check {
	// 		continue
	// 	}

	// 	for i := 0; i <= 14; i++ {
	// 		log.Printf("Column %d: %#v\n", i, rowData.Col(i))
	// 	}

	// 	transferInfo := models.TransferInfo{
	// 		ReceivingBankCode: rowData.Col(0),
	// 		ReceivingACNo:     rowData.Col(1),
	// 		ReceiverName:      rowData.Col(2),
	// 		TransferAmount:    convertFloat(rowData.Col(3)),
	// 		CitizenIDTaxID:    rowData.Col(4),
	// 		DDARef:            rowData.Col(5),
	// 		ReferenceNoDDARef: rowData.Col(6),
	// 		Email:             cleanEmailAddress(rowData.Col(7)),
	// 		MobileNo:          rowData.Col(8),
	// 	}
	// 	transferInfos = append(transferInfos, transferInfo)
	// }

	return transferInfos, nil
}

// ==================================== รพสต. xlsx
func extractSheetSalaryHospital(f *excelize.File) ([]models.Salary, error) {
	rows, err := f.GetRows("เงินเดือน ") // Change "Sheet1" to your sheet name
	if err != nil {
		return nil, err
	}

	// check format col
	formatKeys := []string{
		"ลำดับที่", "ชื่อ-สกุล", "เลขบัญชีธนาคาร", "เงินเดือน", "เงินเพิ่มค่าครองชีพ", "เงินประจำตำแหน่ง", "ค่าตอบแทนรายเดือน", "รวมจ่ายจริง", "ภาษี", "กบข.", "สหกรณ์ออมทรัพย์สาธารณสุข", "กรมสรรพากร(กยศ.)", "ฌกส.", "เงินกู้ธ.กรุงไทย", "รวมรับจริง", "รับจริง", "อื่นๆ",
	}
	if err = CheckFormatCol(rows[3], formatKeys, 21); err != nil {
		return nil, err
	}

	// add data form xlsx to models
	var salaries []models.Salary

	for rowIndex := 7; rowIndex < len(rows)-1; rowIndex++ {
		row := rows[rowIndex]
		if row[0] == "" {
			continue
		}
		if check := isNumeric(row[0]); !check {
			continue
		}

		other := takesFloat(row, 14) + takesFloat(row, 15) + takesFloat(row, 16) + takesFloat(row, 17) + takesFloat(row, 18)

		salary := models.Salary{
			FullName:            takes(row, 1),
			BankAccountNumber:   takes(row, 2),
			Salary:              takesFloat(row, 3),
			AdditionalBenefits:  takesFloat(row, 4),
			FixedIncome:         takesFloat(row, 5),
			MonthlyCompensation: takesFloat(row, 6),
			TotalIncome:         takesFloat(row, 7),

			Tax:                     takesFloat(row, 8),
			Cooperative:             takesFloat(row, 9),
			PublicHealthCooperative: takesFloat(row, 10),
			RevenueDepartment:       takesFloat(row, 11),
			DGS:                     takesFloat(row, 12),
			BangkokBank:             takesFloat(row, 13),
			Other:                   other,
			ActualPay:               takesFloat(row, 19),
			Received:                takesFloat(row, 20),
		}
		salaries = append(salaries, salary)
	}
	return salaries, nil
}

func extractSheetDetailHospital(f *excelize.File) ([]models.TransferInfo, error) {
	rows, err := f.GetRows("Detail")
	if err != nil {
		return nil, err
	}

	// add data form xlsx to models
	var transferInfos []models.TransferInfo

	for rowIndex := 3; rowIndex < len(rows)-10; rowIndex++ {
		row := rows[rowIndex]
		if row[0] == "" {
			continue
		}
		if check := isNumeric(row[0]); !check {
			continue
		}
		transfer := models.TransferInfo{
			ReceivingBankCode: takes(row, 0),
			ReceivingACNo:     takes(row, 1),
			ReceiverName:      takes(row, 2),
			TransferAmount:    takesFloat(row, 3),
			CitizenIDTaxID:    takes(row, 4),
			DDARef:            takes(row, 5),
			ReferenceNoDDARef: takes(row, 6),
			Email:             takes(row, 7),
			MobileNo:          takes(row, 8),
		}
		transferInfos = append(transferInfos, transfer)
	}
	return transferInfos, nil
}

// ==================================== สจ.
func extractSheetSalaryConsultant(xlsFile xls.Workbook) ([]models.Salary, error) {
	targetSheetName := "เงินเดือน"

	// get rows
	rows, err := GetRowsFromSheet(xlsFile, targetSheetName, 12)
	if err != nil {
		return nil, err
	}

	// add data form xlsx to models
	var salaries []models.Salary

	// Iterate through rows and columns to get the data
	for _, row := range rows {
		totalIncome := convertFloat(row.Cells[4]) + convertFloat(row.Cells[5]) + convertFloat(row.Cells[6])
		tax := convertFloat(row.Cells[8])
		savingsBank := convertFloat(row.Cells[9])
		actualPay := tax + savingsBank
		other := convertFloat(row.Cells[11]) + convertFloat(row.Cells[12])
		received := totalIncome - actualPay - other

		salary := models.Salary{
			FullName:            row.Cells[1],
			BankAccountNumber:   row.Cells[2],
			Salary:              convertFloat(row.Cells[4]),
			FixedIncome:         convertFloat(row.Cells[5]),
			MonthlyCompensation: convertFloat(row.Cells[6]),
			TotalIncome:         totalIncome,
			Tax:                 tax,
			Other:               other,
			ActualPay:           actualPay,
			Received:            received,
		}
		salaries = append(salaries, salary)
	}

	return salaries, nil
}

// ==================================== ฝ่ายประจำ
func extractSheetSalaryDepartment(xlsFile xls.Workbook) ([]models.Salary, error) {
	targetSheetName := "เงินเดือน"

	// get rows
	rows, err := GetRowsFromSheet(xlsFile, targetSheetName, 32)
	if err != nil {
		return nil, err
	}
	for i, row := range rows {
		for j, col := range row.Cells {
			log.Printf("row: %d cell: %d => %v", i, j, col)
		}
	}

	// add data form xlsx to models
	var salaries []models.Salary

	// Iterate through rows and columns to get the data
	// for _, row := range rows {
	// 	totalIncome := convertFloat(row.Cells[4]) + convertFloat(row.Cells[5]) + convertFloat(row.Cells[6])
	// 	tax := convertFloat(row.Cells[8])
	// 	savingsBank := convertFloat(row.Cells[9])
	// 	actualPay := tax + savingsBank
	// 	other := convertFloat(row.Cells[11]) + convertFloat(row.Cells[12])
	// 	received := totalIncome - actualPay - other

	// 	salary := models.Salary{
	// 		FullName:            row.Cells[1],
	// 		BankAccountNumber:   row.Cells[2],
	// 		Salary:              convertFloat(row.Cells[4]),
	// 		FixedIncome:         convertFloat(row.Cells[5]),
	// 		MonthlyCompensation: convertFloat(row.Cells[6]),
	// 		TotalIncome:         totalIncome,
	// 		Tax:                 tax,
	// 		Other:               other,
	// 		ActualPay:           actualPay,
	// 		Received:            received,
	// 	}
	// 	salaries = append(salaries, salary)
	// }

	return salaries, nil
}
