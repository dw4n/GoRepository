package repository

import "gorepository/model"

type UserRepository interface {
	GetAll() ([]model.User, error)
	FindByID(id uint) (model.User, error)
	Create(user model.User) (model.User, error)
	Update(user model.User) (model.User, error)
	Delete(id uint) error
}
