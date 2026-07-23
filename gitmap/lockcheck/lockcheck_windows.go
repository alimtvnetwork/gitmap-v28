//go:build windows

package lockcheck

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
	"syscall"
)

// FindLockingProcesses detects processes with open handles on the given directory.
// Uses PowerShell to query via handle.exe or falls back to tasklist-based heuristics.
func FindLockingProcesses(dirPath string) ([]LockingProcess, error) {
	// Try using PowerShell's Get-Process with file handle check.
	// This approach uses handle.exe from Sysinternals if available,
	// otherwise falls back to checking common culprits via WMI.
	procs, err := findViaHandle(dirPath)
	if err == nil && len(procs) > 0 {
		return procs, nil
	}

	// Fallback: use PowerShell WMI query for processes with the directory in their path.
	return findViaWMI(dirPath)
}

// findViaHandle uses Sysinternals handle.exe if available on PATH.
func findViaHandle(dirPath string) ([]LockingProcess, error) {
	handlePath, err := exec.LookPath("handle.exe")
	if err != nil {
		handlePath, err = exec.LookPath("handle64.exe")
		if err != nil {
			return nil, fmt.Errorf("handle.exe not found")
		}
	}

	cmd := exec.Command(handlePath, "-accepteula", "-nobanner", dirPath)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
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

		// handle.exe format: "processname pid: NNNN ..."
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
		if err != nil || seen[pid] {
			continue
		}

		seen[pid] = true
		procs = append(procs, LockingProcess{PID: pid, Name: name})
	}

	return procs
}

// findViaWMI uses PowerShell to find processes whose executable path or
// current directory overlaps with the target directory.
func findViaWMI(dirPath string) ([]LockingProcess, error) {
	// Escape backslashes for WMI LIKE query.
	escaped := strings.ReplaceAll(dirPath, `\`, `\\`)

	script := fmt.Sprintf(
		`Get-CimInstance Win32_Process | Where-Object { $_.ExecutablePath -like '%s*' -or $_.CommandLine -like '*%s*' } | Select-Object ProcessId, Name | ForEach-Object { "$($_.ProcessId)|$($_.Name)" }`,
		escaped, escaped,
	)

	cmd := exec.Command("powershell", "-NoProfile", "-Command", script)
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	out, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	return parseWMIOutput(string(out)), nil
}

func parseWMIOutput(output string) []LockingProcess {
	seen := make(map[int]bool)
	var procs []LockingProcess

	for _, line := range strings.Split(output, "\n") {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, "|", 2)
		if len(parts) != 2 {
			continue
		}
		pid, err := strconv.Atoi(parts[0])
		if err != nil || seen[pid] {
			continue
		}
		seen[pid] = true
		procs = append(procs, LockingProcess{PID: pid, Name: parts[1]})
	}

	return procs
}

// KillProcess terminates a process by PID on Windows.
func KillProcess(pid int) error {
	cmd := exec.Command("taskkill", "/F", "/PID", strconv.Itoa(pid))
	cmd.SysProcAttr = &syscall.SysProcAttr{HideWindow: true}
	return cmd.Run()
}
