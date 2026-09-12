package cmd

import (
	"testing"
)

func TestDeriveSlug_Integration(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{"github.com/org/repo", "github-com-org-repo"},
		{"go.example.dev/pkg", "go-example-dev-pkg"},
		{"github.com/a/b@v2", "github-com-a-b-v2"},
		{"simple", "simple"},
	}

	for _, tt := range tests {
		got := deriveSlug(tt.input)
		if got != tt.expected {
			t.Errorf("deriveSlug(%q) = %q, want %q", tt.input, got, tt.expected)
		}
	}
}
