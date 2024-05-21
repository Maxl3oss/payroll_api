package models

type Salary struct {
	CommonModelFields
	FullName string `json:"full_name"` // ชื่อ-สกุล

	BankAccountNumber   string  `json:"bank_account_number"`  // เลขบัญชีธนาคาร
	Salary              float64 `json:"salary"`               // เงินเดือน
	AdditionalBenefits  float64 `json:"additional_benefits"`  // เงินเพิ่มค่าครองชีพ
	FixedIncome         float64 `json:"fixed_income"`         // เงินประจําตําแหน่ง
	MonthlyCompensation float64 `json:"monthly_compensation"` // ค่าตอบแทนรายเดือน
	TotalIncome         float64 `json:"total_income"`         // รวมรับจริง

	Tax                     float64 `json:"tax"`                       // ภาษี
	Cooperative             float64 `json:"coop"`                      // กบข.
	PublicHealthCooperative float64 `json:"public_health_cooperative"` // สหกรณ์ออมทรัพย์สาธารณสุข
	RevenueDepartment       float64 `json:"revenue_department"`        // กรมสรรพากร(กยศ.)
	DGS                     float64 `json:"dgs"`                       // ฌกส.
	BangkokBank             float64 `json:"bangkok_bank"`              // เงินกู้ธ.กรุงไทย
	Other                   float64 `json:"other"`                     // อื่นๆ (รวม)
	ActualPay               float64 `json:"actual_pay"`                // รวมจ่ายจริง

	Received float64 `json:"received"` // รับจริง

	SalaryTypeID uint       `json:"type_id"`
	UserID       *string    `json:"user_id"`     // รหัส user
	SalaryType   SalaryType `json:"salary_type"` // ข้อมูลที่ relation กัน
	User         User       `json:"user"`        // ข้อมูลที่ relation กัน
}

type SalaryType struct {
	ID   uint   `json:"id"`
	Name string `json:"name"`
}
