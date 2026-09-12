package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	processLockMu        sync.Mutex
	processLockRefCounts = make(map[string]int)
)

// acquireLock creates an advisory lock file in the given directory with a retry backoff.
func acquireLock(dbDir string) error {
	cleanDir := filepath.Clean(dbDir)
	processLockMu.Lock()
	if processLockRefCounts[cleanDir] > 0 {
		processLockRefCounts[cleanDir]++
		processLockMu.Unlock()

		return nil
	}

	processLockMu.Unlock()

	return retryAcquireLock(cleanDir)
}

func retryAcquireLock(cleanDir string) error {
	lockPath := filepath.Join(cleanDir, constants.LockFileName)
	var lastErr error
	for i := 0; i < 50; i++ {
		if lockExists(lockPath) {
			lastErr = handleExistingLock(cleanDir, lockPath)
		} else {
			lastErr = writeLock(cleanDir, lockPath)
		}

		if lastErr == nil {
			return nil
		}

		time.Sleep(100 * time.Millisecond)
	}

	return lastErr
}

// releaseLock removes the lock file from the given directory.
func releaseLock(dbDir string) {
	cleanDir := filepath.Clean(dbDir)
	processLockMu.Lock()
	defer processLockMu.Unlock()

	processLockRefCounts[cleanDir]--
	if processLockRefCounts[cleanDir] <= 0 {
		delete(processLockRefCounts, cleanDir)
		lockPath := filepath.Join(cleanDir, constants.LockFileName)
		_ = os.Remove(lockPath)
	}
}

// lockExists checks if the lock file is present on disk.
func lockExists(lockPath string) bool {
	_, err := os.Stat(lockPath)

	return err == nil
}

// handleExistingLock reads the PID from the lock and checks liveness.
func handleExistingLock(cleanDir, lockPath string) error {
	pid, err := readLockPID(lockPath)
	if err != nil {
		_ = os.Remove(lockPath)

		return writeLock(cleanDir, lockPath)
	}

	if pid == os.Getpid() {
		processLockMu.Lock()
		processLockRefCounts[cleanDir]++
		processLockMu.Unlock()

		return nil
	}

	if processRunning(pid) {
		return fmt.Errorf(constants.ErrLockHeld, pid, lockPath)
	}

	_ = os.Remove(lockPath)

	return writeLock(cleanDir, lockPath)
}

// readLockPID reads and parses the PID from a lock file.
func readLockPID(lockPath string) (int, error) {
	data, err := os.ReadFile(lockPath)
	if err != nil {
		return 0, err
	}

	pid, err := strconv.Atoi(strings.TrimSpace(string(data)))
	if err != nil {
		return 0, err
	}

	return pid, nil
}

// writeLock writes the current process PID to the lock file.
func writeLock(cleanDir, lockPath string) error {
	pid := os.Getpid()
	data := []byte(strconv.Itoa(pid))
	if err := os.WriteFile(lockPath, data, constants.LockFilePermission); err != nil {
		return err
	}

	processLockMu.Lock()
	processLockRefCounts[cleanDir] = 1
	processLockMu.Unlock()

	return nil
}

// processRunning checks if a process with the given PID exists.
func processRunning(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}

	err = proc.Signal(syscall.Signal(0))

	return err == nil
}
