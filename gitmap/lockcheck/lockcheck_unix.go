//go:build !windows

package lockcheck

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

// FindLockingProcesses uses lsof to detect processes with open files in dirPath.
func FindLockingProcesses(dirPath string) ([]LockingProcess, error) {
	lsofPath, err := exec.LookPath("lsof")
	if err != nil {
		return nil, fmt.Errorf("lsof not found")
	}

	cmd := exec.Command(lsofPath, "+D", dirPath)
	out, _ := cmd.Output() // lsof returns exit 1 when no results

	return parseLsofOutput(string(out)), nil
}

func parseLsofOutput(output string) []LockingProcess {
	seen := make(map[int]bool)
	var procs []LockingProcess

	for i, line := range strings.Split(output, "\n") {
		if i == 0 || strings.TrimSpace(line) == "" {
			continue // skip header
		}
		fields := strings.Fields(line)
		if len(fields) < 2 {
			continue
		}
		name := fields[0]
		pid, err := strconv.Atoi(fields[1])
		if err != nil || seen[pid] {
			continue
		}
		seen[pid] = true
		procs = append(procs, LockingProcess{PID: pid, Name: name})
	}

	return procs
}

// KillProcess terminates a process by PID on Unix.
func KillProcess(pid int) error {
	return exec.Command("kill", "-9", strconv.Itoa(pid)).Run()
}
