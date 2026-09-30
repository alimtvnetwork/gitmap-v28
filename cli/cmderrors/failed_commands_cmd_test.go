package cmderrors

import (
	"testing"
)

func TestIsFailedSubcommandToken(t *testing.T) {
	cases := []struct {
		token    string
		expected bool
	}{
		{"failed", true},
		{"failed-commands", true},
		{"fc", true},
		{"unknown-commands", true},
		{"failed-to-detect", true},
		{"clear", false},
		{"show", false},
	}

	for _, tc := range cases {
		got := isFailedSubcommandToken(tc.token)
		if got != tc.expected {
			t.Errorf("isFailedSubcommandToken(%q) = %v, want %v", tc.token, got, tc.expected)
		}
	}
}

func TestIsFailedCountRequested(t *testing.T) {
	if !isFailedCountRequested([]string{"count"}) {
		t.Fatal("expected count to be detected")
	}
	if !isFailedCountRequested([]string{"--json", "stats"}) {
		t.Fatal("expected stats to be detected")
	}
	if isFailedCountRequested([]string{"--json"}) {
		t.Fatal("expected --json alone not to trigger count mode")
	}
}
