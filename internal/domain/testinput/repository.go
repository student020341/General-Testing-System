package testinput

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*TestInput, error)
	Save(ctx context.Context, testInput *TestInput) error
	Delete(ctx context.Context, testInput *TestInput) error
}
