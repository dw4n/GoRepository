package repository

import (
	"gorm.io/gorm"
)

type GenericRepository[T any] interface {
	GetAll() ([]T, error)
	GetAllWithConditions(result interface{}, conditions []func(*gorm.DB) *gorm.DB, opts ...Option) error
	FindByID(id uint) (T, error)
	Create(entity T) (T, error)
	Update(entity T) (T, error)
	Delete(id uint) error
}

type genericRepository[T any] struct {
	db *gorm.DB
}

func NewGenericRepository[T any](db *gorm.DB) GenericRepository[T] {
	return &genericRepository[T]{db}
}

func (r *genericRepository[T]) GetAll() ([]T, error) {
	var entities []T
	result := r.db.Find(&entities)
	return entities, result.Error
}

func (r *genericRepository[T]) GetAllWithConditions(result interface{}, conditions []func(*gorm.DB) *gorm.DB, opts ...Option) error {
	var entity T
	query := r.db.Model(&entity)

	for _, condition := range conditions {
		query = condition(query)
	}

	for _, opt := range opts {
		query = opt(query)
	}

	return query.Find(result).Error
}

func (r *genericRepository[T]) FindByID(id uint) (T, error) {
	var entity T
	result := r.db.First(&entity, id)
	return entity, result.Error
}

func (r *genericRepository[T]) Create(entity T) (T, error) {
	result := r.db.Create(&entity)
	return entity, result.Error
}

func (r *genericRepository[T]) Update(entity T) (T, error) {
	result := r.db.Save(&entity)
	return entity, result.Error
}

func (r *genericRepository[T]) Delete(id uint) error {
	var entity T
	result := r.db.Delete(&entity, id)
	return result.Error
}

type Option func(*gorm.DB) *gorm.DB

func WithPaging(page, pageSize int) Option {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func WithSorting(columns []string) Option {
	return func(db *gorm.DB) *gorm.DB {
		for _, column := range columns {
			db = db.Order(column)
		}
		return db
	}
}
