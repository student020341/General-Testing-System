package calculation

import (
	"context"
	"test-system/internal/domain/calculation"
	"test-system/internal/shared/paging"
	"testing"

	"go.uber.org/mock/gomock"
)

func TestCalculationsQuery(t *testing.T) {
	ctrl := gomock.NewController(t)
	mockCalcRepo := NewMockCalculationRepository(ctrl)

	ctx := context.Background()
	expectedSearch := calculation.Search{Name: "test-test", Paging: paging.NewPageRequest(1, 10)}
	expectedResult := []calculation.Calculation{
		{
			ID:     "foo",
			TestID: "bar",
			Name:   "test-test",
		},
	}

	mockCalcRepo.EXPECT().
		Search(ctx, expectedSearch).
		Return(expectedResult, nil)

	handler := NewListHandler(mockCalcRepo)
	list, err := handler.Handle(ctx, calculation.Search{Name: "test-test", Paging: paging.NewPageRequest(1, 10)})
	if err != nil {
		t.Fatalf("list handler: %v", err)
	}

	if len(list) != 1 {
		t.Fatalf("unexpected list length: %d", len(list))
	}

	if list[0].ID != "foo" {
		t.Fatalf("unexpected calculation id: %s", list[0].ID)
	}
}
