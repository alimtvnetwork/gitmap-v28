package release

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// TestHandleOrphanedMetaConfirm verifies that answering "y" removes the
// orphaned release JSON file and returns nil.
func TestHandleOrphanedMetaConfirm(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	v, _ := Parse("v2.3.10")
	meta := ReleaseMeta{Tag: "v2.3.10", Version: "2.3.10"}
	if err := WriteReleaseMeta(meta); err != nil {
		t.Fatalf("WriteReleaseMeta: %v", err)
	}

	if !ReleaseExists(v) {
		t.Fatal("release metadata should exist before test")
	}

	// Simulate stdin with "y\n".
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	w.WriteString("y\n")
	w.Close()

	origStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	err = handleOrphanedMeta(v)
	if err != nil {
		t.Fatalf("handleOrphanedMeta returned error: %v", err)
	}

	if ReleaseExists(v) {
		t.Error("release metadata should be removed after confirmation")
	}
}

// TestHandleOrphanedMetaDecline verifies that answering "n" aborts
// the release and leaves the JSON file intact.
func TestHandleOrphanedMetaDecline(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	v, _ := Parse("v2.3.10")
	meta := ReleaseMeta{Tag: "v2.3.10", Version: "2.3.10"}
	if err := WriteReleaseMeta(meta); err != nil {
		t.Fatalf("WriteReleaseMeta: %v", err)
	}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	w.WriteString("n\n")
	w.Close()

	origStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	err = handleOrphanedMeta(v)
	if err == nil {
		t.Fatal("expected error on decline, got nil")
	}

	if err.Error() != constants.ErrReleaseAborted {
		t.Errorf("expected %q, got %q", constants.ErrReleaseAborted, err.Error())
	}

	if !ReleaseExists(v) {
		t.Error("release metadata should still exist after decline")
	}
}

// TestHandleOrphanedMetaEOF verifies that EOF (no input) aborts
// with the standard "already released" error.
func TestHandleOrphanedMetaEOF(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	v, _ := Parse("v2.3.10")
	meta := ReleaseMeta{Tag: "v2.3.10", Version: "2.3.10"}
	if err := WriteReleaseMeta(meta); err != nil {
		t.Fatalf("WriteReleaseMeta: %v", err)
	}

	// Empty pipe → EOF on first Scan.
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	w.Close()

	origStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	err = handleOrphanedMeta(v)
	if err == nil {
		t.Fatal("expected error on EOF, got nil")
	}

	if !ReleaseExists(v) {
		t.Error("release metadata should still exist after EOF")
	}
}

// TestOrphanedMetaFileRemoval verifies the file is physically deleted
// from the .release directory after confirmation.
func TestOrphanedMetaFileRemoval(t *testing.T) {
	cleanup := setupTempDir(t)
	defer cleanup()

	v, _ := Parse("v5.0.0")
	meta := ReleaseMeta{Tag: "v5.0.0", Version: "5.0.0"}
	if err := WriteReleaseMeta(meta); err != nil {
		t.Fatalf("WriteReleaseMeta: %v", err)
	}

	path := filepath.Join(constants.DefaultReleaseDir, "v5.0.0.json")
	if _, err := os.Stat(path); os.IsNotExist(err) {
		t.Fatal("metadata file should exist on disk")
	}

	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("os.Pipe: %v", err)
	}
	w.WriteString("yes\n")
	w.Close()

	origStdin := os.Stdin
	os.Stdin = r
	defer func() { os.Stdin = origStdin }()

	if err := handleOrphanedMeta(v); err != nil {
		t.Fatalf("handleOrphanedMeta: %v", err)
	}

	if _, err := os.Stat(path); !os.IsNotExist(err) {
		t.Error("metadata file should not exist on disk after removal")
	}
}
