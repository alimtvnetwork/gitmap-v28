// Package lockfile provides a tiny, dependency-free advisory file lock
// scoped to a single named operation (e.g. "selfinstall"). It exists to
// prevent two `gitmap self-install` processes from racing each other —
// concurrent installs would otherwise re-prompt the user, double-write
// PATH entries, and overlap binary downloads in the same dir.
//
// Design choices:
//
//   - PID-based, advisory: stale locks (process exited without cleanup)
//     are auto-recovered on the next acquire by checking process liveness.
//   - Lives in os.TempDir(): works even before any gitmap data dir
//     exists, so it can guard the *very first* install.
//   - Exported `Acquire`+`Release` keep the surface tiny; callers always
//     pair them with `defer release()` and never touch the path directly.
//
// Mirrors the established store/lock.go pattern but is exported so cmd/
// can reuse it without importing the whole store package.
package lockfile

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"
)

// ErrAlreadyHeld is returned by Acquire when a live process already
// holds the named lock. Callers check with errors.Is to distinguish it
// from filesystem errors.
var ErrAlreadyHeld = errors.New("lockfile: already held by another process")

// Releaser is the cleanup callback returned by Acquire. Always defer it
// immediately after a successful Acquire.
type Releaser func()

// Acquire takes the named lock or returns ErrAlreadyHeld. The lock file
// is written under os.TempDir() as `gitmap-<name>.lock` and contains
// the holder's PID. If a lock file exists but its PID no longer maps
// to a live process, the stale file is removed and the lock is reacquired.
//
// Returned Releaser is a no-op-safe func: it deletes the file and can
// be called multiple times.
func Acquire(name string) (Releaser, error) {
	path := lockPath(name)
	err := tryClaim(path)
	if err != nil {
		return func() {}, err
	}

	return func() { os.Remove(path) }, nil
}

// ForceAcquire ignores any existing lock (stale or live) and writes a
// fresh one. Used by `--force-lock` to recover from a crashed installer
// that left the file behind but somehow evaded the PID liveness check
// (e.g. PID was recycled by the OS).
func ForceAcquire(name string) (Releaser, error) {
	path := lockPath(name)
	os.Remove(path)
	err := writePIDFile(path)
	if err != nil {
		return func() {}, fmt.Errorf("lockfile: write %s: %w", path, err)
	}

	return func() { os.Remove(path) }, nil
}

// HolderPID returns the PID recorded in the lock file, or 0 if the
// file is absent / unreadable. Used for error messages so the user
// knows which process to investigate.
func HolderPID(name string) int {
	pid, err := readPID(lockPath(name))
	if err != nil {
		return 0
	}

	return pid
}

// lockPath builds the absolute path for a named lock under TempDir.
func lockPath(name string) string {
	return filepath.Join(os.TempDir(), "gitmap-"+name+".lock")
}

// tryClaim is the core acquire logic: write if absent, recover-if-stale
// otherwise. Split out so Acquire stays a one-liner.
func tryClaim(path string) error {
	if !fileExists(path) {
		return writePIDFile(path)
	}

	return recoverOrFail(path)
}

// recoverOrFail inspects an existing lock; replaces it if the holder
// died, returns ErrAlreadyHeld if the holder is alive.
//
// Special case: when the recorded PID is OUR pid, the lock is already
// held by this process — signal(0) trivially succeeds for the caller's
// own PID on every OS, so without this guard a second in-process
// Acquire would falsely "reclaim" its own lock and silently double-
// claim, defeating the whole point of the advisory lock. Tested by
// TestAcquire_DoubleClaimFails.
func recoverOrFail(path string) error {
	pid, err := readPID(path)
	if err != nil {
		// Unparseable lock file → treat as stale.
		os.Remove(path)

		return writePIDFile(path)
	}
	if pid == os.Getpid() {
		return fmt.Errorf("%w (pid=%d, file=%s)", ErrAlreadyHeld, pid, path)
	}
	if processRunning(pid) {
		return fmt.Errorf("%w (pid=%d, file=%s)", ErrAlreadyHeld, pid, path)
	}
	os.Remove(path)

	return writePIDFile(path)
}

// fileExists reports whether path resolves to any filesystem entry.
func fileExists(path string) bool {
	_, err := os.Stat(path)

	return err == nil
}

// writePIDFile writes the current process PID to path with 0o600 perms.
func writePIDFile(path string) error {
	pid := os.Getpid()

	return os.WriteFile(path, []byte(strconv.Itoa(pid)), 0o600)
}

// readPID parses an integer PID from the lock file body.
func readPID(path string) (int, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return 0, err
	}

	return strconv.Atoi(strings.TrimSpace(string(data)))
}

// processRunning checks PID liveness via signal(0). Works on Unix and
// Windows (Go's syscall layer translates the no-op signal correctly).
func processRunning(pid int) bool {
	if pid <= 0 {
		return false
	}
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	err = proc.Signal(syscall.Signal(0))

	return err == nil
}
