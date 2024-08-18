package database

import (
	"log"
	"maxl3oss/app/models"
	"maxl3oss/pkg/utils"
	"os"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// MariaDBConnection func for connection to MariaDB database.
func MariaDBConnection() (*gorm.DB, error) {
	// Build MariaDB connection URL.
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

	// Connect
	db, err := gorm.Open(mysql.Open(connectionString), &gorm.Config{
		Logger: newLogger, // Add Logger
	})

	// Set the time zone for the session
	if err := db.Exec("SET time_zone = 'Asia/Bangkok'").Error; err != nil {
		log.Fatalf("failed to set time zone: %v", err)
	}
	// AutoMigrate
	// db.Migrator().DropTable(&models.User{}, &models.Role{}, &models.Salary{}, &models.SalaryType{}, &models.SalaryOther{})
	// db.AutoMigrate(&models.User{}, &models.Role{}, &models.Salary{}, &models.SalaryType{}, &models.SalaryOther{})
	// createRole(db)
	// createAdmin(db)
	// createSalaryType(db)

	// If connection fails
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

	// Add data
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

func createSalaryType(db *gorm.DB) {
	var count int64
	response := db.Model(&models.SalaryType{}).Count(&count)
	if response.Error != nil || count > 0 {
		return
	}

	salaryTypes := []*models.SalaryType{
		{Name: "รพสต."},
		{Name: "สจ."},
		{Name: "ฝ่ายประจำ"},
		{Name: "เงินเดือนครู"},
		{Name: "บำนาญครู"},
		{Name: "บำเหน็จรายเดือน"},
		{Name: "บำนาญข้าราชการ"},
	}

	for _, salaryType := range salaryTypes {
		result := db.Model(&models.SalaryType{}).Create(salaryType)
		if result.Error != nil {
			log.Fatalf("Error creating salary type: %v", result.Error)
		}
	}
}
