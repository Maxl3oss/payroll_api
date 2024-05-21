package models

import "time"

type User struct {
	CommonModelKeyStringFields
	FullName  string `json:"full_name"`          // ชื่อ
	Email     string `json:"email"`              // อีเมล
	TaxID     string `json:"taxid"`              // taxId
	Password  string `json:"password,omitempty"` // รหัสผ่าน
	Mobile    string `json:"mobile"`             // เบอร์โทร
	RoleID    uint   `json:"role_id"`            // สิทะฺเข้าใช้
	Role      Role   `json:"role"`
	OTP       string
	OTPExpiry time.Time
}
