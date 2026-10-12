package store

import (
	"os"
	"path/filepath"
	"testing"
)

func TestIsDirWritable(t *testing.T) {
	tempDir := t.TempDir()
	if !isDirWritable(tempDir) {
		t.Errorf("expected tempDir %q to be writable", tempDir)
	}

	nonExistentUnderWritable := filepath.Join(tempDir, "sub", "deep")
	if !isDirWritable(nonExistentUnderWritable) {
		t.Errorf("expected non-existent subfolder under writable parent to be recognized as writable")
	}
}

func TestFindGlobalFallbackDB_WhenBinaryDirNotWritable(t *testing.T) {
	tempDir := t.TempDir()
	readOnlyDir := filepath.Join(tempDir, "readonly")
	if err := os.MkdirAll(readOnlyDir, 0555); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	defer func() {
		_ = os.Chmod(readOnlyDir, 0755)
	}()

	currentDBPath := filepath.Join(readOnlyDir, "gitmap.db")
	fallback := findGlobalFallbackDB("gitmap.db", currentDBPath)
	if fallback == "" {
		t.Errorf("expected findGlobalFallbackDB to return fallback path when dir is not writable")
	}
}

func TestOpenDefault_ReadOnlyBinaryDirFallback(t *testing.T) {
	tempDir := t.TempDir()
	readOnlyBinDir := filepath.Join(tempDir, "fakebin")
	readOnlyDataDir := filepath.Join(readOnlyBinDir, "data")
	if err := os.MkdirAll(readOnlyDataDir, 0755); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	if err := os.Chmod(readOnlyDataDir, 0555); err != nil {
		t.Fatalf("failed to chmod dir: %v", err)
	}

	defer func() {
		_ = os.Chmod(readOnlyDataDir, 0755)
		SetBinaryDataDirForTesting("")
	}()

	SetBinaryDataDirForTesting(readOnlyDataDir)
	db, err := OpenDefault()
	if err != nil {
		t.Fatalf("expected OpenDefault to succeed with fallback, got error: %v", err)
	}
	if db != nil {
		db.Close()
	}
}

func TestEffectiveDataDir_ReadOnlyFallback(t *testing.T) {
	tempDir := t.TempDir()
	readOnlyDir := filepath.Join(tempDir, "readonly_data")
	if err := os.MkdirAll(readOnlyDir, 0555); err != nil {
		t.Fatalf("failed to create dir: %v", err)
	}
	defer func() {
		_ = os.Chmod(readOnlyDir, 0755)
		SetBinaryDataDirForTesting("")
	}()

	SetBinaryDataDirForTesting(readOnlyDir)
	effective := EffectiveDataDir()
	if effective == readOnlyDir {
		t.Errorf("expected EffectiveDataDir to fall back away from read-only dir, got %q", effective)
	}
}
