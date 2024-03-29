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
	connectionString, err := utils.ConnectionURLBuilder("postgres")
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
	db.AutoMigrate(&models.User{}, &models.Role{}, &models.Salary{})

	// if connect fail
	if err != nil {
		panic("failed to connect to database")
	}

	return db, nil
}

// # setup first build database
// func migrationDB(db *gorm.DB) {
// 	db.AutoMigrate(&models.Prefix{}, &models.User{})
// 	fmt.Println("Database migration completed!")
// }

// func createPrefixes(db *gorm.DB) error {
// 	// add data
// 	prefixes := []*models.Prefix{
// 		{TitleTh: "นาย", TitleEn: "Mr."},
// 		{TitleTh: "นาง", TitleEn: "Mrs."},
// 		{TitleTh: "นางสาว", TitleEn: "Miss"},
// 	}
// 	for _, prefix := range prefixes {
// 		result := db.Create(prefix)
// 		if result.Error != nil {
// 			log.Fatalf("Error creating prefix: %v", result.Error)
// 		}
// 	}
// 	fmt.Printf("Create prefix successfully!")
// 	return nil
// }

// func createRole(db *gorm.DB) error {
// 	// add data
// 	roles := []*models.Role{
// 		{Name: "admin"},
// 		{Name: "user"},
// 	}
// 	for _, role := range roles {
// 		result := db.Model(&models.Role{}).Create(role)
// 		if result.Error != nil {
// 			log.Fatalf("Error creating roles: %v", result.Error)
// 		}
// 	}
// 	fmt.Printf("Create roles successfully!")
// 	return nil
// }
