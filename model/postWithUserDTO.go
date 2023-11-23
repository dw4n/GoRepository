package model

import (
	"time"

	"gorm.io/gorm"
)

type PostWithUserName struct {
	gorm.Model
	Title       string
	Content     string
	UserID      uint
	Published   bool
	PublishDate *time.Time
	UserName    string
}
