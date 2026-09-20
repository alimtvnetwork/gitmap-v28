package osclean

import (
	"os/exec"
	"strings"
)

func runToolCommand(name string, args ...string) (string, bool) {
	path, lookErr := exec.LookPath(name)
	if lookErr != nil || len(path) == 0 {
		return "", false
	}

	out, err := exec.Command(path, args...).CombinedOutput()
	if err != nil {
		return strings.TrimSpace(string(out)), false
	}

	return strings.TrimSpace(string(out)), true
}

func hasTool(name string) bool {
	path, lookErr := exec.LookPath(name)

	return lookErr == nil && len(path) > 0
}
