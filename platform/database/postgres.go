package database

import (
	"log"
	"maxl3oss/app/models"
	"maxl3oss/pkg/utils"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// PostgreSQLConnection func for connection to PostgreSQL database.
func PostgreSQLConnection() (*gorm.DB, error) {
	// Build PostgreSQL connection URL.
	connectionString, err := utils.ConnectionURLBuilder(os.Getenv("CONNECT_TYPE"))
	if err != nil {
		return nil, err
	}

	// New logger for detailed SQL logging
	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info, // Log level
			Colorful:      true,        // Enable color
		},
	)

	// connect
	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{
		Logger: newLogger, // add Logger
	})

	// AutoMigrate
	// db.AutoMigrate(&models.User{}, &models.Role{}, &models.Salary{}, &models.SalaryType{})
	// createRole(db)
	// createAdmin(db)

	// if connect fail
	if err != nil {
		panic("failed to connect to database")
	}

	return db, nil
}

func createRole(db *gorm.DB) {
	var count int64
	result := db.Model(&models.Role{}).Count(&count)
	if result.Error != nil || count > 0 {
		return
	}

	// add data
	roles := []*models.Role{
		{Name: "admin"},
		{Name: "user"},
	}
	for _, role := range roles {
		result := db.Model(&models.Role{}).Create(role)
		if result.Error != nil {
			log.Fatalf("Error creating roles: %v", result.Error)
		}
	}
}

func createAdmin(db *gorm.DB) {
	var count int64
	resCount := db.Model(&models.User{}).Count(&count)
	if resCount.Error != nil || count > 0 {
		return
	}

	admin := models.User{
		Email:    "admin@gmail.com",
		FullName: "แอดมิน (ผู้ตรวจสอบ)",
		Password: utils.GeneratePassword("admin"),
		RoleID:   1,
		Mobile:   "",
		TaxID:    "",
	}

	result := db.Model(&models.User{}).Create(&admin)
	if result.Error != nil {
		log.Fatalf("Error creating roles: %v", result.Error)
	}
}
