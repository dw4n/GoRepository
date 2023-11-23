// post.go
package model

import (
	"time"

	"gorm.io/gorm"
)

type Post struct {
	gorm.Model
	Title       string
	Content     string
	UserID      uint // Foreign key for User
	Published   bool
	PublishDate *time.Time
}

// BeforeCreate will set a record's publish date to now
func (post *Post) BeforeCreate(tx *gorm.DB) (err error) {
	if post.Published && post.PublishDate == nil {
		now := time.Now()
		post.PublishDate = &now
	}
	return
}
