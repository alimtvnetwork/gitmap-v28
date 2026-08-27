package release

import (
	"testing"
)

func TestReadWriteLastScannedCommit(t *testing.T) {
	repoDir := t.TempDir()
	commitHash := "abcd123"
	
	if err := WriteLastScannedCommit(repoDir, commitHash); err != nil {
		t.Fatalf("Write failed: %v", err)
	}

	readHash, err := ReadLastScannedCommit(repoDir)
	if err != nil || readHash != commitHash {
		t.Fatalf("Read failed or mismatch. err: %v, got: %s", err, readHash)
	}
}

func TestReadLastScannedCommit_NotFound(t *testing.T) {
	repoDir := t.TempDir()
	_, err := ReadLastScannedCommit(repoDir)
	if err == nil {
		t.Error("Expected error for missing file, got nil")
	}
}
