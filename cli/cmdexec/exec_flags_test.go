package cmdexec

import (
	"testing"
)

func TestParseExecFlags_NoFlags(t *testing.T) {
	group, all, _, gitArgs := parseExecFlags([]string{"fetch", "--prune"})
	if len(group) > 0 {
		t.Errorf("expected empty group, got %q", group)
	}

	if all {
		t.Error("expected all=false")
	}

	if len(gitArgs) != 2 || gitArgs[0] != "fetch" || gitArgs[1] != "--prune" {
		t.Errorf("expected [fetch --prune], got %v", gitArgs)
	}
}

func TestParseExecFlags_GroupLong(t *testing.T) {
	group, all, _, gitArgs := parseExecFlags([]string{"--group", "backend", "status"})
	if group != "backend" {
		t.Errorf("expected group=backend, got %q", group)
	}

	if all {
		t.Error("expected all=false")
	}

	if len(gitArgs) != 1 || gitArgs[0] != "status" {
		t.Errorf("expected [status], got %v", gitArgs)
	}
}

func TestParseExecFlags_GroupShort(t *testing.T) {
	group, _, _, gitArgs := parseExecFlags([]string{"-g", "infra", "pull"})
	if group != "infra" {
		t.Errorf("expected group=infra, got %q", group)
	}

	if len(gitArgs) != 1 || gitArgs[0] != "pull" {
		t.Errorf("expected [pull], got %v", gitArgs)
	}
}

func TestParseExecFlags_All(t *testing.T) {
	group, all, _, gitArgs := parseExecFlags([]string{"--all", "fetch"})
	if len(group) > 0 {
		t.Errorf("expected empty group, got %q", group)
	}

	if all != true {
		t.Error("expected all=true")
	}

	if len(gitArgs) != 1 || gitArgs[0] != "fetch" {
		t.Errorf("expected [fetch], got %v", gitArgs)
	}
}

func TestParseExecFlags_NoArgs(t *testing.T) {
	group, all, _, gitArgs := parseExecFlags([]string{"--all"})
	if len(group) > 0 {
		t.Errorf("expected empty group, got %q", group)
	}

	if all != true {
		t.Error("expected all=true")
	}

	if len(gitArgs) != 0 {
		t.Errorf("expected empty gitArgs, got %v", gitArgs)
	}
}
