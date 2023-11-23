// repository.go
package repository

import (
	"gorm.io/gorm"
)

// Repositories holds references to all repositories
type Repositories struct {
	UserRepo UserRepository
	PostRepo PostRepository
}

// NewRepositories initializes all repositories
func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepo: NewUserRepository(db),
		PostRepo: NewPostRepository(db),
	}
}
