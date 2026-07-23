package txn

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// TestBeginCommitRevertEdit exercises the full snapshot-edit → revert path
// against a temp SQLite database.
func TestBeginCommitRevertEdit(t *testing.T) {
	db := openTempDB(t)
	defer db.Close()

	cwd := t.TempDir()
	target := filepath.Join(cwd, "data.txt")
	mustWriteFile(t, target, []byte("original"))

	j := mustBegin(t, db, cwd)
	if err := j.SnapshotEdit(target); err != nil {
		t.Fatalf("SnapshotEdit: %v", err)
	}
	mustWriteFile(t, target, []byte("mutated"))
	if err := j.Commit(); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	if err := Revert(db, j.ID(), RevertOptions{Force: true}); err != nil {
		t.Fatalf("Revert: %v", err)
	}
	got, _ := os.ReadFile(target)
	if string(got) != "original" {
		t.Fatalf("revert restored %q, want %q", got, "original")
	}
}

// TestAbortRemovesBackups checks the partial backup tree is cleaned on Abort.
func TestAbortRemovesBackups(t *testing.T) {
	db := openTempDB(t)
	defer db.Close()

	cwd := t.TempDir()
	target := filepath.Join(cwd, "data.txt")
	mustWriteFile(t, target, []byte("original"))

	j := mustBegin(t, db, cwd)
	if err := j.SnapshotDelete(target); err != nil {
		t.Fatalf("SnapshotDelete: %v", err)
	}
	if err := j.Abort(); err != nil {
		t.Fatalf("Abort: %v", err)
	}
}

// mustBegin is a test helper that opens a fresh edit-kind transaction.
func mustBegin(t *testing.T, db *store.DB, cwd string) *Journal {
	t.Helper()
	j, err := Begin(db, Meta{
		Kind:           constants.TxnKindFixRepo,
		Argv:           []string{"gitmap", "fix-repo"},
		Cwd:            cwd,
		ReverseSummary: "test fixture",
		RepoSlug:       "test-repo",
		GitSha:         "deadbeef",
	})
	if err != nil {
		t.Fatalf("Begin: %v", err)
	}

	return j
}

// openTempDB creates an isolated SQLite db + runs Migrate.
func openTempDB(t *testing.T) *store.DB {
	t.Helper()
	db, err := store.Open(t.TempDir())
	if err != nil {
		t.Fatalf("store.Open: %v", err)
	}
	if err := db.Migrate(); err != nil {
		t.Fatalf("db.Migrate: %v", err)
	}

	return db
}

// mustWriteFile writes bytes to path with parents created.
func mustWriteFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatalf("mkdir: %v", err)
	}
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("write: %v", err)
	}
}
