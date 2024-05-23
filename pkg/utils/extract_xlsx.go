package utils

import (
	"maxl3oss/app/models"
	"strings"

	"github.com/xuri/excelize/v2"
)

type TypeTargetSheet struct {
	Name  string
	Cols  int8
	isUse bool
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
		// Skip non-numeric index rows
		indexNo := takes(row, 0)
		if !isNumeric(indexNo) || strings.TrimSpace(indexNo) == "" {
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

func extractSheetDetailHospital(f *excelize.File, targetSheet *TypeTargetSheet) ([]models.TransferInfo, error) {
	rows, err := f.GetRows(targetSheet.Name)
	if err != nil {
		return nil, err
	}

	// add data form xlsx to models
	var transferInfos []models.TransferInfo

	for rowIndex := 3; rowIndex < len(rows)-10; rowIndex++ {
		row := rows[rowIndex]
		indexNo := takes(row, 0)
		if !isNumeric(indexNo) || strings.TrimSpace(indexNo) == "" {
			continue
		}

		// นับ 0
		var email string
		var mobileNo string
		if targetSheet.Cols == 8 {
			email = takes(row, 7)
			mobileNo = takes(row, 8)
		} else {
			// if not email
			email = ""
			mobileNo = takes(row, 7)
		}

		transfer := models.TransferInfo{
			ReceivingBankCode: takes(row, 0),
			ReceivingACNo:     takes(row, 1),
			ReceiverName:      takes(row, 2),
			TransferAmount:    takesFloat(row, 3),
			CitizenIDTaxID:    takes(row, 4),
			DDARef:            takes(row, 5),
			ReferenceNoDDARef: takes(row, 6),
			Email:             email,
			MobileNo:          mobileNo,
		}
		transferInfos = append(transferInfos, transfer)
	}
	return transferInfos, nil
}

// ==================================== สจ.
func extractSheetSalaryConsultant(f *excelize.File) ([]models.Salary, error) {
	rows, err := f.GetRows("เงินเดือน")
	if err != nil {
		return nil, err
	}

	// add data form xlsx to models
	var salaries []models.Salary

	for rowIndex := 0; rowIndex < len(rows); rowIndex++ {
		row := rows[rowIndex]

		// Skip non-numeric index rows
		indexNo := takes(row, 0)
		if !isNumeric(indexNo) || strings.TrimSpace(indexNo) == "" {
			continue
		}

		other := takesFloat(row, 11) + takesFloat(row, 12)

		salary := models.Salary{
			FullName:                  takes(row, 1),
			BankAccountNumber:         takes(row, 2),
			Salary:                    takesFloat(row, 4),
			FixedIncome:               takesFloat(row, 5),
			MonthlyCompensation:       takesFloat(row, 6),
			TotalIncome:               takesFloat(row, 7),
			Tax:                       takesFloat(row, 8),
			CentralWorldSavingsBranch: takesFloat(row, 9),
			ActualPay:                 takesFloat(row, 10),
			Other:                     other,
			Received:                  takesFloat(row, 13),
		}
		salaries = append(salaries, salary)
	}

	return salaries, nil
}

// ==================================== บำเน็จรายเดือน
func extractSheetSalaryMonthlyPension(f *excelize.File) ([]models.Salary, error) {
	// get rows
	rows, err := f.GetRows("รายละเอียดประกอบรายงานเช็ค")
	if err != nil {
		return nil, err
	}

	// add data form xlsx to models
	var salaries []models.Salary

	for rowIndex := 0; rowIndex < len(rows); rowIndex++ {

		row := rows[rowIndex]

		// Skip non-numeric index rows
		indexNo := takes(row, 1)
		if !isNumeric(indexNo) || strings.TrimSpace(indexNo) == "" {
			continue
		}

		monthlyPension := takesFloat(row, 7)
		salaryPeriod := takesFloat(row, 8)
		totalIncome := takesFloat(row, 9)

		payDamages := takesFloat(row, 10)
		coopSaving := takesFloat(row, 11)
		GC := takesFloat(row, 12)
		phuketSavingsBranch := takesFloat(row, 13)
		bangkokBankLoan := takesFloat(row, 14)
		other := takesFloat(row, 15) + takesFloat(row, 16) + takesFloat(row, 17) + takesFloat(row, 18) + takesFloat(row, 19)

		actualPay := takesFloat(row, 20)
		received := takesFloat(row, 21)

		salary := models.Salary{
			FullName:            takes(row, 2),
			BankAccountNumber:   takes(row, 3),
			MonthlyPension:      monthlyPension,
			SalaryPeriod:        salaryPeriod,
			TotalIncome:         totalIncome,
			PayDamages:          payDamages,
			OBACoop:             coopSaving,
			GC:                  GC,
			PhuketSavingsBranch: phuketSavingsBranch,
			BangkokBankLoan:     bangkokBankLoan,
			Other:               other,
			ActualPay:           actualPay,
			Received:            received,
		}
		salaries = append(salaries, salary)
	}
	// Skip non-numeric index rows
	return salaries, nil
}

// ==================================== ฝ่ายประจำ
func extractSheetSalaryDepartment(f *excelize.File) ([]models.Salary, error) {
	// get rows
	rows, err := f.GetRows("เงินเดือน")
	if err != nil {
		return nil, err
	}

	// add data form xlsx to models
	var salaries []models.Salary

	// Iterate through rows and columns to get the data
	for rowIndex := 0; rowIndex < len(rows); rowIndex++ {

		row := rows[rowIndex]

		// Skip non-numeric index rows
		indexNo := takes(row, 0)
		if !isNumeric(indexNo) || strings.TrimSpace(indexNo) == "" {
			continue
		}

		salaryInput := takesFloat(row, 3)
		additionalBenefits := takesFloat(row, 4)
		fixedIncome := takesFloat(row, 5)
		monthlyCompensation := takesFloat(row, 6)
		totalIncome := takesFloat(row, 7)

		tax := takesFloat(row, 8)
		socialSecurityDeduction := takesFloat(row, 9)
		socialSecurityWelfare := takesFloat(row, 10)
		OBACoop := takesFloat(row, 11)
		teacherCoop := takesFloat(row, 12)
		ministryOfPublicHealthCoop := takesFloat(row, 13)
		RHSCoop := takesFloat(row, 14)
		phuketSavingsBranch := takesFloat(row, 15)
		centralWorldSavingsBranch := takesFloat(row, 16)
		revenueDepartment := takesFloat(row, 17)
		bangkokBank := takesFloat(row, 18)
		publicHealthCooperative := takesFloat(row, 19)
		CPKP := takesFloat(row, 20)
		CPKS := takesFloat(row, 21)
		DPC := takesFloat(row, 22)
		GC := takesFloat(row, 23)
		islamicBankLoan := takesFloat(row, 24)
		other := takesFloat(row, 25) + takesFloat(row, 26) + takesFloat(row, 27) + takesFloat(row, 28) + takesFloat(row, 29)

		actualPay := takesFloat(row, 30)
		received := takesFloat(row, 31)

		salary := models.Salary{
			FullName:                   takes(row, 1),
			BankAccountNumber:          takes(row, 2),
			Salary:                     salaryInput,
			AdditionalBenefits:         additionalBenefits,
			FixedIncome:                fixedIncome,
			MonthlyCompensation:        monthlyCompensation,
			TotalIncome:                totalIncome,
			Tax:                        tax,
			SocialSecurityDeduction:    socialSecurityDeduction,
			SocialSecurityWelfare:      socialSecurityWelfare,
			OBACoop:                    OBACoop,
			TeacherCoop:                teacherCoop,
			MinistryOfPublicHealthCoop: ministryOfPublicHealthCoop,
			RHSCoop:                    RHSCoop,
			PhuketSavingsBranch:        phuketSavingsBranch,
			CentralWorldSavingsBranch:  centralWorldSavingsBranch,
			RevenueDepartment:          revenueDepartment,
			BangkokBank:                bangkokBank,
			PublicHealthCooperative:    publicHealthCooperative,
			CPKP:                       CPKP,
			CPKS:                       CPKS,
			DPC:                        DPC,
			GC:                         GC,
			IslamicBankLoan:            islamicBankLoan,
			Other:                      other,
			ActualPay:                  actualPay,
			Received:                   received,
		}
		salaries = append(salaries, salary)

	}

	return salaries, nil
}

// ==================================== เงินเดือนครู
func extractSheetSalaryTeacher(f *excelize.File) ([]models.Salary, error) {
	rows, err := f.GetRows("รายละเอียดประกอบรายงานเช็ค")
	if err != nil {
		return nil, err
	}

	// add data form xlsx to models
	var salaries []models.Salary

	for rowIndex := 0; rowIndex < len(rows); rowIndex++ {
		row := rows[rowIndex]

		// Skip non-numeric index rows
		indexNo := takes(row, 0)
		if !isNumeric(indexNo) || strings.TrimSpace(indexNo) == "" {
			continue
		}

		other := takesFloat(row, 33) + takesFloat(row, 34) + takesFloat(row, 35) + takesFloat(row, 36) + takesFloat(row, 37)

		salary := models.Salary{
			FullName:          takes(row, 1),
			BankAccountNumber: takes(row, 2),
			Salary:            takesFloat(row, 4),
			SalaryPeriod:      takesFloat(row, 5),

			LivingAllowance:     takesFloat(row, 7),
			MonthlyCompensation: takesFloat(row, 8),

			AcademicAllowance: takesFloat(row, 11),
			TotalIncome:       takesFloat(row, 13),

			Tax:                         takesFloat(row, 14),
			SocialSecurityDeduction:     takesFloat(row, 15),
			SocialSecurityWelfare:       takesFloat(row, 16),
			CooperativeAdditional:       takesFloat(row, 17),
			CPKP:                        takesFloat(row, 18),
			CPKS:                        takesFloat(row, 19),
			Cooperative:                 takesFloat(row, 20),
			RevenueDepartment:           takesFloat(row, 21),
			GSC:                         takesFloat(row, 22),
			PrivateCompany:              takesFloat(row, 23),
			IslamicBankLoan:             takesFloat(row, 24),
			TeachersSavingsCoopSurat:    takesFloat(row, 25),
			PhuketSavingsBranch:         takesFloat(row, 26),
			PatongSavingsBranch:         takesFloat(row, 27),
			PoonPholSavingsBranch:       takesFloat(row, 28),
			CentralWorldSavingsBranch:   takesFloat(row, 29),
			CherngTalaySavingsBranch:    takesFloat(row, 30),
			HomeProChalongSavingsBranch: takesFloat(row, 31),
			TeachersSavingsCoop:         takesFloat(row, 32),

			Other:     other,
			ActualPay: takesFloat(row, 38),
			Received:  takesFloat(row, 39),
		}
		salaries = append(salaries, salary)
	}

	return salaries, nil
}

// ==================================== ครูบำนาญครู
func extractSheetTeacherPension(f *excelize.File) ([]models.Salary, error) {
	rows, err := f.GetRows("รายละเอียดประกอบรายงานเช็ค")
	if err != nil {
		return nil, err
	}

	// add data form xlsx to models
	var salaries []models.Salary

	for rowIndex := 0; rowIndex < len(rows); rowIndex++ {
		row := rows[rowIndex]

		// Skip non-numeric index rows
		indexNo := takes(row, 0)
		if !isNumeric(indexNo) || strings.TrimSpace(indexNo) == "" {
			continue
		}

		other := takesFloat(row, 17) + takesFloat(row, 18) + takesFloat(row, 19) + takesFloat(row, 20) + takesFloat(row, 21)

		salary := models.Salary{
			FullName:          takes(row, 1),
			BankAccountNumber: takes(row, 2),
			NormalPension:     takesFloat(row, 4),
			IncreasePension:   takesFloat(row, 5),
			LumpSum:           takesFloat(row, 6),
			LumpSumWithdrawal: takesFloat(row, 7),
			TotalIncome:       takesFloat(row, 8),

			Tax:                         takesFloat(row, 9),
			TeachersSavingsCoop:         takesFloat(row, 10),
			PhuketSavingsBranch:         takesFloat(row, 11),
			PatongSavingsBranch:         takesFloat(row, 12),
			HomeProChalongSavingsBranch: takesFloat(row, 13),
			CentralWorldSavingsBranch:   takesFloat(row, 14),
			CPKP:                        takesFloat(row, 15),
			CPKS:                        takesFloat(row, 16),
			Other:                       other,

			ActualPay: takesFloat(row, 22),
			Received:  takesFloat(row, 23),
		}
		salaries = append(salaries, salary)
	}

	return salaries, nil
}

// ==================================== ครูบำนาญครู
func extractSheetCivilServantPension(f *excelize.File) ([]models.Salary, error) {
	rows, err := f.GetRows("รายละเอียดประกอบรายงานเช็ค")
	if err != nil {
		return nil, err
	}

	// add data form xlsx to models
	var salaries []models.Salary

	for rowIndex := 0; rowIndex < len(rows); rowIndex++ {
		row := rows[rowIndex]

		// Skip non-numeric index rows
		indexNo := takes(row, 0)
		if !isNumeric(indexNo) || strings.TrimSpace(indexNo) == "" {
			continue
		}

		other := takesFloat(row, 22) + takesFloat(row, 23) + takesFloat(row, 24) + takesFloat(row, 25) + takesFloat(row, 26) + takesFloat(row, 27) + takesFloat(row, 28) + takesFloat(row, 29)

		salary := models.Salary{
			FullName:          takes(row, 1),
			BankAccountNumber: takes(row, 2),
			IncreasePension:   takesFloat(row, 4),
			NormalPension:     takesFloat(row, 5),
			LumpSumWithdrawal: takesFloat(row, 6),
			LumpSum:           takesFloat(row, 7),
			TotalIncome:       takesFloat(row, 8),

			Tax:         takesFloat(row, 9),
			OBACoop:     takesFloat(row, 10),
			TeacherCoop: takesFloat(row, 11),
			RHSCoop:     takesFloat(row, 12),
			CPKS:        takesFloat(row, 13),

			CPKP:                      takesFloat(row, 15),
			DPC:                       takesFloat(row, 16),
			GC:                        takesFloat(row, 17),
			BangkokBank:               takesFloat(row, 18),
			IslamicBankLoan:           takesFloat(row, 19),
			PhuketSavingsBranch:       takesFloat(row, 20),
			CentralWorldSavingsBranch: takesFloat(row, 21),

			Other:     other,
			ActualPay: takesFloat(row, 30),
			Received:  takesFloat(row, 31),
		}
		salaries = append(salaries, salary)
	}

	return salaries, nil
}
