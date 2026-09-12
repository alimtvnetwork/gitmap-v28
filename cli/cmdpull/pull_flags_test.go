package cmdpull

import (
	"testing"
)

func TestParsePullFlags_NoFlags(t *testing.T) {
	opts := parsePullFlags([]string{"my-repo"})
	if opts.slug != "my-repo" {
		t.Errorf("expected slug=my-repo, got %q", opts.slug)
	}

	if len(opts.group) > 0 || opts.all || opts.verbose {
		t.Error("expected no group/all/verbose")
	}

	if opts.parallel != 0 {
		t.Errorf("expected default parallel=0, got %d", opts.parallel)
	}

	if opts.onlyAvailable {
		t.Error("expected onlyAvailable=false by default")
	}
}

func TestParsePullFlags_GroupLong(t *testing.T) {
	opts := parsePullFlags([]string{"--group", "backend"})
	if len(opts.slug) > 0 {
		t.Errorf("expected empty slug, got %q", opts.slug)
	}

	if opts.group != "backend" {
		t.Errorf("expected group=backend, got %q", opts.group)
	}

	if opts.all {
		t.Error("expected all=false")
	}
}

func TestParsePullFlags_GroupShort(t *testing.T) {
	opts := parsePullFlags([]string{"-g", "infra"})
	if opts.group != "infra" {
		t.Errorf("expected group=infra, got %q", opts.group)
	}
}

func TestParsePullFlags_All(t *testing.T) {
	opts := parsePullFlags([]string{"--all"})
	if len(opts.slug) > 0 || len(opts.group) > 0 {
		t.Error("expected empty slug and group")
	}

	if opts.all != true {
		t.Error("expected all=true")
	}
}

func TestParsePullFlags_AllWithVerbose(t *testing.T) {
	opts := parsePullFlags([]string{"--all", "--verbose"})
	if opts.all != true {
		t.Error("expected all=true")
	}

	if opts.verbose != true {
		t.Error("expected verbose=true")
	}
}

func TestParsePullFlags_Parallel(t *testing.T) {
	opts := parsePullFlags([]string{"--all", "--parallel", "4"})
	if opts.parallel != 4 {
		t.Errorf("expected parallel=4, got %d", opts.parallel)
	}
}

func TestParsePullFlags_OnlyAvailable(t *testing.T) {
	opts := parsePullFlags([]string{"--all", "--only-available"})
	if !opts.onlyAvailable {
		t.Error("expected onlyAvailable=true")
	}
}
