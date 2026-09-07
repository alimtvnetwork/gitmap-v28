package store

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func TestProcessLockReentrancy(t *testing.T) {
	tempDir := t.TempDir()

	err := acquireLock(tempDir)
	if err != nil {
		t.Fatalf("first acquireLock failed: %v", err)
	}

	lockFile := filepath.Join(tempDir, constants.LockFileName)
	if !lockExists(lockFile) {
		t.Fatalf("expected lock file %s to exist", lockFile)
	}

	// Second acquire within the same process must succeed immediately (reentrant)
	err2 := acquireLock(tempDir)
	if err2 != nil {
		t.Fatalf("second acquireLock in same process failed: %v", err2)
	}

	// First release drops refcount to 1; lock file must still exist
	releaseLock(tempDir)
	if !lockExists(lockFile) {
		t.Fatalf("expected lock file to persist while ref count > 0")
	}

	// Second release drops refcount to 0; lock file must be removed
	releaseLock(tempDir)
	if lockExists(lockFile) {
		t.Fatalf("expected lock file to be removed when ref count reaches 0")
	}
}

func TestProcessLockStaleRemoval(t *testing.T) {
	tempDir := t.TempDir()
	lockFile := filepath.Join(tempDir, constants.LockFileName)

	// Write an invalid / dead PID
	_ = os.WriteFile(lockFile, []byte("99999999"), 0644)

	err := acquireLock(tempDir)
	if err != nil {
		t.Fatalf("expected acquireLock to clear dead lock and succeed, got: %v", err)
	}
	defer releaseLock(tempDir)

	if !lockExists(lockFile) {
		t.Fatalf("expected lock file to exist after acquiring")
	}
}
