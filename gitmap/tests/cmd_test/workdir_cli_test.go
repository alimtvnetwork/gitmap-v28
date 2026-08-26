package cmd_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func TestWorkDirIntegration(t *testing.T) {
	tempDir := t.TempDir()
	db, errDB := store.OpenDefault()
	if errDB != nil {
		t.Skip("sqlite db unavailable in isolated test")
	}
	defer db.Close()

	_, _ = db.SQL().Exec(store.SQLCreateWorkDirsTable)

	wdPath := filepath.Join(tempDir, "work-parent")
	_ = os.MkdirAll(wdPath, 0755)

	wd, errEnsure := db.EnsureWorkDir(wdPath, "test-label", true)
	if errEnsure != nil {
		t.Fatalf("EnsureWorkDir failed: %v", errEnsure)
	}

	if wd.AbsolutePath != wdPath {
		t.Fatalf("expected path %s, got %s", wdPath, wd.AbsolutePath)
	}
}
