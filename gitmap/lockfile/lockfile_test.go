package lockfile

import (
	"errors"
	"os"
	"path/filepath"
	"strconv"
	"testing"
)

// TestAcquire_FreshSucceeds verifies the basic happy path: no lock file,
// Acquire creates one and returns a working Releaser.
func TestAcquire_FreshSucceeds(t *testing.T) {
	name := uniqueName(t)
	release, err := Acquire(name)
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	defer release()

	path := lockPath(name)
	if !fileExists(path) {
		t.Errorf("lock file %s missing after Acquire", path)
	}
}

// TestAcquire_DoubleClaimFails verifies that a second live-process claim
// returns ErrAlreadyHeld. We simulate by manually writing our own PID.
func TestAcquire_DoubleClaimFails(t *testing.T) {
	name := uniqueName(t)
	release, err := Acquire(name)
	if err != nil {
		t.Fatalf("first Acquire: %v", err)
	}
	defer release()

	_, err = Acquire(name)
	if !errors.Is(err, ErrAlreadyHeld) {
		t.Errorf("second Acquire err = %v, want ErrAlreadyHeld", err)
	}
}

// TestAcquire_StalePIDReclaimed verifies the PID-liveness recovery path.
// We write a lock with PID=1 (init/launchd, definitely not gitmap), then
// confirm Acquire steals it.
func TestAcquire_StalePIDReclaimed(t *testing.T) {
	name := uniqueName(t)
	path := lockPath(name)
	// Use a clearly-unrelated PID. PID 1 is always alive but is init/launchd,
	// not a gitmap process — but processRunning only checks liveness, not
	// identity, so we use a definitely-dead PID instead.
	if err := os.WriteFile(path, []byte("999999999"), 0o600); err != nil {
		t.Fatalf("seed lock: %v", err)
	}
	defer os.Remove(path)

	release, err := Acquire(name)
	if err != nil {
		t.Fatalf("Acquire over stale lock: %v", err)
	}
	defer release()

	pid, err := readPID(path)
	if err != nil {
		t.Fatalf("readPID: %v", err)
	}
	if pid != os.Getpid() {
		t.Errorf("PID in lock = %d, want current pid %d", pid, os.Getpid())
	}
}

// TestForceAcquire_OverridesLiveLock verifies the escape hatch works
// even when a "live" PID still holds the lock.
func TestForceAcquire_OverridesLiveLock(t *testing.T) {
	name := uniqueName(t)
	release, err := Acquire(name)
	if err != nil {
		t.Fatalf("first Acquire: %v", err)
	}
	defer release()

	forceRelease, err := ForceAcquire(name)
	if err != nil {
		t.Fatalf("ForceAcquire: %v", err)
	}
	defer forceRelease()

	pid, _ := readPID(lockPath(name))
	if pid != os.Getpid() {
		t.Errorf("PID after ForceAcquire = %d, want %d", pid, os.Getpid())
	}
}

// TestRelease_Idempotent verifies the Releaser can be called twice
// without crashing — important because both `defer release()` and an
// explicit cleanup path may run.
func TestRelease_Idempotent(t *testing.T) {
	name := uniqueName(t)
	release, err := Acquire(name)
	if err != nil {
		t.Fatalf("Acquire: %v", err)
	}
	release()
	release() // must not panic
}

// TestHolderPID_AbsentReturnsZero verifies the missing-file branch.
func TestHolderPID_AbsentReturnsZero(t *testing.T) {
	name := uniqueName(t)
	if got := HolderPID(name); got != 0 {
		t.Errorf("HolderPID on missing lock = %d, want 0", got)
	}
}

// uniqueName returns a per-test lock name + registers cleanup so a
// failed test never leaves stale files behind for the next run.
func uniqueName(t *testing.T) string {
	t.Helper()
	name := "test-" + filepath.Base(t.Name()) + "-" + strconv.Itoa(os.Getpid())
	t.Cleanup(func() { os.Remove(lockPath(name)) })

	return name
}
