package store

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"syscall"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// acquireLock creates an advisory lock file in the given directory.
func acquireLock(dbDir string) error {
	lockPath := filepath.Join(dbDir, constants.LockFileName)

	if lockExists(lockPath) {
		return handleExistingLock(lockPath)
	}

	return writeLock(lockPath)
}

// releaseLock removes the lock file from the given directory.
func releaseLock(dbDir string) {
	lockPath := filepath.Join(dbDir, constants.LockFileName)
	os.Remove(lockPath)
}

// lockExists checks if the lock file is present on disk.
func lockExists(lockPath string) bool {
	_, err := os.Stat(lockPath)

	return err == nil
}

// handleExistingLock reads the PID from the lock and checks liveness.
func handleExistingLock(lockPath string) error {
	pid, err := readLockPID(lockPath)
	if err != nil {
		os.Remove(lockPath)

		return writeLock(lockPath)
	}

	if processRunning(pid) {
		return fmt.Errorf(constants.ErrLockHeld, pid, lockPath)
	}

	os.Remove(lockPath)

	return writeLock(lockPath)
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
func writeLock(lockPath string) error {
	pid := os.Getpid()
	data := []byte(strconv.Itoa(pid))

	return os.WriteFile(lockPath, data, constants.LockFilePermission)
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
