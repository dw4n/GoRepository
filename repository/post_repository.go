// post_repository.go
package repository

import (
	"gorepository/model"

	"gorm.io/gorm"
)

type PostRepository interface {
	Create(post model.Post) (model.Post, error)
	GetAll() ([]model.Post, error)
	FindByID(id uint) (model.Post, error)
	Update(post model.Post) (model.Post, error)
	Delete(id uint) error
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db}
}

func (r *userRepository) GetAll() ([]model.User, error) {
	var users []model.User
	result := r.db.Find(&users)
	return users, result.Error
}

func (r *userRepository) FindByID(id uint) (model.User, error) {
	var user model.User
	result := r.db.First(&user, id)
	return user, result.Error
}

func (r *userRepository) Create(user model.User) (model.User, error) {
	result := r.db.Create(&user)
	return user, result.Error
}

func (r *userRepository) Update(user model.User) (model.User, error) {
	result := r.db.Save(&user)
	return user, result.Error
}

func (r *userRepository) Delete(id uint) error {
	result := r.db.Delete(&model.User{}, id)
	return result.Error
}
