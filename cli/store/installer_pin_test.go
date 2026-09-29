package store

import (
	"path/filepath"
	"testing"
)

func TestInstallerPins_Lifecycle(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test_installation.db")
	db, err := OpenAt(dbPath)
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer db.Close()

	if errMig := RegisterInstallerPinsMigration(db.conn, 1, false); errMig != nil {
		t.Fatalf("failed to register installer pins migration: %v", errMig)
	}

	// 1. Initial lookup should be empty
	ver, err := db.GetPinnedVersion("gitmap")
	if err != nil {
		t.Fatalf("unexpected error querying non-existent pin: %v", err)
	}
	if ver != "" {
		t.Fatalf("expected empty pinned version, got: %s", ver)
	}

	// 2. Pin version
	if errPin := db.PinInstaller("gitmap", "v6.402.0"); errPin != nil {
		t.Fatalf("failed to pin installer: %v", errPin)
	}

	ver, err = db.GetPinnedVersion("gitmap")
	if err != nil {
		t.Fatalf("failed to get pinned version: %v", err)
	}
	if ver != "v6.402.0" {
		t.Fatalf("expected v6.402.0, got: %s", ver)
	}

	// 3. Update pin (upsert)
	if errPin := db.PinInstaller("gitmap", "v6.403.0"); errPin != nil {
		t.Fatalf("failed to update pinned installer: %v", errPin)
	}
	ver, _ = db.GetPinnedVersion("gitmap")
	if ver != "v6.403.0" {
		t.Fatalf("expected v6.403.0 after update, got: %s", ver)
	}

	// 4. Pin second component
	if errPin := db.PinInstaller("agm", "v1.2.0"); errPin != nil {
		t.Fatalf("failed to pin agm: %v", errPin)
	}

	pins, errList := db.ListPinnedVersions()
	if errList != nil {
		t.Fatalf("failed to list pinned versions: %v", errList)
	}
	if len(pins) != 2 {
		t.Fatalf("expected 2 pins, got: %d", len(pins))
	}
	if pins["gitmap"] != "v6.403.0" || pins["agm"] != "v1.2.0" {
		t.Fatalf("unexpected pins map: %+v", pins)
	}

	// 5. Unpin
	if errUnpin := db.UnpinInstaller("gitmap"); errUnpin != nil {
		t.Fatalf("failed to unpin gitmap: %v", errUnpin)
	}
	ver, _ = db.GetPinnedVersion("gitmap")
	if ver != "" {
		t.Fatalf("expected empty version after unpin, got: %s", ver)
	}
}
