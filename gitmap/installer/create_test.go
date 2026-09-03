// Package installer — create_test.go tests installer creation business logic.
package installer

import (
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func setupInstallerTestDB(t *testing.T) *store.DB {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_installer.db")
	db, err := store.Open(dbPath)
	if err != nil {
		t.Fatalf("failed to open store: %v", err)
	}

	t.Cleanup(func() { db.Close() })

	if errMigrate := db.MigrateInstallers(); errMigrate != nil {
		t.Fatalf("failed to migrate: %v", errMigrate)
	}

	return db
}

func TestManagerCreate(t *testing.T) {
	db := setupInstallerTestDB(t)
	mgr, _ := NewManager(db)

	script, err := mgr.Create("Docker Suite", "Installs docker engine")
	if err != nil {
		t.Fatalf("Create failed: %v", err)
	}

	if script.Slug != "docker-suite" || script.Version != "v1.0.0" {
		t.Errorf("unexpected script created: %+v", script)
	}

	// Empty name
	_, errEmpty := mgr.Create("", "")
	if errEmpty == nil {
		t.Fatal("expected error on empty name")
	}
}
