// Package cmd — agy_suite_v2_test.go tests new AGY commands, status gap, and prefix exceptions.
package cmd

import (
	"strings"
	"testing"
)

func TestFormatAgyStatusGap(t *testing.T) {
	active := formatAgyStatus("active", false, 14)
	if !strings.Contains(active, "✔   active") {
		t.Errorf("expected gap between check and active, got: %q", active)
	}

	missing := formatAgyStatus("missing", true, 14)
	if !strings.Contains(missing, "✖   missing") {
		t.Errorf("expected gap between cross and missing, got: %q", missing)
	}

	global := formatAgyStatus("global", false, 14)
	if !strings.Contains(global, "—   global") {
		t.Errorf("expected gap in global status, got: %q", global)
	}
}

func TestNormalizeAgySubcommand(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"fdp", "find-duplicate-projects"},
		{"find-duplicate-projects", "find-duplicate-projects"},
		{"cdp", "optimize-projects"},
		{"cure-duplicate-projects", "optimize-projects"},
		{"remove-misisng-projects", "remove-missing-projects"},
		{"rm-missing", "remove-missing-projects"},
		{"aprmp", "all-projects-read-memory-prompt"},
		{"recon", "reconcile"},
	}

	for _, tc := range cases {
		got := normalizeAgySubcommand(tc.input)
		if got != tc.want {
			t.Errorf("normalizeAgySubcommand(%q) = %q, want %q", tc.input, got, tc.want)
		}
	}
}

func TestMatchPrefixOrSlugExcept(t *testing.T) {
	proj := AgyProject{
		ID:   "46d05021-30e1-4036-acfb-3020489125eb",
		Name: "gitmap-v28",
		ProjectResources: &AgyProjectResources{
			Resources: []AgyResource{
				{GitFolder: &AgyGitFolder{FolderURI: "file:///mock/repos/gitmap"}},
			},
		},
	}

	if !isMatchPrefixOrSlugExcept(proj, []string{"46d0"}) {
		t.Errorf("expected prefix match on short ID 46d0")
	}

	if !isMatchPrefixOrSlugExcept(proj, []string{"gitmap"}) {
		t.Errorf("expected prefix match on slug/name gitmap")
	}

	if isMatchPrefixOrSlugExcept(proj, []string{"other-project"}) {
		t.Errorf("expected no match on other-project")
	}
}
