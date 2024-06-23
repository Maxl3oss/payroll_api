package utils

// func GetRowsFromSheet(xlsFile xls.Workbook, name string, colTotal int, optionalArgs ...int) ([]RowData, error) {
// 	var rows []RowData

// 	// make default index
// 	indexCol := 0
// 	if len(optionalArgs) > 0 {
// 		indexCol = optionalArgs[0]
// 	}
// 	sheet, err := getSheetByName(xlsFile, name)
// 	if err != nil {
// 		return nil, err
// 	}

// 	numRows := sheet.GetNumberRows()
// 	for i := 0; i < numRows; i++ {
// 		rowData, err := extractRowData(sheet, i, colTotal, indexCol)
// 		if err != nil {
// 			return nil, err
// 		}

// 		if rowData != nil {
// 			rows = append(rows, *rowData)
// 		}
// 	}

// 	return rows, nil
// }

// func extractRowData(sheet *xls.Sheet, rowIndex, colTotal int, indexCol int) (*RowData, error) {
// 	rowData := RowData{}

// 	row, err := sheet.GetRow(rowIndex)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// Check if the first cell contains a numeric index
// 	cellFirst, err := row.GetCol(indexCol)
// 	if err != nil {
// 		return nil, err
// 	}
// 	cellIndex := cellFirst.GetString()
// 	if !isNumeric(cellIndex) || strings.TrimSpace(cellIndex) == "" {
// 		// Skip non-numeric index rows
// 		return nil, nil
// 	}

// 	// Extract data from each column
// 	for j := 0; j <= colTotal; j++ {
// 		cell, err := row.GetCol(j)
// 		if err != nil {
// 			return nil, err
// 		}
// 		cellValue := cell.GetString()
// 		rowData.Cells = append(rowData.Cells, cellValue)
// 	}

// 	return &rowData, nil
// }

// func getSheetByName(workbook xls.Workbook, name string) (*xls.Sheet, error) {
// 	numSheets := workbook.GetNumberSheets()
// 	for i := 0; i < numSheets; i++ {
// 		sheet, err := workbook.GetSheet(i)
// 		if err != nil {
// 			return nil, err
// 		}
// 		if sheet.GetName() == name {
// 			return sheet, nil
// 		}
// 	}
// 	return nil, fmt.Errorf("sheet not found: %s", name)
// }

// Function to clean up email addresses
// func cleanEmailAddress(email string) string {
// 	// Define a regular expression to match valid email characters
// 	validEmailRegex := regexp.MustCompile(`[[:alnum:]!#$%&'*+/=?^_` + "`" + `{|}~-]+(?:\.[[:alnum:]!#$%&'*+/=?^_` + "`" + `{|}~-]+)*@(?:[[:alnum:]-]+\.)+[[:alpha:]]{2,7}`)

// 	// Find valid email addresses in the string
// 	validEmail := validEmailRegex.FindString(email)

// 	// Remove any leading or trailing whitespace
// 	cleanEmail := strings.TrimSpace(validEmail)

// 	return cleanEmail
// }

// ==================================== สจ.
// func extractSheetSalaryConsultantXLS(xlsFile xls.Workbook) ([]models.Salary, error) {
// 	targetSheetName := "เงินเดือน"

// 	// get rows
// 	rows, err := GetRowsFromSheet(xlsFile, targetSheetName, 12)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// add data form xlsx to models
// 	var salaries []models.Salary

// 	// Iterate through rows and columns to get the data
// 	for _, row := range rows {
// 		totalIncome := convertFloat(row.Cells[4]) + convertFloat(row.Cells[5]) + convertFloat(row.Cells[6])
// 		tax := convertFloat(row.Cells[8])
// 		savingsBank := convertFloat(row.Cells[9])
// 		actualPay := tax + savingsBank
// 		other := convertFloat(row.Cells[11]) + convertFloat(row.Cells[12])
// 		received := totalIncome - actualPay - other

// 		salary := models.Salary{
// 			FullName:            row.Cells[1],
// 			BankAccountNumber:   row.Cells[2],
// 			Salary:              convertFloat(row.Cells[4]),
// 			FixedIncome:         convertFloat(row.Cells[5]),
// 			MonthlyCompensation: convertFloat(row.Cells[6]),
// 			TotalIncome:         totalIncome,
// 			Tax:                 tax,
// 			Other:               other,
// 			ActualPay:           actualPay,
// 			Received:            received,
// 		}
// 		salaries = append(salaries, salary)
// 	}

// 	return salaries, nil
// }

// ==================================== ฝ่ายประจำ
// func extractSheetSalaryDepartment(xlsFile xls.Workbook) ([]models.Salary, error) {
// 	targetSheetName := "เงินเดือน"

// 	// get rows
// 	rows, err := GetRowsFromSheet(xlsFile, targetSheetName, 29)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// add data form xlsx to models
// 	var salaries []models.Salary

// 	// Iterate through rows and columns to get the data
// 	for _, row := range rows {

// 		salaryInput := convertFloat(row.Cells[3])
// 		additionalBenefits := convertFloat(row.Cells[4])
// 		fixedIncome := convertFloat(row.Cells[5])
// 		monthlyCompensation := convertFloat(row.Cells[6])
// 		totalIncome := salaryInput + additionalBenefits + fixedIncome + monthlyCompensation

// 		tax := convertFloat(row.Cells[8])
// 		socialSecurityDeduction := convertFloat(row.Cells[9])
// 		socialSecurityWelfare := convertFloat(row.Cells[10])
// 		OBACoop := convertFloat(row.Cells[11])
// 		teacherCoop := convertFloat(row.Cells[12])
// 		ministryOfPublicHealthCoop := convertFloat(row.Cells[13])
// 		RHSCoop := convertFloat(row.Cells[14])
// 		phuketSavingsBranch := convertFloat(row.Cells[15])
// 		centralWorldSavingsBranch := convertFloat(row.Cells[16])
// 		revenueDepartment := convertFloat(row.Cells[17])
// 		bangkokBank := convertFloat(row.Cells[18])
// 		publicHealthCooperative := convertFloat(row.Cells[19])
// 		CPKP := convertFloat(row.Cells[20])
// 		CPKS := convertFloat(row.Cells[21])
// 		DPC := convertFloat(row.Cells[22])
// 		GC := convertFloat(row.Cells[23])
// 		islamicBankLoan := convertFloat(row.Cells[24])
// 		other := convertFloat(row.Cells[25]) + convertFloat(row.Cells[26]) + convertFloat(row.Cells[27]) + convertFloat(row.Cells[28]) + convertFloat(row.Cells[29])

// 		actualPay := tax + socialSecurityDeduction + socialSecurityWelfare + OBACoop + teacherCoop + ministryOfPublicHealthCoop + RHSCoop + phuketSavingsBranch + centralWorldSavingsBranch + revenueDepartment + bangkokBank + publicHealthCooperative + CPKP + CPKS + DPC + GC + islamicBankLoan + other
// 		received := totalIncome - actualPay

// 		salary := models.Salary{
// 			FullName:                   row.Cells[1],
// 			BankAccountNumber:          row.Cells[2],
// 			Salary:                     salaryInput,
// 			AdditionalBenefits:         additionalBenefits,
// 			FixedIncome:                fixedIncome,
// 			MonthlyCompensation:        monthlyCompensation,
// 			TotalIncome:                totalIncome,
// 			Tax:                        tax,
// 			SocialSecurityDeduction:    socialSecurityDeduction,
// 			SocialSecurityWelfare:      socialSecurityWelfare,
// 			OBACoop:                    OBACoop,
// 			TeacherCoop:                teacherCoop,
// 			MinistryOfPublicHealthCoop: ministryOfPublicHealthCoop,
// 			RHSCoop:                    RHSCoop,
// 			PhuketSavingsBranch:        phuketSavingsBranch,
// 			CentralWorldSavingsBranch:  centralWorldSavingsBranch,
// 			RevenueDepartment:          revenueDepartment,
// 			BangkokBank:                bangkokBank,
// 			PublicHealthCooperative:    publicHealthCooperative,
// 			CPKP:                       CPKP,
// 			CPKS:                       CPKS,
// 			DPC:                        DPC,
// 			GC:                         GC,
// 			IslamicBankLoan:            islamicBankLoan,
// 			Other:                      other,
// 			ActualPay:                  actualPay,
// 			Received:                   received,
// 		}
// 		salaries = append(salaries, salary)
// 	}

// 	return salaries, nil
// }

// ==================================== บำเน็จรายเดือน
// func extractSheetSalaryMonthlyPension(xlsFile xls.Workbook) ([]models.Salary, error) {
// 	targetSheetName := "รายละเอียดประกอบรายงานเช็ค"

// 	// get rows
// 	rows, err := GetRowsFromSheet(xlsFile, targetSheetName, 18, 1)
// 	if err != nil {
// 		return nil, err
// 	}

// 	// add data form xlsx to models
// 	var salaries []models.Salary

// 	// Iterate through rows and columns to get the data
// 	for _, row := range rows {

// 		monthlyPension := convertFloat(row.Cells[6])
// 		salaryPeriod := convertFloat(row.Cells[7])
// 		totalIncome := monthlyPension + salaryPeriod

// 		payDamages := convertFloat(row.Cells[9])
// 		coopSaving := convertFloat(row.Cells[10])
// 		GC := convertFloat(row.Cells[11])
// 		phuketSavingsBranch := convertFloat(row.Cells[12])
// 		bangkokBankLoan := convertFloat(row.Cells[13])
// 		other := convertFloat(row.Cells[14]) + convertFloat(row.Cells[15]) + convertFloat(row.Cells[16]) + convertFloat(row.Cells[17]) + convertFloat(row.Cells[18])

// 		actualPay := payDamages + coopSaving + GC + phuketSavingsBranch + bangkokBankLoan + other
// 		received := totalIncome - actualPay

// 		salary := models.Salary{
// 			FullName:            row.Cells[2],
// 			BankAccountNumber:   row.Cells[3],
// 			MonthlyPension:      monthlyPension,
// 			SalaryPeriod:        salaryPeriod,
// 			TotalIncome:         totalIncome,
// 			PayDamages:          payDamages,
// 			CoopSaving:          coopSaving,
// 			GC:                  GC,
// 			PhuketSavingsBranch: phuketSavingsBranch,
// 			BangkokBankLoan:     bangkokBankLoan,
// 			Other:               other,
// 			ActualPay:           actualPay,
// 			Received:            received,
// 		}
// 		salaries = append(salaries, salary)
// 	}

// 	return salaries, nil
// }

// ==================================== sheet detail xls
// func extractSheetDetail(xlsFile xls.Workbook, name string, optionalArgs ...bool) ([]models.TransferInfo, error) {
// 	targetSheetName := name
// 	isHaveEmail := true
// 	if len(optionalArgs) > 0 {
// 		isHaveEmail = optionalArgs[0]
// 	}
// 	// get rows
// 	rows, err := GetRowsFromSheet(xlsFile, targetSheetName, 8)
// 	if err != nil {
// 		return nil, err
// 	}

// 	var transferInfos []models.TransferInfo

// 	for _, row := range rows {
// 		var email string
// 		var MobileNo string

// 		if isHaveEmail {
// 			email = cleanEmailAddress(row.Cells[7])
// 			MobileNo = row.Cells[8]
// 		} else {
// 			email = ""
// 			MobileNo = row.Cells[7]
// 		}

// 		transferInfo := models.TransferInfo{
// 			ReceivingBankCode: row.Cells[0],
// 			ReceivingACNo:     row.Cells[1],
// 			ReceiverName:      row.Cells[2],
// 			TransferAmount:    convertFloat(row.Cells[3]),
// 			CitizenIDTaxID:    row.Cells[4],
// 			DDARef:            row.Cells[5],
// 			ReferenceNoDDARef: row.Cells[6],
// 			Email:             email,
// 			MobileNo:          MobileNo,
// 		}
// 		transferInfos = append(transferInfos, transferInfo)
// 	}

// 	return transferInfos, nil
// }
