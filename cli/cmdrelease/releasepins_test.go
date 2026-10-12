package cmdrelease

import (
	"os"
	"path/filepath"
	"testing"
)

func TestPinnedTagsSaveAndLoad(t *testing.T) {
	tempDir := t.TempDir()

	pins := []string{"v4.162.0", "v4.190.0", "4.162.0"}
	if err := SavePinnedTags(tempDir, pins); err != nil {
		t.Fatalf("SavePinnedTags failed: %v", err)
	}

	loaded, err := LoadPinnedTags(tempDir)
	if err != nil {
		t.Fatalf("LoadPinnedTags failed: %v", err)
	}

	if len(loaded) < 2 {
		t.Errorf("Expected at least 2 deduplicated pins, got %d: %v", len(loaded), loaded)
	}

	path := PinnedTagsFilePath(tempDir)
	if _, err := os.Stat(path); err != nil {
		t.Errorf("Expected pinned-tags.json to exist at %s, got error: %v", path, err)
	}
}

func TestIsTagPinned(t *testing.T) {
	pins := []string{"v4.162.0", "v4.190.0"}

	cases := []struct {
		tag      string
		expected bool
	}{
		{"v4.162.0", true},
		{"4.162.0", true},
		{"v4.190.0", true},
		{"v4.175.0", false},
		{"v4.163.0", false},
		{"", false},
	}

	for _, tc := range cases {
		got := IsTagPinned(tc.tag, pins)
		if got != tc.expected {
			t.Errorf("IsTagPinned(%q) = %v, expected %v", tc.tag, got, tc.expected)
		}
	}
}

func TestPinnedTagsNonExistentFile(t *testing.T) {
	nonExistent := filepath.Join(t.TempDir(), "sub")
	loaded, err := LoadPinnedTags(nonExistent)
	if err != nil {
		t.Errorf("Expected nil error for non-existent file, got %v", err)
	}
	if len(loaded) != 0 {
		t.Errorf("Expected empty slice, got %v", loaded)
	}
}
