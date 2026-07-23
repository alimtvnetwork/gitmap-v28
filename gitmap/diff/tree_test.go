package diff

import (
	"os"
	"path/filepath"
	"testing"
)

func TestDiffTrees_AllFourKinds(t *testing.T) {
	left := t.TempDir()
	right := t.TempDir()

	mustWrite(t, filepath.Join(left, "same.txt"), "hello")
	mustWrite(t, filepath.Join(right, "same.txt"), "hello")
	mustWrite(t, filepath.Join(left, "only-left.txt"), "L")
	mustWrite(t, filepath.Join(right, "only-right.txt"), "R")
	mustWrite(t, filepath.Join(left, "conflict.txt"), "version-A")
	mustWrite(t, filepath.Join(right, "conflict.txt"), "version-B")

	entries, err := DiffTrees(left, right, WalkOptions{})
	if err != nil {
		t.Fatalf("DiffTrees: %v", err)
	}
	got := map[string]EntryKind{}
	for _, e := range entries {
		got[e.RelPath] = e.Kind
	}
	if got["same.txt"] != Identical {
		t.Errorf("same.txt: want Identical, got %v", got["same.txt"])
	}
	if got["only-left.txt"] != MissingRight {
		t.Errorf("only-left.txt: want MissingRight, got %v", got["only-left.txt"])
	}
	if got["only-right.txt"] != MissingLeft {
		t.Errorf("only-right.txt: want MissingLeft, got %v", got["only-right.txt"])
	}
	if got["conflict.txt"] != Conflict {
		t.Errorf("conflict.txt: want Conflict, got %v", got["conflict.txt"])
	}
}

func TestDiffTrees_IgnoresGitAndNodeModulesByDefault(t *testing.T) {
	left := t.TempDir()
	right := t.TempDir()
	mustWrite(t, filepath.Join(left, ".git", "HEAD"), "ref: refs/heads/main")
	mustWrite(t, filepath.Join(left, "node_modules", "pkg", "index.js"), "x")

	entries, err := DiffTrees(left, right, WalkOptions{})
	if err != nil {
		t.Fatalf("DiffTrees: %v", err)
	}
	if len(entries) != 0 {
		t.Errorf("expected 0 entries (everything ignored), got %d: %+v", len(entries), entries)
	}
}

func TestSummaryFor_Counts(t *testing.T) {
	in := []Entry{
		{Kind: Conflict}, {Kind: Conflict},
		{Kind: MissingLeft},
		{Kind: MissingRight}, {Kind: MissingRight}, {Kind: MissingRight},
		{Kind: Identical},
	}
	s := SummaryFor(in)
	if s.Conflicts != 2 || s.MissingLeft != 1 || s.MissingRight != 3 || s.Identical != 1 {
		t.Errorf("summary mismatch: %+v", s)
	}
}

// mustWrite creates parent dirs and writes content.
func mustWrite(t *testing.T, path, content string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, []byte(content), 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}
