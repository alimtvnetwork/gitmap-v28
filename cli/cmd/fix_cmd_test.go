package cmd

import (
	"testing"
)

func TestIsFixAgyRequest_Variations(t *testing.T) {
	testCases := []struct {
		name     string
		args     []string
		expected bool
	}{
		{
			name:     "fix_agy",
			args:     []string{"agy"},
			expected: true,
		},
		{
			name:     "fix_errors_agy",
			args:     []string{"errors", "agy"},
			expected: true,
		},
		{
			name:     "fix_agy_errors",
			args:     []string{"agy", "errors"},
			expected: true,
		},
		{
			name:     "fix_aef",
			args:     []string{"aef"},
			expected: true,
		},
		{
			name:     "fix_agy_errors_fix",
			args:     []string{"agy-errors-fix"},
			expected: true,
		},
		{
			name:     "fix_normal_repo",
			args:     []string{"my-repo"},
			expected: false,
		},
		{
			name:     "empty_args",
			args:     []string{},
			expected: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := isFixAgyRequest(tc.args)
			if actual != tc.expected {
				t.Fatalf("expected %v for %v, got %v", tc.expected, tc.args, actual)
			}
		})
	}
}
