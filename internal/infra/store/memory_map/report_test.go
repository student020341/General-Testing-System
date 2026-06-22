package memorymap

import (
	"context"
	"fmt"
	"test-system/internal/domain/report"
	"test-system/internal/shared/paging"
	"testing"
)

func TestReportRepoBasics(t *testing.T) {
	repo := NewReportRepository()
	ctx := context.Background()

	if err := repo.Save(ctx, &report.Report{
		ID:   "r-1",
		Name: "test3",
	}); err != nil {
		t.Fatalf("saving report: %v", err)
	}

	e, err := repo.GetByID(ctx, "r-1")
	if err != nil {
		t.Fatalf("fetching report: %v", err)
	}

	if e.ID != "r-1" || e.Name != "test3" {
		t.Fatalf("unexpected report: %s, %s", e.ID, e.Name)
	}

	if err := repo.Delete(ctx, e); err != nil {
		t.Fatalf("delete report: %v", err)
	}

	// quick paging test
	s := report.Search{Paging: paging.NewPageRequest(1, 2)}
	list, err := repo.Search(ctx, s)
	if err != nil {
		t.Fatalf("empty collection search: %v", err)
	}

	if len(list) != 0 {
		t.Fatalf("collection should be empty")
	}

	for i := 0; i < 5; i++ {
		status := report.StatusOpen
		if i%2 == 0 {
			status = report.StatusClosed
		}
		repo.Save(ctx, &report.Report{
			ID:     fmt.Sprintf("r-%d", i),
			Name:   fmt.Sprintf("Test %d", i),
			Status: status,
		})
	}

	s.Paging.PageSize = 10
	list, err = repo.Search(ctx, s)
	if err != nil {
		t.Fatalf("search 10: %v", err)
	}
	if len(list) != 5 {
		t.Fatalf("should have entire collection: %d", len(list))
	}

	// fetch page 1
	s.Paging.PageSize = 2
	list, err = repo.Search(ctx, s)
	if err != nil {
		t.Fatalf("search page 1 size 2: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("unexpected page size: %d", len(list))
	}

	if !(list[0].ID == "r-0" && list[1].ID == "r-1") {
		t.Fatalf("unexpected page 1 results: %s, %s", list[0].ID, list[1].ID)
	}

	// fetch page 2
	s.Paging.Page = 2
	list, err = repo.Search(ctx, s)
	if err != nil {
		t.Fatalf("search page 2 size 2: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("unexpected page size: %d", len(list))
	}

	if !(list[0].ID == "r-2" && list[1].ID == "r-3") {
		t.Fatalf("unexpected page 2 results: %s, %s", list[0].ID, list[1].ID)
	}

	// fetch page 3
	s.Paging.Page = 3
	list, err = repo.Search(ctx, s)
	if err != nil {
		t.Fatalf("search page 3 size 2: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("unexpected page size: %d", len(list))
	}

	if !(list[0].ID == "r-4") {
		t.Fatalf("unexpected page 2 results: %s", list[0].ID)
	}

	// fetch page 4
	s.Paging.Page = 4
	list, err = repo.Search(ctx, s)
	if err != nil {
		t.Fatalf("search page 4, size 2: %v", err)
	}

	if len(list) != 0 {
		t.Fatalf("unexpected page size: %d", len(list))
	}

	// status search
	s.Paging.Page = 1
	s.Paging.PageSize = 10
	s.Status = report.StatusOpen
	list, err = repo.Search(ctx, s)
	if err != nil {
		t.Fatalf("status search: %v", err)
	}

	if len(list) != 2 {
		t.Fatalf("unexpected status response: %d", len(list))
	}

	if !(list[0].ID == "r-1" && list[1].ID == "r-3") {
		t.Fatalf("unexpected status response: %s, %s", list[0].ID, list[1].ID)
	}
}
