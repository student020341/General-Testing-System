//go:generate mockgen -destination=mock_calc_test.go -package=calculation -mock_names Repository=MockCalculationRepository test-system/internal/domain/calculation Repository

package calculation
