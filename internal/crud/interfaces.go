// Package crud provides generic CRUD operations for entities.
package crud

import (
	"context"
	"errors"
	"time"
)

var (
	ErrNotFound     = errors.New("entity not found")
	ErrInvalidInput = errors.New("invalid input")
)

type Entity interface {
	GetID() int64
	GetName() string
	GetCreatedAt() time.Time
	GetUpdatedAt() time.Time
}

type Manager[T Entity] interface {
	Create(ctx context.Context, name string) (T, error)
	Read(ctx context.Context, id int64) (T, error)
	Update(ctx context.Context, id int64, name string) (T, error)
	Delete(ctx context.Context, id int64) error
	List(ctx context.Context) ([]T, error)
}
