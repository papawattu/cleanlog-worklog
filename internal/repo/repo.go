package repo

import "context"

type Repository[T any, S comparable] interface {
	Save(ctx context.Context, entity T) error
	Get(ctx context.Context, ID S) (T, error)
	GetAll(ctx context.Context) ([]T, error)
	Delete(ctx context.Context, e T) error
}
