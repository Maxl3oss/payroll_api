package models

type Salary struct {
	CommonModelFields
	FullName string `json:"full_name"` // ชื่อ-สกุล

	// สีเขียว / รายจ่าย
	BankAccountNumber   string  `json:"bank_account_number"`  // เลขบัญชีธนาคาร
	Salary              float64 `json:"salary"`               // เงินเดือน
	SalaryPeriod        float64 `json:"salary_period"`        // เงินเดือน(ตกเบิก)
	AdditionalBenefits  float64 `json:"additional_benefits"`  // เงินเพิ่มค่าครองชีพ
	FixedIncome         float64 `json:"fixed_income"`         // เงินประจําตําแหน่ง
	MonthlyCompensation float64 `json:"monthly_compensation"` // ค่าตอบแทนรายเดือน
	IncreasePension     float64 `json:"increase_pension"`     // เงินเพิ่มบำนาญ
	NormalPension       float64 `json:"normal_pension"`       // บำนาญปกติ
	LumpSumWithdrawal   float64 `json:"lump_sum_withdrawal"`  // ช.ค.บ.(ตบเบิก)
	LumpSum             float64 `json:"lump_sum"`             // ช.ค.บ.
	MonthlyPension      float64 `json:"monthly_pension"`      // เงินบำเหน็จรายเดือน
	SpecialCompensation float64 `json:"special_compensation"` // ค่าตอบแทนพิเศษ
	LivingAllowance     float64 `json:"living_allowance"`     // ค่าครองชีพ/เงินเพิ่มการครองชีพชั่วคราว
	AcademicAllowance   float64 `json:"academic_allowance"`   // เงินค่าวิทยฐานะ
	TotalIncome         float64 `json:"total_income"`         // รวมรับจริง

	// สีแดง / เงินหัก
	Tax                         float64 `json:"tax"`                             // ภาษี
	PublicHealthCooperative     float64 `json:"public_health_cooperative"`       // สหกรณ์ออมทรัพย์สาธารณสุข
	RevenueDepartment           float64 `json:"revenue_department"`              // กรมสรรพากร(กยศ.)
	SocialSecurityDeduction     float64 `json:"social_security_deduction"`       // ประกันสังคม(งด) #เงินเดือน
	SocialSecurityWelfare       float64 `json:"social_security_welfare"`         // ประกันสังคม(งพ) #ค่าครองชีพ
	MinistryOfPublicHealthCoop  float64 `json:"ministry_of_public_health_coop"`  // สหกรณ์สำนักงานปลัดกระทรวงสาธารณสุข
	OBACoop                     float64 `json:"oba_coop"`                        // สหกรณ์ฯ(อบจ.)
	RHSCoop                     float64 `json:"rhs_coop"`                        // สหกรณ์ฯ๖(รพช.)
	PhuketSavingsBranch         float64 `json:"phuket_savings_branch"`           // เงินกู้ธ.ออมสิน สาขาภูเก็ต
	CentralWorldSavingsBranch   float64 `json:"central_world_savings_branch"`    // เงินกู้ธ.ออมสิน สาขาเซ็นทรัลฯ
	PatongSavingsBranch         float64 `json:"patong_savings_branch"`           // เงินกู้ธ.ออมสิน ป่าตอง
	PoonPholSavingsBranch       float64 `json:"poon_phol_savings_branch"`        // เงินกู้ธ.ออมสิน พูนผล
	CherngTalaySavingsBranch    float64 `json:"cherng_talay_savings_branch"`     // เงินกู้ธ.ออมสิน สาขาเชิงทะเล
	HomeProChalongSavingsBranch float64 `json:"home_pro_chalong_savings_branch"` // เงินกู้ธ.ออมสิน สาขาโฮมโปรห้าแยกฉลอง
	BangkokBankLoan             float64 `json:"bangkok_bank_loan"`               // เงินกู้ธ.กรุงเทพ
	IslamicBankLoan             float64 `json:"islamic_bank_loan"`               // เงินกู้ธ.อิสลาม
	KrungThaiBank               float64 `json:"krung_thai_bank"`                 // เงินกู้ธ.กรุงไทย
	DGS                         float64 `json:"dgs"`                             // ฌกส.
	Cooperative                 float64 `json:"cooperative"`                     // กบข.
	CooperativeAdditional       float64 `json:"cooperative_additional"`          // กบข. หักเพิ่มเติม (2%)
	CPKP                        float64 `json:"cpkp"`                            // ช.พ.ค.
	CPKS                        float64 `json:"cpks"`                            // ช.พ.ส.
	DPC                         float64 `json:"dpc"`                             // ฌปค.
	GC                          float64 `json:"gc"`                              // ก.ฌ.
	GSC                         float64 `json:"gsc"`                             // กสจ.
	PrivateCompany              float64 `json:"private_company"`                 // บมจ.
	PayDamages                  float64 `json:"pay_damages"`                     // ชดเชยค่าเสียหาย
	ActualPay                   float64 `json:"actual_pay"`                      // รวมจ่ายจริง
	TeachersSavingsCoop         float64 `json:"teachers_savings_coop"`           // สหกรณ์ออมทรัพย์ครู
	TeachersSavingsCoopSurat    float64 `json:"teachers_savings_coop_surat"`     // สหกรณ์ออมทรัพย์ครูสุราษฎร์ธานี

	// อื่นๆ
	Other1 float64 `json:"other1"` // อื่นๆ 1
	Other2 float64 `json:"other2"` // อื่นๆ 2
	Other3 float64 `json:"other3"` // อื่นๆ 3
	Other4 float64 `json:"other4"` // อื่นๆ 4
	Other5 float64 `json:"other5"` // อื่นๆ 5
	Other6 float64 `json:"other6"` // อื่นๆ 6
	Other7 float64 `json:"other7"` // อื่นๆ 7
	Other8 float64 `json:"other8"` // อื่นๆ 8

	Received float64 `json:"received"` // รับจริง

	SalaryTypeID uint       `json:"type_id"`     // รหัส type
	SalaryType   SalaryType `json:"salary_type"` // ข้อมูลที่ relation กัน
	UserID       *string    `json:"user_id"`     // รหัส user
	User         User       `json:"user"`        // ข้อมูลที่ relation กัน

	SalaryOtherId uint        `json:"salary_other_id"` // รหัส other
	SalaryOther   SalaryOther `json:"salary_other"`
}

type SalaryType struct {
	ID   uint   `gorm:"primary_key" json:"id"`
	Name string `json:"name"`
}

type SalaryOther struct {
	CommonModelFields
	SalaryTypeID uint `json:"type_id"`
	TypeOthersName
}

type TypeOthersName struct {
	Other1Name string `json:"other1_name"`
	Other2Name string `json:"other2_name"`
	Other3Name string `json:"other3_name"`
	Other4Name string `json:"other4_name"`
	Other5Name string `json:"other5_name"`
	Other6Name string `json:"other6_name"`
	Other7Name string `json:"other7_name"`
	Other8Name string `json:"other8_name"`
}
