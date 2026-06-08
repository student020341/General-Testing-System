package labtest

import "context"

type Repository interface {
	GetByID(ctx context.Context, id string) (*Test, error)
	Save(ctx context.Context, test *Test) error
}
