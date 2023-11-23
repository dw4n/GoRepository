// post_repository.go
package repository

import "gorepository/model"

type PostRepository interface {
	Create(post model.Post) (model.Post, error)
	GetAll() ([]model.Post, error)
	FindByID(id uint) (model.Post, error)
	Update(post model.Post) (model.Post, error)
	Delete(id uint) error
}
