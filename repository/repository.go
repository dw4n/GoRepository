package repository

import (
	"gorepository/model"

	"gorm.io/gorm"
)

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
