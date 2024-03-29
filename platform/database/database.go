package database

import "gorm.io/gorm"

type dbController struct {
	Db *gorm.DB
}
