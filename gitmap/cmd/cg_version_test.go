package cmd

import (
	"testing"
)

func TestCGVersion(t *testing.T) {
	tempDir := t.TempDir()

	// Initially not installed
	if _, err := ReadCGMetadata(tempDir); err == nil {
		t.Fatal("expected error for uninstalled repo")
	}

	meta := CGMetadata{
		Version: "v24.0.0",
		Status:  "active",
	}
	if err := WriteCGMetadata(tempDir, meta); err != nil {
		t.Fatalf("WriteCGMetadata failed: %v", err)
	}

	read, errRead := ReadCGMetadata(tempDir)
	if errRead != nil {
		t.Fatalf("ReadCGMetadata failed: %v", errRead)
	}

	if read.Version != "v24.0.0" {
		t.Fatalf("unexpected version: %s", read.Version)
	}

	// Test print functions
	PrintCGStatus([]string{tempDir})
	PrintCGVersion([]string{tempDir})
}
