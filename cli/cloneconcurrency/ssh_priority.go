package cloneconcurrency

import (
	"os"
	"runtime"
)

// IsSSHSession returns true if running inside an SSH session or delegated fleet execution.
func IsSSHSession() bool {
	return os.Getenv("SSH_CLIENT") != "" || os.Getenv("SSH_CONNECTION") != "" || os.Getenv("GITMAP_SSH_DELEGATED") != ""
}

// ResolveWithPriority translates concurrency respecting SSH half-priority throttling.
func ResolveWithPriority(n int, isSSHRequest bool) (int, bool) {
	if n < 0 {
		return 0, false
	}
	if n > 0 {
		return n, true
	}
	if isSSHRequest || IsSSHSession() {
		return resolveHalfPriority(), true
	}

	return resolveStandardPriority(), true
}

func resolveHalfPriority() int {
	half := runtime.NumCPU() / 2
	if half < 1 {
		return 1
	}
	if half > 2 {
		return 2
	}

	return half
}

func resolveStandardPriority() int {
	w := runtime.NumCPU()
	if w < 1 {
		return 1
	}

	return w
}
