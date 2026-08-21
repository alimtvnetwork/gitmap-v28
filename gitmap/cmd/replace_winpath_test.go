package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestReplaceScanningOnWindowsPaths(t *testing.T) {
	tempDir := t.TempDir()
	filePath := filepath.Join(tempDir, "sample.go")
	content := []byte("package main\n// errorwrapper-v3 is imported\n")
	if err := os.WriteFile(filePath, content, 0644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}

	pair := replacePair{old: "errorwrapper-v3", new: "errorwrapper-v4"}
	hits, total := scanReplacements([]string{filePath}, []replacePair{pair})
	if total != 1 || len(hits) != 1 {
		t.Fatalf("scanReplacements got total=%d, hits=%d; want 1", total, len(hits))
	}

	if string(hits[0].updated) != "package main\n// errorwrapper-v4 is imported\n" {
		t.Errorf("Unexpected updated content: %s", string(hits[0].updated))
	}
}
