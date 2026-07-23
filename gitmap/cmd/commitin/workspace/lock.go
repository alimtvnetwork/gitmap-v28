package workspace

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
	"syscall"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// LockHandle owns a single advisory lock acquired via AcquireLock.
// Release MUST be called (typically deferred) so the lock file is
// removed even on early-exit paths.
type LockHandle struct {
	Path string
}

// AcquireLock implements spec §3.1 stage 02. Returns *LockHandle on
// success or an error whose .Error() is the spec §2.7 message string
// for CommitInExitLockBusy. Idempotent across crashed predecessors:
// stale lock files (PID no longer running) are reclaimed.
func AcquireLock(p *Paths) (*LockHandle, error) {
	if isLockHeldByLive(p.LockFile) {
		return nil, fmt.Errorf(constants.CommitInErrLockBusy, p.LockFile)
	}
	if err := writeLockPid(p.LockFile); err != nil {
		return nil, fmt.Errorf("commit-in: lock write failed: %w", err)
	}
	return &LockHandle{Path: p.LockFile}, nil
}

// Release deletes the lock file. Errors are intentionally ignored —
// the file may already be gone and the run is shutting down anyway.
func (h *LockHandle) Release() {
	if h == nil || h.Path == "" {
		return
	}
	_ = os.Remove(h.Path)
}

// isLockHeldByLive reports whether the lock file exists AND its PID
// belongs to a running process. Stale lock files (PID dead, or file
// unreadable) are removed and treated as not-held.
//
// Special case: when the recorded PID is the CURRENT process, treat
// the lock as live without probing — signal(0) trivially succeeds
// for the caller's own PID, but more importantly: if WE wrote the
// PID, the lock is by definition held by us. Without this guard a
// second in-process AcquireLock would never return LockBusy, which
// is exactly what TestAcquireLockBlocksDoubleAcquire pins down.
func isLockHeldByLive(lockPath string) bool {
	data, err := os.ReadFile(lockPath)
	if err != nil {
		return false
	}
	pid, parseErr := strconv.Atoi(strings.TrimSpace(string(data)))
	if parseErr != nil {
		_ = os.Remove(lockPath)
		return false
	}
	if pid == os.Getpid() {
		return true
	}
	if isProcessAlive(pid) {
		return true
	}
	_ = os.Remove(lockPath)
	return false
}

// writeLockPid writes the current PID to the lock file. Mode 0o644 to
// match store/lock.go's existing pattern.
func writeLockPid(lockPath string) error {
	pid := strconv.Itoa(os.Getpid())
	return os.WriteFile(lockPath, []byte(pid), 0o644)
}

// isProcessAlive uses signal-0 probing — works on POSIX and Windows
// (Windows os.FindProcess always succeeds; signal(0) returns nil for
// running processes and an error otherwise). EPERM is treated as
// "alive": the process exists, we just lack permission to signal it
// (e.g. PID 1 / root-owned processes when running as a normal user
// in CI). Without this, `kill(0, pid)` on PID 1 from a non-root
// runner falsely reports the process as dead.
func isProcessAlive(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	sigErr := proc.Signal(syscall.Signal(0))
	if sigErr == nil {
		return true
	}
	return errors.Is(sigErr, syscall.EPERM)
}
