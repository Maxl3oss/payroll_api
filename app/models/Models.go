package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type CommonModelFields struct {
	ID        uint       `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

type CommonModelKeyStringFields struct {
	ID        string     `gorm:"primary_key" json:"id"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at"`
}

func (cm *CommonModelKeyStringFields) BeforeCreate(tx *gorm.DB) error {
	uuid := uuid.New().String()
	tx.Statement.SetColumn("ID", uuid)
	return nil
}

type Dashboard struct {
	Received        float64           `json:"received"`          // รับจริง
	User            int64             `json:"user"`              // รับจริง
	ReceivedByMonth []ReceivedByMonth `json:"received_by_month"` // รายการตามเดือน
}

type ReceivedByMonth struct {
	Year  int     `json:"year"`  // ปี
	Month int     `json:"month"` // เดือน
	Sum   float64 `json:"sum"`   // รวม
}
