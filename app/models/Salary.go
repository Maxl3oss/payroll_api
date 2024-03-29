package models

type Salary struct {
	CommonModelFields
	FullName            string  `json:"full_name"`            // ชื่อ-สกุล
	BankAccountNumber   string  `json:"bank_account_number"`  // เลขบัญชีธนาคาร
	Salary              float64 `json:"salary"`               // เงินเดือน
	AdditionalBenefits  float64 `json:"additional_benefits"`  // เงินเพิ่มค่าครองชีพ
	FixedIncome         float64 `json:"fixed_income"`         // เงินประจําตําแหน่ง
	MonthlyCompensation float64 `json:"monthly_compensation"` // ค่าตอบแทนรายเดือน
	TotalIncome         float64 `json:"total_income"`         // รวมรับจริง

	Tax                     float64 `json:"tax"`                       // ภาษี
	Cooperative             float64 `json:"coop"`                      // กบข.
	PublicHealthCooperative float64 `json:"public_health_cooperative"` // สหกรณฯ(สาธารณสุข)
	RevenueDepartment       float64 `json:"revenue_department"`        // กรมสรรพากร(กยศ.)
	DGS                     float64 `json:"dgs"`                       // ฌกส.
	BangkokBank             float64 `json:"bangkok_bank"`              // ธ.กรุงไทย
	ActualPay               float64 `json:"actual_pay"`                // รวมจ่ายจริง

	Received float64 `json:"received"` // รับจริง

	SocialSecurity string `json:"social_security"` // รวมประกันสังคม
	BankTransfer   string `json:"bank_transfer"`   // ส่งธนาคาร

	UserID *string `json:"user_id"` // รหัส user
	User   User    `json:"user"`    // ข้อมูลที่ relation กัน
}
