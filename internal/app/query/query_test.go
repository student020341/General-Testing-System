//go:generate mockgen -destination=mock_calc_test.go -package=query -mock_names Repository=MockCalculationRepository test-system/internal/domain/calculation Repository

package query
