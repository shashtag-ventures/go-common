package gormutil

import (
	"context"
	"errors"
	"fmt"

	"github.com/shashtag-ventures/go-common/middleware"
	"gorm.io/gorm"
)

// Repository is a generic GORM repository that provides basic CRUD operations.
type Repository[T any] struct {
	db *gorm.DB
}

// NewRepository creates a new instance of the generic repository for a given model type.
func NewRepository[T any](db *gorm.DB) *Repository[T] {
	return &Repository[T]{db: db}
}

// Create creates a new entity in the database.
// It uses the provided transaction `tx` if not nil, otherwise it uses the base DB connection.
func (r *Repository[T]) Create(ctx context.Context, tx *gorm.DB, entity *T) error {
	logger := middleware.GetLoggerFromContext(ctx)
	db := r.getDB(tx)

	if err := db.WithContext(ctx).Create(entity).Error; err != nil {
		logger.Error("Failed to create entity in DB", "error", err)
		return fmt.Errorf("failed to create entity: %w", err)
	}
	return nil
}

// FindByID finds an entity by its primary key.
func (r *Repository[T]) FindByID(ctx context.Context, id any) (*T, error) {
	logger := middleware.GetLoggerFromContext(ctx)
	var entity T
	if err := r.db.WithContext(ctx).First(&entity, id).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("Failed to find entity by ID in DB", "id", id, "error", err)
		}
		return nil, fmt.Errorf("failed to find entity by ID %v: %w", id, err)
	}
	return &entity, nil
}

// Update updates an existing entity in the database.
// It uses the provided transaction `tx` if not nil, otherwise it uses the base DB connection.
func (r *Repository[T]) Update(ctx context.Context, tx *gorm.DB, entity *T) error {
	logger := middleware.GetLoggerFromContext(ctx)
	db := r.getDB(tx)

	if err := db.WithContext(ctx).Save(entity).Error; err != nil {
		logger.Error("Failed to update entity in DB", "error", err)
		return fmt.Errorf("failed to update entity: %w", err)
	}
	return nil
}

// Delete deletes an entity from the database.
// It uses the provided transaction `tx` if not nil, otherwise it uses the base DB connection.
func (r *Repository[T]) Delete(ctx context.Context, tx *gorm.DB, entity *T) error {
	logger := middleware.GetLoggerFromContext(ctx)
	db := r.getDB(tx)

	if err := db.WithContext(ctx).Delete(entity).Error; err != nil {
		logger.Error("Failed to delete entity from DB", "error", err)
		return fmt.Errorf("failed to delete entity: %w", err)
	}
	return nil
}

// FindOneBy is a generic finder for a single record.
func (r *Repository[T]) FindOneBy(ctx context.Context, query string, args ...any) (*T, error) {
	logger := middleware.GetLoggerFromContext(ctx)
	var entity T
	if err := r.db.WithContext(ctx).Where(query, args...).First(&entity).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("Failed to find entity in DB", "query", query, "args", args, "error", err)
		}
		return nil, fmt.Errorf("failed to find entity: %w", err)
	}
	return &entity, nil
}

// Find is a generic finder for multiple records.
func (r *Repository[T]) Find(ctx context.Context, query string, args ...any) ([]*T, error) {
	logger := middleware.GetLoggerFromContext(ctx)
	var entities []*T
	if err := r.db.WithContext(ctx).Where(query, args...).Find(&entities).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("Failed to find entities in DB", "query", query, "args", args, "error", err)
		}
		return nil, fmt.Errorf("failed to find entities: %w", err)
	}
	return entities, nil
}

// FindAll finds all entities of a given type.
func (r *Repository[T]) FindAll(ctx context.Context) ([]*T, error) {
	logger := middleware.GetLoggerFromContext(ctx)
	var entities []*T
	if err := r.db.WithContext(ctx).Find(&entities).Error; err != nil {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error("Failed to find all entities in DB", "error", err)
		}
		return nil, fmt.Errorf("failed to find all entities: %w", err)
	}
	return entities, nil
}

// DB returns the underlying gorm.DB instance with the provided context.
func (r *Repository[T]) DB(ctx context.Context) *gorm.DB {
	return r.db.WithContext(ctx)
}

// getDB returns the transaction if it's not nil, otherwise it returns the base DB connection.
func (r *Repository[T]) getDB(tx *gorm.DB) *gorm.DB {
	if tx != nil {
		return tx
	}
	return r.db
}