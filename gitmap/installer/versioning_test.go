// Package installer — versioning_test.go tests semantic versioning helpers.
package installer

import (
	"testing"
)

func TestSemanticVersioning(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"v1.0.0", "v1.0.1"},
		{"1.0.0", "1.0.1"},
		{"v2.4.9", "v2.4.10"},
		{"invalid", "1.0.1"},
		{"vinvalid", "v1.0.1"},
	}

	for _, c := range cases {
		out := NextSemanticVersion(c.input)
		if out != c.expected {
			t.Errorf("NextSemanticVersion(%q) = %q, expected %q", c.input, out, c.expected)
		}
	}
}
