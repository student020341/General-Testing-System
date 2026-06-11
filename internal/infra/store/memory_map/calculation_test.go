package memorymap

import (
	"context"
	"test-system/internal/domain/calculation"
	"testing"
)

func TestCalculationRepoBasics(t *testing.T) {
	repo := NewCalculationRepository()
	ctx := context.Background()

	if err := repo.Save(ctx, &calculation.Calculation{}); err != ErrInvalidID {
		t.Fatal("blank calculation ID should be an error")
	}

	if _, err := repo.GetByID(ctx, "not-found"); err != ErrNotFound {
		t.Fatal("calculation should not be found")
	}

	if err := repo.Save(ctx, &calculation.Calculation{
		ID:   "c-1",
		Name: "test",
	}); err != nil {
		t.Fatalf("saving calculation: %v", err)
	}

	c, err := repo.GetByID(ctx, "c-1")
	if err != nil {
		t.Fatalf("fetching calculation: %v", err)
	}

	if c.ID != "c-1" || c.Name != "test" {
		t.Fatalf("unexpected calculation: %s, %s", c.ID, c.Name)
	}

	if err := repo.Delete(ctx, c); err != nil {
		t.Fatalf("should delete: %v", err)
	}

	if err := repo.Delete(ctx, c); err != nil {
		t.Fatalf("deleting the same entity should be ok: %v", err)
	}

	c, err = repo.GetByID(ctx, "c-1")
	if c != nil {
		t.Fatalf("calculation should be nil")
	}

	if err != ErrNotFound {
		t.Fatalf("should have not found error: %v", err)
	}
}
