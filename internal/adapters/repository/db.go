package repository

import (
	"github.com/jinzhu/gorm"
)

// new database
func NewDB(db *gorm.DB) *gorm.DB {
	return db
}
