package cmd

import (
	"testing"
)

func TestNormalizeVersion(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"empty string", "", ""},
		{"already normal", "v1.0.0", "v1.0.0"},
		{"missing v", "1.0.0", "v1.0.0"},
		{"capital V", "V1.0.0", "v1.0.0"},
		{"multiple v", "vv1.0.0", "v1.0.0"},
		{"multiple V", "VV1.0.0", "v1.0.0"},
		{"mixed vV", "vVv1.0.0", "v1.0.0"},
		{"only v", "v", "v"},
		{"only V", "V", "v"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeVersion(tt.input)
			if result != tt.expected {
				t.Errorf("normalizeVersion(%q) = %q; want %q", tt.input, result, tt.expected)
			}
		})
	}
}
