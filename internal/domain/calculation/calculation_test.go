package calculation

import (
	"fmt"
	"slices"
	"strings"
	"testing"
)

func TestCalculationClosureParsing(t *testing.T) {
	tests := []struct {
		name     string
		src      string
		expected closureDetails
	}{
		// simple
		{
			name:     "no parameters, no return",
			src:      `() => {}`,
			expected: closureDetails{},
		},
		{
			name:     "no parameters, single return",
			src:      `() => 1`,
			expected: closureDetails{HasSingleReturn: true},
		},
		{
			name:     "no parameters, array return",
			src:      `() => ([])`,
			expected: closureDetails{HasArrayReturn: true},
		},
		{
			name:     "no parameters, array return 2",
			src:      `() => { return []; }`,
			expected: closureDetails{HasArrayReturn: true},
		},
		{
			name:     "no parameters, keyed returned",
			src:      `() => ({a: 1, "b": "f"})`,
			expected: closureDetails{KeyedReturnFields: []string{"a", "b"}},
		},
		{
			name: "no parameters, variable return (single)",
			src: `() => {
				var x = {};
				return x;
			}`,
			expected: closureDetails{HasSingleReturn: true},
		},
		// more complicated
		{
			name: "parameters and all returns",
			src: `(width, height, length) => {
				if (width == 0) {
					return "fail";
				} else if (height == 0) {
					return ["fail"];
				} else if (length == 0) {
					return { error: "fail" };
				}
				
				return width * height * length;
			}`,
			expected: closureDetails{
				Parameters:      []string{"width", "height", "length"},
				HasSingleReturn: true,
				HasArrayReturn:  true,
				KeyedReturnFields: []string{
					"error",
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(*testing.T) {
			cd, err := parseAndWalk(tt.src)
			if err != nil {
				t.Errorf("parsing source: %v", err)
			}

			if err := compareDetails(*cd, tt.expected); err != nil {
				t.Error(err)
			}
		})
	}

}

func compareDetails(a, b closureDetails) error {
	if a.HasSingleReturn != b.HasSingleReturn {
		return fmt.Errorf(
			"expected a.HasSingleReturn(%v) to equal b.HasSingleReturn(%v)",
			a.HasSingleReturn,
			b.HasSingleReturn,
		)
	}

	if a.HasArrayReturn != b.HasArrayReturn {
		return fmt.Errorf(
			"expected a.HasArrayReturn(%v) to equal b.HasArrayReturn(%v)",
			a.HasArrayReturn,
			b.HasArrayReturn,
		)
	}

	if len(a.KeyedReturnFields) != len(b.KeyedReturnFields) {
		return fmt.Errorf(
			"a.KeyedReturnFields(%d) has a different number of fields than b.KeyedReturnFields(%d)",
			len(a.KeyedReturnFields),
			len(b.KeyedReturnFields),
		)
	}

	slices.Sort(a.KeyedReturnFields)
	slices.Sort(b.KeyedReturnFields)

	if slices.Compare(a.KeyedReturnFields, b.KeyedReturnFields) != 0 {
		return fmt.Errorf(
			"KeyedReturnFields values are different between a and b\na(%s)\nb(%s)",
			strings.Join(a.KeyedReturnFields, ", "),
			strings.Join(b.KeyedReturnFields, ", "),
		)
	}

	if len(a.Parameters) != len(b.Parameters) {
		return fmt.Errorf(
			"a.Parameters(%d) has a different number of fields than b.Parameters(%d)",
			len(a.Parameters),
			len(b.Parameters),
		)
	}

	slices.Sort(a.Parameters)
	slices.Sort(b.Parameters)

	if slices.Compare(a.Parameters, b.Parameters) != 0 {
		return fmt.Errorf(
			"Parameters values are different between a and b\na(%s)\nb(%s)",
			strings.Join(a.Parameters, ", "),
			strings.Join(b.Parameters, ", "),
		)
	}

	return nil
}
