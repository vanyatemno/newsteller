package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Repository[T any] interface {
	All(ctx context.Context) ([]T, error)
	FindByID(ctx context.Context, id string) (*T, error)
	FindPaginated(ctx context.Context, query *PaginatedSearchQuery) ([]T, int64, error)
	Create(ctx context.Context, model *T) (*primitive.ObjectID, error)
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, model *T) error
}

type PaginatedSearchQuery struct {
	Page    int    `json:"page" valid:"required,gte=1"`
	Limit   int    `json:"limit" validate:"required,gte=1"`
	Keyword string `json:"keyword"`
}
