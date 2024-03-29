package models

type TransferInfo struct {
	CommonModelFields
	ReceivingBankCode string  `json:"receivingBankCode"` // รหัสธนาคาร
	ReceivingACNo     string  `json:"receivingACNo"`     // เลขที่บัญชี
	ReceiverName      string  `json:"receiverName"`      // ชื่อบัญชี
	TransferAmount    float64 `json:"transferAmount"`    // จำนวนเงิน
	CitizenIDTaxID    string  `json:"citizenIdTaxid"`    // taxId
	DDARef            string  `json:"ddaRef"`            // อ้างอิง
	ReferenceNoDDARef string  `json:"referenceNo"`       // รายการ
	Email             string  `json:"email"`             // อีเมล
	MobileNo          string  `json:"mobileNo"`          // เบอร์โทร
}
