package state

import (
	"context"
	"newsteller/internal/models"
	"newsteller/internal/repositories"
)

type State[T models.Model] interface {
	FindAll(ctx context.Context) ([]T, error)
	FindPaginated(ctx context.Context, q *repositories.PaginatedSearchQuery) ([]T, int64, error)
	FindByID(ctx context.Context, id string) (*T, error)
	Insert(ctx context.Context, model *T) error
	Delete(ctx context.Context, id string) error
	Update(ctx context.Context, model *T) error
}
