package memorymap

import (
	"context"
	"test-system/internal/domain/labtest"
	"testing"
)

func TestLabTestRepoBasics(t *testing.T) {
	repo := NewTestRepository()
	ctx := context.Background()

	if err := repo.Save(ctx, &labtest.Test{
		ID:   "t-1",
		Name: "test2",
	}); err != nil {
		t.Fatalf("saving test: %v", err)
	}

	e, err := repo.GetByID(ctx, "t-1")
	if err != nil {
		t.Fatalf("fetching test: %v", err)
	}

	if e.ID != "t-1" || e.Name != "test2" {
		t.Fatalf("unexpected test: %s, %s", e.ID, e.Name)
	}

	// TODO some calculation specific stuff
}
