// Package lockcheck detects processes that hold file locks on a directory.
package lockcheck

import (
	"fmt"
	"strings"
)

// LockingProcess represents a process that is locking a path.
type LockingProcess struct {
	PID  int
	Name string
}

// String returns a human-readable representation.
func (p LockingProcess) String() string {
	return fmt.Sprintf("%s (PID %d)", p.Name, p.PID)
}

// FormatProcessList returns a formatted string of locking processes.
func FormatProcessList(procs []LockingProcess) string {
	if len(procs) == 0 {
		return ""
	}

	var b strings.Builder
	for i, p := range procs {
		if i > 0 {
			b.WriteString("\n")
		}

		b.WriteString(fmt.Sprintf("  • %s", p.String()))
	}

	return b.String()
}

// IsProtectedProcess returns true if a process PID or name belongs to protected IDE, editor, or current process.
func IsProtectedProcess(name string, pid int) bool {
	if pid <= 0 || pid == 4 { // PID 4 is Windows System
		return true
	}

	return isProtectedName(name)
}

func isProtectedName(name string) bool {
	lower := strings.ToLower(name)
	protectedList := []string{
		"antigravity", "code", "node", "electron",
		"msedgewebview2", "python", "pwsh", "powershell", "bash",
	}
	for _, protected := range protectedList {
		if strings.Contains(lower, protected) {
			return true
		}
	}

	return false
}
