//go:build windows

package fspath

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
)

// FindLockingProcesses detects processes with open handles on the given directory
// if Sysinternals handle.exe is installed on PATH.
func FindLockingProcesses(dirPath string) ([]LockingProcess, error) {
	procs, err := findViaHandle(dirPath)
	if err == nil && len(procs) > 0 {
		return procs, nil
	}
	return nil, nil
}

// findViaHandle uses Sysinternals handle.exe if available on PATH.
func findViaHandle(dirPath string) ([]LockingProcess, error) {
	handlePath, err := exec.LookPath("handle.exe")
	if err != nil {
		handlePath, err = exec.LookPath("handle64.exe")
	}
	if err != nil {
		return nil, fmt.Errorf("handle.exe not found: %w", err)
	}

	out, err := exec.Command(handlePath, "-accepteula", "-nobanner", dirPath).Output()
	if err != nil {
		return nil, err
	}

	return parseHandleOutput(string(out)), nil
}

// parseHandleOutput parses handle.exe output lines like:
// explorer.exe       pid: 1234  type: File  ...
func parseHandleOutput(output string) []LockingProcess {
	seen := make(map[int]bool)
	var procs []LockingProcess

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		pidIdx := strings.Index(line, "pid: ")
		if pidIdx < 0 {
			continue
		}

		name := strings.TrimSpace(line[:pidIdx])
		rest := line[pidIdx+5:]
		fields := strings.Fields(rest)
		if len(fields) == 0 {
			continue
		}

		pid, err := strconv.Atoi(fields[0])
		if err != nil || seen[pid] || IsProtectedProcess(name, pid) {
			continue
		}

		seen[pid] = true
		procs = append(procs, LockingProcess{PID: pid, Name: name})
	}

	return procs
}

// KillProcess terminates a process by PID on Windows using standard os.Process.
func KillProcess(pid int) error {
	if IsProtectedProcess("", pid) {
		return fmt.Errorf("refusing to terminate protected process PID %d", pid)
	}

	proc, err := os.FindProcess(pid)
	if err != nil {
		return err
	}
	return proc.Kill()
}
