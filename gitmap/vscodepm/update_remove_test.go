package vscodepm

import (
	"os"
	"path/filepath"
	"testing"
)

func TestUpdateRootPathAndRemoveEntry(t *testing.T) {
	tempDir := t.TempDir()
	pJson := filepath.Join(tempDir, "projects.json")

	initial := []Entry{
		{Name: "my-app", RootPath: `D:\work\my-app`, Paths: []string{}, Tags: []string{}, Enabled: true},
		{Name: "other-app", RootPath: `D:\work\other-app`, Paths: []string{}, Tags: []string{}, Enabled: true},
	}
	if err := writeEntriesAtomic(pJson, initial); err != nil {
		t.Fatalf("write initial failed: %v", err)
	}

	// Update root path
	if err := UpdateRootPathAt(pJson, `D:/work/my-app`, `D:\work\new-home\my-app`, "my-app-new"); err != nil {
		t.Fatalf("UpdateRootPathAt failed: %v", err)
	}

	entries, err := readEntries(pJson)
	if err != nil {
		t.Fatalf("readEntries failed: %v", err)
	}
	if len(entries) != 2 || entries[0].RootPath != `D:\work\new-home\my-app` || entries[0].Name != "my-app-new" {
		t.Fatalf("Update failed: %+v", entries)
	}

	// Remove entry
	if err := RemoveEntryAt(pJson, `D:/work/other-app`); err != nil {
		t.Fatalf("RemoveEntryAt failed: %v", err)
	}
	entries, err = readEntries(pJson)
	if err != nil {
		t.Fatalf("readEntries after remove failed: %v", err)
	}
	if len(entries) != 1 || entries[0].Name != "my-app-new" {
		t.Fatalf("Remove failed: %+v", entries)
	}
	_ = os.Remove(pJson)
}
