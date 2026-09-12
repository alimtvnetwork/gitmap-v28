package cmd

import (
	"testing"
)

func TestParseStatusFlags_NoFlags(t *testing.T) {
	group, all, dirty := parseStatusFlags([]string{})
	if len(group) > 0 {
		t.Errorf("expected empty group, got %q", group)
	}

	if all {
		t.Error("expected all=false")
	}

	if dirty {
		t.Error("expected dirty=false")
	}
}

func TestParseStatusFlags_GroupLong(t *testing.T) {
	group, all, dirty := parseStatusFlags([]string{"--group", "backend"})
	if group != "backend" {
		t.Errorf("expected group=backend, got %q", group)
	}

	if all {
		t.Error("expected all=false")
	}

	if dirty {
		t.Error("expected dirty=false")
	}
}

func TestParseStatusFlags_GroupShort(t *testing.T) {
	group, all, dirty := parseStatusFlags([]string{"-g", "frontend"})
	if group != "frontend" {
		t.Errorf("expected group=frontend, got %q", group)
	}

	if all {
		t.Error("expected all=false")
	}

	if dirty {
		t.Error("expected dirty=false")
	}
}

func TestParseStatusFlags_All(t *testing.T) {
	group, all, dirty := parseStatusFlags([]string{"--all"})
	if len(group) > 0 {
		t.Errorf("expected empty group, got %q", group)
	}

	if all != true {
		t.Error("expected all=true")
	}

	if dirty {
		t.Error("expected dirty=false")
	}
}

func TestParseStatusFlags_GroupAndAll(t *testing.T) {
	group, all, dirty := parseStatusFlags([]string{"--group", "ops", "--all"})
	if group != "ops" {
		t.Errorf("expected group=ops, got %q", group)
	}

	if all != true {
		t.Error("expected all=true")
	}

	if dirty {
		t.Error("expected dirty=false")
	}
}

func TestParseStatusFlags_Dirty(t *testing.T) {
	_, _, dirty := parseStatusFlags([]string{"--dirty"})
	if !dirty {
		t.Error("expected dirty=true")
	}
}

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
