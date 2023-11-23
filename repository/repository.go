// repository.go
package repository

import (
	"gorepository/model"

	"gorm.io/gorm"
)

type Repositories struct {
	UserRepo GenericRepository[model.User]
	PostRepo GenericRepository[model.Post]
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepo: NewGenericRepository[model.User](db),
		PostRepo: NewGenericRepository[model.Post](db),
	}
}
