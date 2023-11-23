package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name      string
	Email     string `gorm:"unique"`
	IsDeleted bool
	// Add other fields as needed
}
