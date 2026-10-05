package cmdautomation

import (
	"testing"
)

func TestCleanSearchPattern(t *testing.T) {
	testCases := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "escaped double quotes",
			input:    `\"all\"`,
			expected: "all",
		},
		{
			name:     "wrapped escaped double quotes",
			input:    `"\"all\""`,
			expected: "all",
		},
		{
			name:     "single quotes",
			input:    `'all'`,
			expected: "all",
		},
		{
			name:     "double quoted single quotes",
			input:    `"'all'"`,
			expected: "all",
		},
		{
			name:     "escaped double quotes phrase",
			input:    `\"import os\"`,
			expected: "import os",
		},
		{
			name:     "wrapped escaped double quotes phrase",
			input:    `"\"import os\""`,
			expected: "import os",
		},
		{
			name:     "normal string without quotes",
			input:    "normal",
			expected: "normal",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "nested escaped quotes",
			input:    `\"\"nested\"\"`,
			expected: "nested",
		},
		{
			name:     "backtick quotes",
			input:    "`backticks`",
			expected: "backticks",
		},
		{
			name:     "internal escaped quotes unescaped",
			input:    `hello\"world`,
			expected: `hello"world`,
		},
		{
			name:     "single character quote preserved",
			input:    `"`,
			expected: `"`,
		},
		{
			name:     "single character single quote preserved",
			input:    `'`,
			expected: `'`,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			actual := cleanSearchPattern(tc.input)
			if actual != tc.expected {
				t.Errorf("cleanSearchPattern(%q) = %q, want %q", tc.input, actual, tc.expected)
			}
		})
	}
}

func TestSearchWorkerCleanPattern(t *testing.T) {
	opts := SearchOptions{
		Pattern: `"\"import os\""`,
	}
	ctx := buildWorkerContext(opts)
	if string(ctx.patBytes) != "import os" {
		t.Errorf("expected ctx.patBytes to be %q, got %q", "import os", string(ctx.patBytes))
	}
}
