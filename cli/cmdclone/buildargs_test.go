package cmdclone

import (
	"testing"
)

func TestBuildCommandArgs_JoinsArgs(t *testing.T) {
	result := buildCommandArgs([]string{"scan", "--all", "/repos"})
	if result != "scan --all /repos" {
		t.Errorf("expected 'scan --all /repos', got %q", result)
	}
}

func TestBuildCommandArgs_Empty(t *testing.T) {
	result := buildCommandArgs([]string{})
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestBuildCommandArgs_SingleArg(t *testing.T) {
	result := buildCommandArgs([]string{"pull"})
	if result != "pull" {
		t.Errorf("expected 'pull', got %q", result)
	}
}

// --- findDuplicate ---
