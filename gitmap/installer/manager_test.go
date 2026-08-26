// Package installer — manager_test.go tests the installer Manager struct initialization.
package installer

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func TestManagerStruct(t *testing.T) {
	_, errNil := NewManager(nil)
	if errNil == nil {
		t.Fatal("expected error on nil db")
	}

	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_mgr.db")
	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}
	defer db.Close()

	mgr, errNew := NewManager(db)
	if errNew != nil || mgr == nil {
		t.Fatalf("expected NewManager to succeed: %v", errNew)
	}
}
