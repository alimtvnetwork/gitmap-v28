package cmdstatus

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
