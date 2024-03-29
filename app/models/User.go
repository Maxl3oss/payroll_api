package models

type User struct {
	CommonModelKeyStringFields
	FullName       string `json:"full_name"`          // ชื่อ
	Email          string `json:"email"`              // อีเมล
	CitizenIDTaxID string `json:"taxid"`              // taxId
	Password       string `json:"password,omitempty"` // รหัสผ่าน
	MobileNo       string `json:"mobile"`             // เบอร์โทร
	RoleID         int16  `json:"role_id"`            // สิทะฺเข้าใช้
	Role           Role   `json:"role"`
}