package cmdcg

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

func TestCGCLIParseFlagsNoArgs(t *testing.T) {
	opts := parseCGFlags([]string{})
	if opts.Action != "help" {
		t.Fatalf("expected help action on empty args, got %s", opts.Action)
	}
}

func TestCGCLIParseFlagsHelpTokens(t *testing.T) {
	tokens := []string{"help", "-h", "--help"}
	for _, tok := range tokens {
		opts := parseCGFlags([]string{tok})
		if opts.Action != "help" {
			t.Fatalf("expected help action for %s, got %s", tok, opts.Action)
		}
	}
}

func TestCGCLIParseFlagsNewMethods(t *testing.T) {
	actions := []string{"audit", "doctor", "sync", "diff"}
	for _, action := range actions {
		assertParsedAction(t, action)
	}
}

func assertParsedAction(t *testing.T, action string) {
	opts := parseCGFlags([]string{action, "repo1"})
	if opts.Action != action {
		t.Fatalf("expected %s action, got %s", action, opts.Action)
	}

	if len(opts.Repos) != 1 || opts.Repos[0] != "repo1" {
		t.Fatalf("expected repo1 in repos, got %v", opts.Repos)
	}
}

func TestRenderCGHelp(t *testing.T) {
	RenderCGHelp()
}
