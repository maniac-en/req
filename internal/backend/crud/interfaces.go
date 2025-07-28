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

type PaginationMetadata struct {
	Total       int64 `json:"total"`
	HasNext     bool  `json:"has_next"`
	HasPrev     bool  `json:"has_prev"`
	Limit       int   `json:"limit"`
	Offset      int   `json:"offset"`
	TotalPages  int   `json:"total_pages"`
	CurrentPage int   `json:"current_page"`
}

func CalculatePagination(total int64, limit, offset int) PaginationMetadata {
	totalPages := 1
	if total > 0 {
		totalPages = int((total + int64(limit) - 1) / int64(limit)) // Ceiling division
	}
	currentPage := (offset / limit) + 1
	hasNext := (offset + limit) < int(total)
	hasPrev := offset > 0

	return PaginationMetadata{
		Total:       total,
		HasNext:     hasNext,
		HasPrev:     hasPrev,
		Limit:       limit,
		Offset:      offset,
		TotalPages:  totalPages,
		CurrentPage: currentPage,
	}
}
