package cmdpull

import (
	"os"
	"strings"
	"testing"
)

func TestResolvePullInvocationHeader_CWD(t *testing.T) {
	header := resolvePullInvocationHeader(false, []string{})
	top := getGitRepoCWDTopLevel()
	if top == "" {
		cwd, _ := os.Getwd()
		top = cwd
	}
	if !strings.Contains(header, top) {
		t.Fatalf("expected header to contain repo path %q, got %q", top, header)
	}
	if strings.Contains(header, "gitmap pull") {
		t.Fatalf("expected header to not repeat 'gitmap pull', got %q", header)
	}
}

func TestResolvePullInvocationHeader_PullAll(t *testing.T) {
	header := resolvePullInvocationHeader(true, []string{})
	if !strings.Contains(header, "pull-all") || !strings.Contains(header, "(all workspaces)") {
		t.Fatalf("expected pull-all (all workspaces), got %q", header)
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
