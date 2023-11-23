package model

import "gorm.io/gorm"

type User struct {
	gorm.Model
	Name  string
	Email string `gorm:"unique"`
	// Add other fields as needed
}

// TableName sets the insert table name for this struct type
func (u *User) TableName() string {
	return "user"
}
