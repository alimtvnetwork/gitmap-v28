package txn

import (
	"bytes"
	"os"
	"path/filepath"
	"testing"
)

// TestRevertActionsEditFileRestoresBytes proves the typed-action layer
// restores the file byte-for-byte after an edit + commit + revert cycle.
func TestRevertActionsEditFileRestoresBytes(t *testing.T) {
	db := openTempDB(t)
	defer db.Close()

	cwd := t.TempDir()
	target := filepath.Join(cwd, "config.json")
	original := []byte(`{"k":"v"}`)
	mustWriteFile(t, target, original)

	j := mustBegin(t, db, cwd)
	if err := j.RecordEditFile(target); err != nil {
		t.Fatalf("RecordEditFile: %v", err)
	}
	mustWriteFile(t, target, []byte(`{"k":"MUTATED"}`))
	if err := j.Commit(); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	if err := RevertActions(db, j.ID()); err != nil {
		t.Fatalf("RevertActions: %v", err)
	}
	got, _ := os.ReadFile(target)
	if !bytes.Equal(got, original) {
		t.Fatalf("revert restored %q, want %q", got, original)
	}
}

// TestRevertActionsRenamePathRestoresLocation proves rename_path is
// reversed deterministically (live-state validated before the swap).
func TestRevertActionsRenamePathRestoresLocation(t *testing.T) {
	db := openTempDB(t)
	defer db.Close()

	cwd := t.TempDir()
	from := filepath.Join(cwd, "vscode-sync.json")
	to := filepath.Join(cwd, "vscode-sync.disabled.json")
	mustWriteFile(t, from, []byte("payload"))
	if err := os.Rename(from, to); err != nil {
		t.Fatalf("setup rename: %v", err)
	}

	j := mustBegin(t, db, cwd)
	if err := j.RecordRenamePath(from, to); err != nil {
		t.Fatalf("RecordRenamePath: %v", err)
	}
	if err := j.Commit(); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	if err := RevertActions(db, j.ID()); err != nil {
		t.Fatalf("RevertActions: %v", err)
	}
	if _, err := os.Stat(from); err != nil {
		t.Fatalf("expected %q to exist after revert: %v", from, err)
	}
	if _, err := os.Stat(to); !os.IsNotExist(err) {
		t.Fatalf("expected %q to be gone after revert, stat err=%v", to, err)
	}
}

// TestRevertActionsIsIdempotent proves replaying RevertActions on an
// already-reverted transaction is a no-op (no error, no double-restore).
func TestRevertActionsIsIdempotent(t *testing.T) {
	db := openTempDB(t)
	defer db.Close()

	cwd := t.TempDir()
	target := filepath.Join(cwd, "data.txt")
	mustWriteFile(t, target, []byte("v1"))

	j := mustBegin(t, db, cwd)
	if err := j.RecordEditFile(target); err != nil {
		t.Fatalf("RecordEditFile: %v", err)
	}
	mustWriteFile(t, target, []byte("v2"))
	if err := j.Commit(); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	if err := RevertActions(db, j.ID()); err != nil {
		t.Fatalf("RevertActions #1: %v", err)
	}
	mustWriteFile(t, target, []byte("v3-after-revert"))
	if err := RevertActions(db, j.ID()); err != nil {
		t.Fatalf("RevertActions #2 should no-op: %v", err)
	}
	got, _ := os.ReadFile(target)
	if string(got) != "v3-after-revert" {
		t.Fatalf("idempotent revert clobbered live state: got %q", got)
	}
}

// TestRevertActionsMultiStepReversesInSeqDescOrder proves a transaction
// with multiple actions (rename → edit) reverses in Seq DESC order so
// the edit is undone BEFORE the rename, restoring the original repo
// state byte-perfect.
func TestRevertActionsMultiStepReversesInSeqDescOrder(t *testing.T) {
	db := openTempDB(t)
	defer db.Close()

	cwd := t.TempDir()
	from := filepath.Join(cwd, "settings.json")
	to := filepath.Join(cwd, "settings.renamed.json")
	mustWriteFile(t, from, []byte("ORIGINAL"))

	// Forward: rename, then edit the renamed file.
	if err := os.Rename(from, to); err != nil {
		t.Fatalf("setup rename: %v", err)
	}
	j := mustBegin(t, db, cwd)
	if err := j.RecordRenamePath(from, to); err != nil {
		t.Fatalf("RecordRenamePath: %v", err)
	}
	if err := j.RecordEditFile(to); err != nil {
		t.Fatalf("RecordEditFile: %v", err)
	}
	mustWriteFile(t, to, []byte("MUTATED"))
	if err := j.Commit(); err != nil {
		t.Fatalf("Commit: %v", err)
	}

	if err := RevertActions(db, j.ID()); err != nil {
		t.Fatalf("RevertActions: %v", err)
	}
	got, err := os.ReadFile(from)
	if err != nil {
		t.Fatalf("expected %q to exist after multi-step revert: %v", from, err)
	}
	if string(got) != "ORIGINAL" {
		t.Fatalf("multi-step revert restored %q, want %q", got, "ORIGINAL")
	}
}
