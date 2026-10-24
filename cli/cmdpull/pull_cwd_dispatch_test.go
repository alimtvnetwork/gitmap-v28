package cmdpull

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolvePullInvocationTarget_CWD(t *testing.T) {
	cmdName, scope := resolvePullInvocationTarget(false, []string{})
	if cmdName != "pull" {
		t.Fatalf("expected cmdName=pull, got %s", cmdName)
	}
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("unexpected getwd error: %v", err)
	}
	expectedScope := "(" + filepath.Base(cwd) + ")"
	if scope != expectedScope {
		t.Fatalf("expected scope=%s, got %s", expectedScope, scope)
	}
}

func TestResolvePullInvocationTarget_ExplicitRepo(t *testing.T) {
	cmdName, scope := resolvePullInvocationTarget(false, []string{"my-repo"})
	if cmdName != "pull" {
		t.Fatalf("expected cmdName=pull, got %s", cmdName)
	}
	if scope != "(my-repo)" {
		t.Fatalf("expected scope=(my-repo), got %s", scope)
	}
}

func TestResolvePullInvocationTarget_PullAll(t *testing.T) {
	cmdName, scope := resolvePullInvocationTarget(true, []string{})
	if cmdName != "pull-all" {
		t.Fatalf("expected cmdName=pull-all, got %s", cmdName)
	}
	if scope != "(all workspaces)" {
		t.Fatalf("expected scope=(all workspaces), got %s", scope)
	}
}

func TestResolveParsedPullOptions_DefaultsToCWD(t *testing.T) {
	opts := resolveParsedPullOptions([]string{}, false, false)
	if opts.all {
		t.Fatalf("expected opts.all=false for plain pull inside repo, got true")
	}
	if opts.slug != "." {
		t.Fatalf("expected opts.slug='.', got %q", opts.slug)
	}
	if !isPullCWDEnabled(opts) {
		t.Fatalf("expected isPullCWDEnabled=true for opts, got false")
	}
}

func TestResolveDirectPathTarget_CurrentRepo(t *testing.T) {
	cwd, err := os.Getwd()
	if err != nil {
		t.Fatalf("unexpected getwd error: %v", err)
	}
	records := resolveDirectPathTarget(cwd)
	if len(records) != 1 {
		t.Fatalf("expected 1 record for cwd direct path, got %d", len(records))
	}
	if !strings.EqualFold(records[0].RepoName, "gitmap") {
		t.Fatalf("expected repo name gitmap, got %s", records[0].RepoName)
	}
}
