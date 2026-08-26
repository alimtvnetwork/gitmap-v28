package cmd

import (
	"testing"
)

func TestCGCLIParseFlags(t *testing.T) {
	opts := parseCGFlags([]string{"version", "repo1", "repo2"})
	if opts.Action != "version" {
		t.Fatalf("expected version action, got %s", opts.Action)
	}
	if len(opts.Repos) != 2 {
		t.Fatalf("expected 2 repos, got %d", len(opts.Repos))
	}

	optsRepo := parseCGFlags([]string{"repo", "update", "repo1"})
	if optsRepo.Action != "update" {
		t.Fatalf("expected update action from repo command, got %s", optsRepo.Action)
	}
}
