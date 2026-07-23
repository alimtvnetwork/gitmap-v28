package templates

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const testTag = "lfs/common"

// TestMergeCreatesFileWhenMissing covers the cold-start path: target
// does not exist, Merge writes the block as the only contents.
func TestMergeCreatesFileWhenMissing(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, ".gitattributes")
	body := []byte("*.png filter=lfs diff=lfs merge=lfs -text\n")

	res, err := Merge(target, testTag, body)
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	if res.Outcome != MergeCreated || !res.Changed {
		t.Fatalf("want created+changed, got outcome=%d changed=%v", res.Outcome, res.Changed)
	}

	got, _ := os.ReadFile(target)
	wantPrefix := "# >>> gitmap:lfs/common >>>\n"
	wantSuffix := "# <<< gitmap:lfs/common <<<\n"
	if !bytes.HasPrefix(got, []byte(wantPrefix)) || !bytes.HasSuffix(got, []byte(wantSuffix)) {
		t.Fatalf("unexpected file contents:\n%s", got)
	}
}

// TestMergeIsIdempotentOnSecondRun is the key contract: re-running with
// the same body must NOT touch the file (Changed=false) and must keep
// the block byte-for-byte stable.
func TestMergeIsIdempotentOnSecondRun(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, ".gitattributes")
	body := []byte("*.png filter=lfs diff=lfs merge=lfs -text\n")

	if _, err := Merge(target, testTag, body); err != nil {
		t.Fatalf("first Merge: %v", err)
	}
	first, _ := os.ReadFile(target)

	res, err := Merge(target, testTag, body)
	if err != nil {
		t.Fatalf("second Merge: %v", err)
	}
	if res.Changed {
		t.Fatalf("second run should be a no-op, got Changed=true")
	}

	second, _ := os.ReadFile(target)
	if !bytes.Equal(first, second) {
		t.Fatalf("file drifted between runs:\nfirst:\n%s\nsecond:\n%s", first, second)
	}
}

// TestMergeUpdatesBlockInPlace verifies that changing the body rewrites
// only the block, leaving surrounding hand-written content alone.
func TestMergeUpdatesBlockInPlace(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, ".gitattributes")
	preamble := "# my custom rules\n*.txt text eol=lf\n\n"
	postamble := "\n# user-managed footer\n*.bin binary\n"

	if err := os.WriteFile(target, []byte(preamble), 0o644); err != nil {
		t.Fatal(err)
	}
	if _, err := Merge(target, testTag, []byte("*.png filter=lfs\n")); err != nil {
		t.Fatalf("first Merge: %v", err)
	}

	// Append user footer AFTER the block, then re-run with new body.
	contents, _ := os.ReadFile(target)
	if err := os.WriteFile(target, append(contents, []byte(postamble)...), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := Merge(target, testTag, []byte("*.jpg filter=lfs\n"))
	if err != nil {
		t.Fatalf("second Merge: %v", err)
	}
	if res.Outcome != MergeUpdated || !res.Changed {
		t.Fatalf("want updated+changed, got outcome=%d changed=%v", res.Outcome, res.Changed)
	}

	got, _ := os.ReadFile(target)
	s := string(got)
	if !strings.Contains(s, "# my custom rules") {
		t.Errorf("preamble lost:\n%s", s)
	}
	if !strings.Contains(s, "# user-managed footer") {
		t.Errorf("postamble lost:\n%s", s)
	}
	if strings.Contains(s, "*.png") {
		t.Errorf("old block body survived:\n%s", s)
	}
	if !strings.Contains(s, "*.jpg") {
		t.Errorf("new block body missing:\n%s", s)
	}
}

// TestMergeAppendsBlockToExistingFileWithoutMarkers covers the warm
// path: a hand-written .gitattributes already exists with no gitmap
// block. The block must be appended cleanly with at most one blank
// line of separation.
func TestMergeAppendsBlockToExistingFileWithoutMarkers(t *testing.T) {
	dir := t.TempDir()
	target := filepath.Join(dir, ".gitattributes")
	existing := "*.txt text eol=lf\n"
	if err := os.WriteFile(target, []byte(existing), 0o644); err != nil {
		t.Fatal(err)
	}

	res, err := Merge(target, testTag, []byte("*.png filter=lfs\n"))
	if err != nil {
		t.Fatalf("Merge: %v", err)
	}
	if res.Outcome != MergeInserted || !res.Changed {
		t.Fatalf("want inserted+changed, got outcome=%d changed=%v", res.Outcome, res.Changed)
	}

	got, _ := os.ReadFile(target)
	if !bytes.HasPrefix(got, []byte(existing)) {
		t.Fatalf("existing content not preserved at head:\n%s", got)
	}
	if !bytes.Contains(got, []byte("# >>> gitmap:lfs/common >>>")) {
		t.Fatalf("marker missing:\n%s", got)
	}
}
