package cmdide

import (
	"path/filepath"
	"testing"
)

func TestUpdateStateVscdb(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "state.vscdb")

	errFirst := UpdateStateVscdb(dbPath, "file:///d%3A/work/repo1", "repo1")
	if errFirst != nil {
		t.Fatalf("first UpdateStateVscdb failed: %v", errFirst)
	}

	errSecond := UpdateStateVscdb(dbPath, "file:///d%3A/work/repo2", "repo2")
	if errSecond != nil {
		t.Fatalf("second UpdateStateVscdb failed: %v", errSecond)
	}

	conn, errOpen := openVscdb(dbPath)
	if errOpen != nil {
		t.Fatalf("openVscdb failed: %v", errOpen)
	}
	defer conn.Close()

	entries := loadRecentEntries(conn)
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}
	if entries[0].FolderURI != "file:///d%3A/work/repo2" {
		t.Errorf("expected first entry to be repo2, got %s", entries[0].FolderURI)
	}
	if entries[1].FolderURI != "file:///d%3A/work/repo1" {
		t.Errorf("expected second entry to be repo1, got %s", entries[1].FolderURI)
	}

	// Re-add repo1 to ensure it moves to the front without duplicates
	errReadd := UpdateStateVscdb(dbPath, "file:///d%3A/work/repo1", "repo1")
	if errReadd != nil {
		t.Fatalf("re-add UpdateStateVscdb failed: %v", errReadd)
	}

	entriesAfter := loadRecentEntries(conn)
	if len(entriesAfter) != 2 {
		t.Fatalf("expected 2 entries after re-add, got %d", len(entriesAfter))
	}
	if entriesAfter[0].FolderURI != "file:///d%3A/work/repo1" {
		t.Errorf("expected repo1 to move to the front, got %s", entriesAfter[0].FolderURI)
	}
}

func TestUpsertPinnedEntry(t *testing.T) {
	initial := []PinnedProjectEntry{
		{ID: "uuid-1", Name: "alpha", Path: "/work/alpha"},
	}

	// Update existing by path
	updated := upsertPinnedEntry(initial, PinnedProjectEntry{ID: "uuid-1-new", Name: "alpha", Path: "/work/alpha"})
	if len(updated) != 1 {
		t.Fatalf("expected 1 entry after update, got %d", len(updated))
	}
	if updated[0].ID != "uuid-1-new" {
		t.Errorf("expected updated ID, got %s", updated[0].ID)
	}

	// Insert new
	appended := upsertPinnedEntry(updated, PinnedProjectEntry{ID: "uuid-2", Name: "beta", Path: "/work/beta"})
	if len(appended) != 2 {
		t.Fatalf("expected 2 entries after insert, got %d", len(appended))
	}
}
