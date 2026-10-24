package cmdpull

import (
	"bytes"
	"path/filepath"
	"strings"
	"testing"
)

func TestRenderFailedGroup_NonGitFolder(t *testing.T) {
	tempDir := t.TempDir()
	states := []*PullRepoState{
		{
			RepoName: "repo-cache",
			RepoPath: filepath.ToSlash(tempDir),
			Changes:  "failed",
			ErrorMsg: "fatal: not a git repository (or any of the parent directories): .git",
		},
	}

	var buf bytes.Buffer
	RenderConciseActiveResultsTo(&buf, states)
	out := buf.String()

	if !strings.Contains(out, "directory exists but is not a Git repository (missing .git)") {
		t.Fatalf("expected non-git folder reason, got:\n%s", out)
	}
	if !strings.Contains(out, "Option 1 (Clone from Upstream):") {
		t.Fatalf("expected Option 1 Clone from Upstream, got:\n%s", out)
	}
	if !strings.Contains(out, "Option 2 (Initialize Repository):") {
		t.Fatalf("expected Option 2 Initialize Repository, got:\n%s", out)
	}
	if !strings.Contains(out, "gitmap fix") {
		t.Fatalf("expected gitmap fix in Option 2, got:\n%s", out)
	}
}

func TestRenderStructuredSolutionsSubtree_SingleOption(t *testing.T) {
	var buf bytes.Buffer
	state := &PullRepoState{
		RepoName: "single-opt-repo",
	}
	renderStructuredSolutionsSubtree(&buf, state)
	out := buf.String()

	if !strings.Contains(out, "├── Solutions:\n") {
		t.Fatalf("expected Solutions header, got:\n%s", out)
	}
	if !strings.Contains(out, "Option 1") {
		t.Fatalf("expected Option 1, got:\n%s", out)
	}
}
