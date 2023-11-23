# Golang Repository
This is a simple POC of Golang Repository Pattern
Prepare your postgre, just connect!

##
How to use
```
1. The model in model folder
2. register your model in repository.go
3. Most of query, pagination, sort is already there
4. In case of weird stuff, you can always use raw query
```

## Generic Repository
This is using Generic Repository.
If you need to create any custom specific to your model, you can do this
```
// repository.go
package repository

import (
	"gorepository/model"

	"gorm.io/gorm"
)

type Repositories struct {
	UserRepo       GenericRepository[model.User]
	PostRepo       GenericRepository[model.Post]
	PostRepoCustom PostRepository
}

func NewRepositories(db *gorm.DB) *Repositories {
	return &Repositories{
		UserRepo:       NewGenericRepository[model.User](db),
		PostRepo:       NewGenericRepository[model.Post](db),
		PostRepoCustom: NewPostRepository(db),
	}
}

```

And put your custom post_repository.go here
```
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

```