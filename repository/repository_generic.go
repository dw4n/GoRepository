package repository

import (
	"gorepository/publisher"
	"reflect"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type GenericRepository[T any] interface {
	GetAll() ([]T, error)
	GetWithConditions(result interface{}, conditions []func(*gorm.DB) *gorm.DB, opts ...GORMOption) error
	CountWithConditions(result *int64, conditions []func(*gorm.DB) *gorm.DB, opts ...GORMOption) error
	FindByID(id uuid.UUID) (T, error)
	Create(entity T, opts ...Option) (T, error)
	Update(entity T, opts ...Option) (T, error)
	Delete(id uuid.UUID, opts ...Option) error
}

type genericRepository[T any] struct {
	db        *gorm.DB
	publisher publisher.Publisher
}

func NewGenericRepository[T any](db *gorm.DB, publisher publisher.Publisher) GenericRepository[T] {
	return &genericRepository[T]{db, publisher}
}

func (r *genericRepository[T]) GetAll() ([]T, error) {
	var entities []T
	result := r.db.Find(&entities)
	return entities, result.Error
}

// hints : if gets error 'golang "reflect: reflect.Value.Set using unaddressable value"' add '&' in front result e.g. GetWithConditions(&result, etc...)
func (r *genericRepository[T]) GetWithConditions(result interface{}, conditions []func(*gorm.DB) *gorm.DB, gormOpts ...GORMOption) error {
	var entity T
	query := r.db.Model(&entity)

	for _, condition := range conditions {
		query = condition(query)
	}

	for _, opt := range gormOpts {
		query = opt(query)
	}

	return query.Find(result).Error
}
func (r *genericRepository[T]) CountWithConditions(result *int64, conditions []func(*gorm.DB) *gorm.DB, gormOpts ...GORMOption) error {
	var entity T
	query := r.db.Model(&entity)

	for _, condition := range conditions {
		query = condition(query)
	}

	for _, opt := range gormOpts {
		query = opt(query)
	}

	return query.Count(result).Error
}

func (r *genericRepository[T]) FindByID(id uuid.UUID) (T, error) {
	var entity T
	result := r.db.First(&entity, id)
	return entity, result.Error
}

func (r *genericRepository[T]) Create(entity T, opts ...Option) (T, error) {
	options := defaultOptions
	options.gormDB = r.db // Set the initial GORM DB

	for _, opt := range opts {
		opt(&options)
	}

	result := options.gormDB.Create(&entity)
	entityId := getIdFromEntity(&entity)
	if result.Error == nil && options.publish {
		r.publisher.PublishMessage(entity, "create", entityId) // create, no id have been formed yet
	}
	return entity, result.Error
}

func (r *genericRepository[T]) Update(entity T, opts ...Option) (T, error) {
	options := defaultOptions
	options.gormDB = r.db

	for _, opt := range opts {
		opt(&options)
	}

	result := options.gormDB.Save(&entity)
	entityId := getIdFromEntity(entity)
	if result.Error == nil && options.publish {
		r.publisher.PublishMessage(entity, "update", entityId)
	}
	return entity, result.Error
}

func (r *genericRepository[T]) Delete(id uuid.UUID, opts ...Option) error {
	options := defaultOptions
	options.gormDB = r.db

	for _, opt := range opts {
		opt(&options)
	}

	var entity T
	result := options.gormDB.Delete(&entity, id)
	if result.Error == nil && options.publish {
		r.publisher.PublishMessage(entity, "delete", id.String())
	}
	return result.Error
}

type Option func(*operationOptions)
type GORMOption func(*gorm.DB) *gorm.DB

type operationOptions struct {
	gormDB  *gorm.DB
	publish bool
}

var defaultOptions = operationOptions{
	publish: true, // By default, publishing is enabled
}

func WithPaging(page, pageSize int) GORMOption {
	return func(db *gorm.DB) *gorm.DB {
		offset := (page - 1) * pageSize
		return db.Offset(offset).Limit(pageSize)
	}
}

func WithSorting(columns []string) GORMOption {
	return func(db *gorm.DB) *gorm.DB {
		for _, column := range columns {
			db = db.Order(column)
		}
		return db
	}
}

func WithPublishing(publish bool) Option {
	return func(o *operationOptions) {
		o.publish = publish
	}
}

func WithPreload(s string) GORMOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Preload(s)
	}
}

func WithJoin(s string) GORMOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Joins(s)
	}
}
func WithSelect(s []string) GORMOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Select(s)
	}
}
func GroupBy(s string) GORMOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Group(s)
	}
}

func DistinctBy(s string) GORMOption {
	return func(db *gorm.DB) *gorm.DB {
		return db.Distinct(s)
	}
}

func getIdFromEntity[T any](entity T) string {
	val := reflect.ValueOf(entity)
	// If entity is a pointer, get the value it points to
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	// Check if the entity has an 'Id' field
	idField := val.FieldByName("Id")
	if !idField.IsValid() {
		// The entity does not have an 'Id' field
		return ""
	}

	// Check if the 'Id' field is of type uuid.UUID
	if idField.Type() == reflect.TypeOf(uuid.UUID{}) {
		// Return the 'Id' field's value as a string
		return idField.Interface().(uuid.UUID).String()
	}

	// If the 'Id' field is not of type uuid.UUID, or any other condition you wish to check
	return ""
}
