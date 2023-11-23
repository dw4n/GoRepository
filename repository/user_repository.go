package repository

import (
	"gorepository/model"

	"gorm.io/gorm"
)

type UserRepository interface {
	GetAll() ([]model.User, error)
	FindByID(id uint) (model.User, error)
	Create(user model.User) (model.User, error)
	Update(user model.User) (model.User, error)
	Delete(id uint) error
}

type postRepository struct {
	db *gorm.DB
}

func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{db}
}

func (r *postRepository) Create(post model.Post) (model.Post, error) {
	result := r.db.Create(&post)
	return post, result.Error
}

func (r *postRepository) GetAll() ([]model.Post, error) {
	var posts []model.Post
	result := r.db.Find(&posts)
	return posts, result.Error
}

func (r *postRepository) FindByID(id uint) (model.Post, error) {
	var post model.Post
	result := r.db.First(&post, id)
	return post, result.Error
}

func (r *postRepository) Update(post model.Post) (model.Post, error) {
	result := r.db.Save(&post)
	return post, result.Error
}

func (r *postRepository) Delete(id uint) error {
	result := r.db.Delete(&model.Post{}, id)
	return result.Error
}
